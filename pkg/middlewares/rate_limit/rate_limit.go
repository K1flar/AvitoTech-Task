package ratelimit

import (
	"net/http"
	"time"
)

func New(rps int) func(http.Handler) http.Handler {
	requests := make(chan struct{}, rps)
	ticker := time.NewTicker(time.Second)
	go func() {
		for range ticker.C {
			for len(requests) > 0 {
				<-requests
			}
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			select {
			case requests <- struct{}{}:
				next.ServeHTTP(w, r)
			default:
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
		})
	}
}
