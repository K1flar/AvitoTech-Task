package auth

import (
	"banner_service/internal/domains"
	tokenManager "banner_service/internal/token_manager"
	"context"
	"net/http"
)

func New(tm *tokenManager.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			inputToken := r.Header.Get("token")
			if inputToken == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			token := tm.Parse(inputToken)
			if !token.Valid {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), domains.RoleKey("role"), token.GetRole())
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
