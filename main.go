package main

import (
	"encoding/json"
	"fmt"
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
		Key:    []byte("secret"),
		TTL:    time.Hour * 24 * 7,
		Issuer: "some-issuer-name",
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
			AuthorizedEndpoint: "http://localhost:3001",
		},
	)

	http.HandleFunc("/auth/discord", discordRedirect)
	http.HandleFunc("/auth/discord/callback", discordCallback)

	http.HandleFunc("/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		cookie, err := r.Cookie("token")

		fmt.Println(1, err)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized: token missing"))
			return
		}

		fmt.Println(2)

		identity, err := jwt.Verify(cookie.Value)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized: token invalid"))
			return
		}

		fmt.Println(3)

		data, _ := json.Marshal(identity)
		fmt.Println(string(data))

		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(data))
	})

	http.ListenAndServe("localhost:3000", nil)
}
