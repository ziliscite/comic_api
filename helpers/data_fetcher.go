package helpers

import (
	"bookstore/database"
	"bookstore/middlewares"
	"context"
	"fmt"
	"github.com/gosimple/slug"
	"net/http"
)

// GetComicBySlugHelper Helps validate the request. Also, it'd be quite a hassle to repeat these lines of code for at least 3 more times.
func GetComicBySlugHelper(ctx context.Context, queries *database.Queries, logger *middlewares.Logger, comicSlug string) (*database.GetComicBySlugRow, int, error) {
	if !slug.IsSlug(comicSlug) {
		logger.Printf("Invalid comic slug: %s", comicSlug)
		return nil, http.StatusBadRequest, fmt.Errorf("invalid comic slug")
	}

	comic, err := queries.GetComicBySlug(ctx, &comicSlug)
	if err != nil {
		logger.Printf("Comic with slug %s is not found: %s", comicSlug, err.Error())
		return nil, http.StatusNotFound, fmt.Errorf("comic is not found")
	}

	return comic, http.StatusOK, nil
}
