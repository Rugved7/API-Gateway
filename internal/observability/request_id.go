package observability

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

type requestIDKeyType struct{}

var requestIDKey = requestIDKeyType{}

const RequestIDHeader = "X-Request-ID"

func GenerateRequestID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func GetRequestID(ctx context.Context) string {
	if v := ctx.Value(requestIDKey); v != nil {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return ""
}

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

func ExtractOrGenerateRequestID(r *http.Request) string {
	if id := r.Header.Get(RequestIDHeader); id != "" {
		return id
	}
	return GenerateRequestID()
}
