package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTManager struct {
	Key    []byte
	TTL    time.Duration
	Issuer string
}

func (mgr *JWTManager) Verify(token string) (*Identity, error) {
	data := &Token{}
	_, err := jwt.ParseWithClaims(token, data, func(t *jwt.Token) (interface{}, error) {
		return mgr.Key, nil
	})

	if err != nil {
		return nil, err
	}

	if data.Issuer != mgr.Issuer {
		return nil, errors.New("invalid issuer")
	}

	if err != nil {
		return nil, err
	}

	return data.Identity, nil
}

func (mgr *JWTManager) New(identity *Identity) string {
	now := time.Now()

	claims := &Token{
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: now.Add(mgr.TTL).Unix(),
			IssuedAt:  now.Unix(),
			Issuer:    mgr.Issuer,
		},
		Identity: identity,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString(mgr.Key)
	return signed
}
