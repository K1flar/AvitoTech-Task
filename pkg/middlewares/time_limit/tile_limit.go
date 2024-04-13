package timelimit

import (
	"banner_service/internal/config"
	"context"
	"net/http"
)

func New(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), cfg.Server.ResponseTime)
			defer cancel()

			done := make(chan struct{})

			go func() {
				next.ServeHTTP(w, r)
				done <- struct{}{}
			}()

			select {
			case <-ctx.Done():
				http.Error(w, "Request Timeout", http.StatusRequestTimeout)
			case <-done:
				return
			}
		})
	}
}
