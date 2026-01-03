// Package middleware allows to chain the middleware
package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler
