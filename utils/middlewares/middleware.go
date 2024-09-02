package middlewares

import (
	"log"
	"net/http"
	"os"
)

type Middleware struct {
	Log *log.Logger
}

func NewMiddleware() *Middleware {
	return &Middleware{
		Log: log.New(os.Stdout, "", log.LstdFlags),
	}
}

type MiddlewareFunc func(http.Handler) http.Handler

func CreateStack(xs ...MiddlewareFunc) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}

		return next
	}
}
