package handler

import (
	"bookstore/database"
	"bookstore/utils/middlewares"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
)

type Handler struct {
	Queries     *database.Queries
	Middlewares *middlewares.Middleware
	Context     context.Context
	ListenAddr  string
	JWTSecret   string
}

func NewHandler() *Handler {
	connStr := fmt.Sprintf("postgres://postgres:%s@localhost:5432/comic?sslmode=disable", os.Getenv("POSTGRESQL_SECRET"))

	ctx := context.Background()

	db, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}

	queries := database.New(db)

	return &Handler{
		Queries:     queries,
		Middlewares: middlewares.NewMiddleware(),
		Context:     ctx,
		ListenAddr:  os.Getenv("LISTEN_ADDR"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
	}
}
