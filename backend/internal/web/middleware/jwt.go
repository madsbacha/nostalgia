package middleware

import (
	"github.com/go-chi/jwtauth/v5"
	"net/http"
)

func Verifier(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return jwtauth.Verify(ja, jwtauth.TokenFromHeader, jwtauth.TokenFromQuery, jwtauth.TokenFromCookie)(next)
	}
}
