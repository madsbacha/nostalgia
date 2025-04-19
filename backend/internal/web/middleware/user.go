package middleware

import (
	"context"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
)

var (
	UserIdCtxKey = &contextKey{"UserId"}
)

func SetUserIdFromJwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		if userId, ok := claims["user_id"].(string); ok {
			ctx = context.WithValue(ctx, UserIdCtxKey, userId)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserIdFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(UserIdCtxKey).(string); ok {
		return v
	}
	return ""
}
