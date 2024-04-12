package adminmw

import (
	"banner_service/internal/domains"
	"net/http"
)

func New() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(domains.RoleKey("role")).(domains.Role)
			if !ok || role != domains.AdminRole {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
