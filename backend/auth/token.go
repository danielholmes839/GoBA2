package auth

import (
	"errors"
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

func (identity *Identity) ID() string {
	return identity.UserID
}

/* IdentityFunc requests the identity information from the oauth provider using an authenticated http client */
type IdentityFunc func(client *http.Client) (*Identity, error)

func GetIdentity(verifier TokenVerifier, r *http.Request) (*Identity, error) {
	cookie, err := r.Cookie("token")

	if err != nil {
		// the token does exist
		return nil, errors.New("token does not exist")
	}

	identity, err := verifier.Verify(cookie.Value)

	if err != nil {
		// the token could not be verifier
		return nil, errors.New("token not valid")
	}

	return identity, err
}
