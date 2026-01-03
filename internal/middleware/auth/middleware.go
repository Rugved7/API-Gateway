package auth

import (
	"context"
	"net/http"
	"strings"
)

type claimsKeyType struct{}

var claimsKey = claimsKeyType{}

func WithClaims(ctx context.Context, claims map[string]any) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

func ClaimsFromContext(ctx context.Context) map[string]any {
	if v := ctx.Value(claimsKey); v != nil {
		if claims, ok := v.(map[string]any); ok {
			return claims
		}
	}
	return nil
}

func Middleware(validator *Validator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			claims, err := validator.Validate(parts[1])
			if err != nil {
				http.Error(w, "invalid token", http.StatusForbidden)
				return
			}

			ctx := WithClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
