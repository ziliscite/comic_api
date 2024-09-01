package main

import (
	"bookstore/authenticator"
	"bookstore/handlers"
	"bookstore/helpers"
	"bookstore/middlewares"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	// Must import this for the driver to work
	_ "github.com/lib/pq"
)

// 	1. Use SQLC

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	handler := handlers.NewHandler()
	router := http.NewServeMux()

	router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("OK"))
		if err != nil {
			handler.Logger.Printf("Failed health check: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	adminRouter := http.NewServeMux()

	// The creation handler (POST request) should be authorized by either admin or moderator
	adminRouter.Handle("POST /comics", helpers.ServeHandler(handler.CreateComic))
	router.Handle("GET /comics", helpers.ServeHandler(handler.GetComics))
	router.Handle("GET /comics/{comic_slug}", helpers.ServeHandler(handler.GetComicBySlug))

	// Alas, I input all the genres by hand
	router.Handle("GET /genres", helpers.ServeHandler(handler.GetGenres))

	adminRouter.Handle("POST /comics/{comic_slug}/{genre_name}", helpers.ServeHandler(handler.AddGenreToComic))

	adminRouter.Handle("POST /comics/{comic_slug}", helpers.ServeHandler(handler.CreateChapter))
	router.Handle("GET /comics/{comic_slug}/chapters/{chapter_number}", helpers.ServeHandler(handler.GetChapterByNumber))

	router.Handle("POST /register", helpers.ServeHandler(handler.RegisterUser))
	router.Handle("POST /login", helpers.ServeHandler(handler.Login))

	adminMiddleware := middlewares.CreateStack(
		// Yeah, the other way around. This one's correct
		authenticator.AuthenticateMiddleware,
		authenticator.EnsureAdminMiddleware,
	)

	router.Handle("/", adminMiddleware(adminRouter))

	server := http.Server{
		Addr:    handler.ListenAddr,
		Handler: handler.Logger.LoggingMiddleware(router),
	}

	err = server.ListenAndServe()
	if err != nil {
		return
	}
}
