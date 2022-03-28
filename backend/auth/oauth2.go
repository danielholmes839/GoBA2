package auth

import (
	"context"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

type URLFunc func(token string) url.URL

type OAuth2EndpointConfig struct {
	IdentityFunc
	TokenProvider
	AuthorizedEndpoint string
}

type OAuth2Config struct {
	Provider *oauth2.Config
	Endpoint *OAuth2EndpointConfig
}

func OAuth2Endpoints(conf *OAuth2Config) (redirect, callback http.HandlerFunc) {
	state := "optional" // TODO
	provider := conf.Provider
	endpoint := conf.Endpoint

	redirect = func(w http.ResponseWriter, r *http.Request) {
		// redirect to the OAuth2 provider endpoint
		http.Redirect(w, r, provider.AuthCodeURL(state), http.StatusTemporaryRedirect)
	}

	callback = func(w http.ResponseWriter, r *http.Request) {
		// receive token
		if r.FormValue("state") != state {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("State does not match."))
			return
		}

		// get an http client using the OAuth2 token
		code := r.URL.Query().Get("code")
		providerToken, _ := provider.Exchange(context.Background(), code)
		client := provider.Client(context.Background(), providerToken)

		// use the client to get an identity
		identity, err := endpoint.IdentityFunc(client)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("identity failed"))
			return
		}

		token := endpoint.TokenProvider.New(identity)

		http.SetCookie(w, &http.Cookie{
			Name:  "token",
			Value: token,
			Path:  "/",
		})

		http.Redirect(w, r, endpoint.AuthorizedEndpoint, http.StatusSeeOther)
	}

	return redirect, callback
}
