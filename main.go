package main

import (
	"goba2/backend"
	"goba2/backend/auth"
	"net/http"
	"time"

	"github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"
)

func main() {
	jwt := &auth.JWTManager{
		Key:    []byte("secret"),
		TTL:    time.Hour * 24 * 7,
		Issuer: "some-issuer-name",
	}

	server := backend.NewServer(&backend.Config{
		TokenVerifier: jwt,
		Discord: &auth.OAuth2Config{
			Provider: &oauth2.Config{
				Endpoint:     discord.Endpoint,
				Scopes:       []string{discord.ScopeIdentify},
				RedirectURL:  "http://localhost:3000/auth/discord/callback",
				ClientID:     "956783521526063204",
				ClientSecret: "GX1ftHqewT7K0ARw87NLtcaCfijqk8Pq", // TODO DELETE THIS BEFORE PROD
			},
			Endpoint: &auth.OAuth2EndpointConfig{
				IdentityFunc:       auth.DiscordIdentity,
				TokenProvider:      jwt,
				AuthorizedEndpoint: "http://localhost:3001",
			},
		},
	})

	router := server.Routes()
	http.ListenAndServe("localhost:3000", router)
}
