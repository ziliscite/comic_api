package main

import (
	"bookstore/server/handlers"
	router "bookstore/server/routers"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	// Must import this for the driver to work
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = runServer()
	if err != nil {
		log.Fatalf("Something went wrong: %s", err.Error())
	}
}

func runServer() error {
	muxHandler := handler.NewHandler()

	routers := router.GetRouters(muxHandler)

	muxServer := http.Server{
		Addr:    muxHandler.ListenAddr,
		Handler: muxHandler.Middlewares.LoggingMiddleware(routers),
	}

	err := muxServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
