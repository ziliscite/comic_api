package router

import (
	server "bookstore/server/handlers"
	"bookstore/utils/helpers"
	"net/http"
)

func UserRouters(h *server.Handler) http.Handler {
	userRouter := http.NewServeMux()

	userRouter.Handle("POST /bookmark/{comic_slug}", helpers.ServeHandler(h.AddComicBookmark))
	userRouter.Handle("DELETE /bookmark/{comic_slug}", helpers.ServeHandler(h.RemoveComicBookmark))

	return userRouter
}
