package main

import (
	"encoding/json"
	"goba2/backend/auth"
	"goba2/game"
	"net/http"
	"time"

	"github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"
)

func main() {

	http.HandleFunc("/game/connect", game.GameEndpoint())

	jwt := &auth.JWTManager{
		Key: []byte("secret"),
		TTL: time.Hour * 24 * 7,
	}

	discordRedirect, discordCallback := auth.OAuth2Endpoints(
		&oauth2.Config{
			Endpoint:     discord.Endpoint,
			Scopes:       []string{discord.ScopeIdentify},
			RedirectURL:  "http://localhost:3000/auth/discord/callback",
			ClientID:     "956783521526063204",
			ClientSecret: "GX1ftHqewT7K0ARw87NLtcaCfijqk8Pq", // TODO DELETE THIS BEFORE PROD
		},
		&auth.OAuth2EndpointsConfig{
			IdentityFunc:       auth.DiscordIdentity,
			TokenProvider:      jwt,
			AuthorizedEndpoint: "/me",
		},
	)

	http.HandleFunc("/auth/discord", discordRedirect)
	http.HandleFunc("/auth/discord/callback", discordCallback)

	http.HandleFunc("/me", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized: \"token\" cookie missing"))
			return
		}

		identity, err := jwt.Verify(cookie.Value)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized: cookie invalid"))
			return
		}

		data, _ := json.Marshal(identity)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(data))
	})

	http.HandleFunc("/auth/check", func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		identity, err := jwt.Verify(token)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
			return
		}

		data, _ := json.Marshal(identity)
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})

	http.ListenAndServe("localhost:3000", nil)
}
