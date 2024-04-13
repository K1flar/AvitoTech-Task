package timelimit

import (
	"context"
	"net/http"
	"time"
)

func New(responseTime time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), responseTime)
			defer cancel()

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
