package middleware

import (
	"context"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
)

var (
	FileIdCtxKey = &contextKey{"FileId"}
)

func SetFileIdFromJwt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := r.Context()
		if fileId, ok := claims["file_id"].(string); ok {
			ctx = context.WithValue(ctx, FileIdCtxKey, fileId)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func FileIdFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(FileIdCtxKey).(string); ok {
		return v
	}
	return ""
}
