package handlers

import (
	"bookstore/database"
	"bookstore/middlewares"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
)

type Handler struct {
	Queries    *database.Queries
	Logger     *middlewares.Logger
	ListenAddr string
}

func NewHandler() *Handler {
	connStr := fmt.Sprintf("postgres://postgres:%s@localhost:5432/comic?sslmode=disable", os.Getenv("POSTGRESQL_SECRET"))
	ctx := context.Background()
	db, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatal(err)
	}

	// I just realized this ctx thing is probably for timeout or something
	// Like, taking too long, return error
	err = db.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to database")

	queries := database.New(db)
	logger := &middlewares.Logger{Log: log.New(os.Stdout, "", log.LstdFlags)}

	return &Handler{
		Queries:    queries,
		Logger:     logger,
		ListenAddr: os.Getenv("LISTEN_ADDR"),
	}
}
