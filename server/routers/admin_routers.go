package router

import (
	handler "bookstore/server/handlers"
	"bookstore/utils/helpers"
	"net/http"
)

func AdminRouters(muxHandler *handler.Handler) *http.ServeMux {
	adminRouter := http.NewServeMux()

	adminRouter.Handle("POST /comics", helpers.ServeHandler(muxHandler.CreateComic))
	adminRouter.Handle("POST /comics/{comic_slug}", helpers.ServeHandler(muxHandler.CreateChapter))

	adminRouter.Handle("POST /genres", helpers.ServeHandler(muxHandler.CreateGenres))
	adminRouter.Handle("POST /genres/{genre_name}/{comic_slug}", helpers.ServeHandler(muxHandler.AddGenreToComic))

	adminRouter.Handle("POST /authors", helpers.ServeHandler(muxHandler.CreateAuthors))
	adminRouter.Handle("POST /authors/{author_name}/{comic_slug}", helpers.ServeHandler(muxHandler.AddComicAuthor))

	return adminRouter
}
