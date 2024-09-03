package router

import (
	handler "bookstore/server/handlers"
	"bookstore/utils/helpers"
	"net/http"
)

func GuessRouter(muxHandler *handler.Handler) *http.ServeMux {
	guessRouter := http.NewServeMux()

	guessRouter.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			muxHandler.Middlewares.Printf("Failed health check: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	guessRouter.Handle("GET /comics", helpers.ServeHandler(muxHandler.GetComics))

	// This will also return all of its genres & chapters
	guessRouter.Handle("GET /comics/{comic_slug}", helpers.ServeHandler(muxHandler.GetComicBySlug))

	guessRouter.Handle("GET /genres", helpers.ServeHandler(muxHandler.GetGenres))

	guessRouter.Handle("GET /comics/{comic_slug}/chapters/{chapter_number}", helpers.ServeHandler(muxHandler.GetChapterByNumber))

	guessRouter.Handle("GET /authors", helpers.ServeHandler(muxHandler.GetAuthors))

	guessRouter.Handle("POST /register", helpers.ServeHandler(muxHandler.Register))
	guessRouter.Handle("POST /login", helpers.ServeHandler(muxHandler.Login))

	// I guess technically we have to be verified or something
	guessRouter.Handle("POST /refresh", helpers.ServeHandler(muxHandler.Refresh))
	guessRouter.Handle("POST /revoke", helpers.ServeHandler(muxHandler.Revoke))

	return guessRouter
}
