package server

import (
	"net/http"
)

// Chain creates a middleware chain from multiple middleware functions
func Chain(next http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, mw := range middlewares {
		next = mw(next)
	}
	return next
}
