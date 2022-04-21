package backend

import (
	"context"
	"errors"
	"goba2/backend/auth"
	"net/http"
)

type Middleware func(next http.Handler) http.Handler

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		next.ServeHTTP(w, r)
	})
}

func AuthenticationMiddleware(verifier auth.TokenVerifier) Middleware {
	// Add an "identity" key containing an *auth.Identity to the request context
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read the cookie
			identity, err := auth.GetIdentity(verifier, r)

			if err != nil {
				// the token could not be verifier
				write(w, http.StatusUnauthorized, err.Error())
				return
			}

			// add the identity to the context
			ctx := context.WithValue(r.Context(), "identity", identity)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type AuthenticatedHandlerFunc func(w http.ResponseWriter, r *http.Request, identity *auth.Identity)

// Adds the auth.Identity to the handler func
func AuthHandler(handler AuthenticatedHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		identity, ok := r.Context().Value("identity").(*auth.Identity)
		if !ok {
			panic(errors.New("identity could not be found check the router setup. add auth middleware"))
		}

		handler(w, r, identity)
	}
}
