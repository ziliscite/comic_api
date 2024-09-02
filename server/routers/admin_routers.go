package router

import (
	handler "bookstore/server/handlers"
	"bookstore/utils/helpers"
	"net/http"
)

func AdminRouters(muxHandler *handler.Handler) *http.ServeMux {
	adminRouter := http.NewServeMux()

	adminRouter.Handle("POST /comics", helpers.ServeHandler(muxHandler.CreateComic))
	adminRouter.Handle("POST /comics/{comic_slug}/{genre_name}", helpers.ServeHandler(muxHandler.AddGenreToComic))
	adminRouter.Handle("POST /comics/{comic_slug}", helpers.ServeHandler(muxHandler.CreateChapter))

	return adminRouter
}
