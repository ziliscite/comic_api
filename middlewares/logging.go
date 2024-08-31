package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Logger struct {
	Log *log.Logger
}

func (l *Logger) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)
		l.Log.Println(wrapped.statusCode, r.Method, r.URL.Path, fmt.Sprintf("%v ms", time.Since(start).Milliseconds()))
	})
}

func (l *Logger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (l *Logger) Println(v ...interface{}) {
	log.Println(v...)
}

// To get http response status code of the response
// Wrapping wrapping stuff
type wrappedResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// We can use this in the actual router is the way I see it
// Oh.. well, it is the same named function that actually writes the header
// We just wrap it with another layer

func (w *wrappedResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode

	// This
	w.ResponseWriter.WriteHeader(statusCode)
}

// Middleware middlewares chaining
type Middleware func(http.Handler) http.Handler

func CreateStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}

		return next
	}
}

/* To use it
// The first should be logging btw
// Because the stack will go from the bottom to top
stack := middlewares.CreateStack(
	middlewares.Logging, // Called last in chain
	middlewares.Something, // Called 3rd
	middlewares.Something, // Called 2nd
	middlewares.Something, // Called 1st
)

server := http.Server{
	Addr: ":8080"
	Handler: stack(router)
}
*/

// SubRouting
/*
handler := &something.Handler{}

// Router
router := http.NewServerMux()
router.HandlerFunc("POST /books", handler.Create)
router.HandlerFunc("GET /books/{id}", handler.Get)

// Sub
v1 := http.NewServerMux()
v1.Handle("/v1/", http.StripPrefix("/v1", router))
// All the router will have /v1 prefix

// /v1/books

server := http.Server {
	Addr: ":8080",
	Handler: router,
}
*/

/*
SubRouter can also be used with middlewares

handler := &invoice.Handler{}

// Router
router := http.NewServerMux()

// Public
router.HandlerFunc("GET /invoice/{id}", handler.Get)

adminRouter := http.NewServerMux()
adminRouter.HandlerFunc("POST /invoices", handler.Create)
adminRouter.HandlerFunc("DELETE /invoices", handler.Delete)

// Subroute through a middlewares that ensure user is admin
router.Handle("/", middlewares.EnsureAdmin(adminRouter))

func EnsureAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Checking if user is admin")
		if !strings.Contains(r.Header.Get("Authorization"), "Admin") {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
			return
		}
		next.ServeHTTP(w, r)
	})
}

server := http.Server {
	Addr: ":8080",
	Handler: router,
}
*/

/*
Passes through context

const AuthUserID = "middlewares.auth.userID"

func writeUnauthed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
}

func IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")

		// Check that the header begins with a prefix of Bearer
		if !strings.HasPrefix(authorization, "Bearer ") {
			writeUnauthed(w)
			return
		}

		// Pull out the token
		encodedToken := strings.TrimPrefix(authorization, "Bearer ")

		// Decode the token from base 64
		token, err := base64.StdEncoding.DecodeString(encodedToken)
		if err != nil {
			writeUnauthed(w)
			return
		}

		// We're just assuming a valid base64 token is a valid user id.
		userID := string(token)
		fmt.Println("userID:", userID)

		next.ServeHTTP(w, r)
	})
}

// We gonna be making the user id available to any downstream handlers
// using context.Context using key value
// Every http request have assosiated context
// Which we can extend and override

// 1. Define unique key to both set and get the user id from the request context
const AuthUserId = "middlewares.auth.userId"

// 2. Replace
func IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")

		if !strings.HasPrefix(authorization, "Bearer ") {
			writeUnauthed(w)
			return
		}

		encodedToken := strings.TrimPrefix(authorization, "Bearer ")

		token, err := base64.StdEncoding.DecodeString(encodedToken)
		if err != nil {
			writeUnauthed(w)
			return
		}

		userID := string(token)

		// This will return a child context, which contain the key value pair
		ctx := context.WithValue(r.Context(), AuthUserId, userId)

		// Passing the new context
		// return a new copy of the request with additional context
		req := r.WithContext(ctx)

		// Pass it
		next.ServeHTTP(w, req)
	})
}

// Pull the userId
type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middlewares.AuthUserId).(string) // Add typecast
	if !ok {
		log.Println("invalid userId")
		w.WriteHeader(http.StatusBadRequest)
	}
}
*/
