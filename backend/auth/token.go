package auth

import (
	"net/http"

	"github.com/golang-jwt/jwt"
)

type Token struct {
	*jwt.StandardClaims
	*Identity `json:"identity"`
}

type TokenVerifier interface {
	Verify(token string) (*Identity, error)
}

type TokenProvider interface {
	New(identity *Identity) string
}

type Identity struct {
	Provider string `json:"provider"`  // discord maybe google
	UserID   string `json:"user_id"`   // 126715321...
	Username string `json:"user_name"` // daniels_ego
	AvatarID string `json:"avatar_id"` // 12386127361...
	Color    string `json:"color"`     // #ffffff
}

type IdentityFunc func(client *http.Client) (*Identity, error)
