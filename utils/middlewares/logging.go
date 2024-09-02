package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func (m *Middleware) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)
		m.Log.Println(wrapped.statusCode, r.Method, r.URL.Path, fmt.Sprintf("%v ms", time.Since(start).Milliseconds()))
	})
}

func (m *Middleware) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (m *Middleware) Println(v ...interface{}) {
	log.Println(v...)
}

type wrappedResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
