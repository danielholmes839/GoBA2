package auth

import (
	"context"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

type URLFunc func(token string) url.URL

type OAuth2EndpointsConfig struct {
	IdentityFunc
	TokenProvider
	AuthorizedEndpoint string
}

func OAuth2Endpoints(provider *oauth2.Config, config *OAuth2EndpointsConfig) (redirect, callback http.HandlerFunc) {
	state := "optional" // TODO

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
		identity, err := config.IdentityFunc(client)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("identity failed"))
			return
		}

		token := config.TokenProvider.New(identity)

		http.SetCookie(w, &http.Cookie{
			Name:  "token",
			Value: token,
			Path:  "/",
		})

		http.Redirect(w, r, config.AuthorizedEndpoint, http.StatusSeeOther)
	}

	return redirect, callback
}
