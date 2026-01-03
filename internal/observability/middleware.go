package observability

import (
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		reqID := ExtractOrGenerateRequestID(r)
		ctx := WithRequestID(r.Context(), reqID)

		r = r.WithContext(ctx)
		w.Header().Set(RequestIDHeader, reqID)

		next.ServeHTTP(w, r)

		latency := time.Since(start)

		Logger.Printf(
			`request_id=%s method=%s path=%s remote_addr=%s latency_ms=%d`,
			reqID,
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			latency.Microseconds(),
		)
	})
}
