package router

import (
	handler "bookstore/server/handlers"
	"bookstore/utils/middlewares"
	"net/http"
)

func GetRouters(muxHandler *handler.Handler) *http.ServeMux {
	router := http.NewServeMux()

	guessRouter := GuessRouter(muxHandler)
	router.Handle("/", guessRouter)

	userRouters := UserRouters(muxHandler)
	userMiddleware := middlewares.CreateStack(
		muxHandler.Middlewares.AuthenticateMiddleware,
	)

	router.Handle("/user/", http.StripPrefix("/user", userMiddleware(userRouters)))

	adminRouter := AdminRouters(muxHandler)
	adminMiddleware := middlewares.CreateStack(
		muxHandler.Middlewares.AuthenticateMiddleware,
		muxHandler.Middlewares.EnsureAdminMiddleware,
	)

	router.Handle("/admin/", http.StripPrefix("/admin", adminMiddleware(adminRouter)))
	return router
}
