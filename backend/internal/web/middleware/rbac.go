package middleware

import (
	"net/http"
	"nostalgia/internal/common/rbac"
)

type RbacMiddleware struct {
	rbac rbac.Context
}

func NewRbacMiddleware(rbac rbac.Context) RbacMiddleware {
	return RbacMiddleware{rbac: rbac}
}

func (m RbacMiddleware) Guard(condition rbac.Condition) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userId := UserIdFromContext(r.Context())
			if userId == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			roles, err := m.rbac.GetUserRoles(r.Context(), userId)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if !condition(roles...) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
