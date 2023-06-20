package core

import (
	"net/http"
	"strings"
)

func Authenticator(api_token string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := TokenFromHeader(r)

			if token == "" || token != api_token {
				http.Error(w, "", http.StatusNotFound)
				return
			}

			// Token is authenticated, pass it through
			next.ServeHTTP(w, r)
		})
	}
}

func TokenFromHeader(r *http.Request) string {
	// Get token from authorization header.
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}
