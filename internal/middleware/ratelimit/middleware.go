package ratelimit

import (
	"net/http"

	"github.com/Rugved7/api-gateway/internal/observability"
)

func Middleware(limiter *Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := observability.GetRequestID(r.Context())
			if key == "" {
				key = r.RemoteAddr
			}

			if !limiter.allow(key) {
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
