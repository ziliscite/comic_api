package main

import (
	"bookstore/handlers"
	"bookstore/helpers"
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

	router.Handle("POST /comics", helpers.ServeHandler(handler.CreateComic))
	router.Handle("GET /comics", helpers.ServeHandler(handler.GetComics))
	router.Handle("GET /comics/{comic_slug}", helpers.ServeHandler(handler.GetComicBySlug))

	// Alas, I input all the genres by hand
	router.Handle("GET /genres", helpers.ServeHandler(handler.GetGenres))

	router.Handle("POST /comics/{comic_slug}/{genre_name}", helpers.ServeHandler(handler.AddGenreToComic))

	router.Handle("POST /comics/{comic_slug}", helpers.ServeHandler(handler.CreateChapter))
	router.Handle("GET /comics/{comic_slug}/chapters/{chapter_number}", helpers.ServeHandler(handler.GetChapterByNumber))

	server := http.Server{
		Addr:    handler.ListenAddr,
		Handler: handler.Logger.LoggingMiddleware(router),
	}

	err = server.ListenAndServe()
	if err != nil {
		return
	}
}
