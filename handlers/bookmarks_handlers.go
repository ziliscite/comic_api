package handlers

import (
	"bookstore/database"
	"bookstore/helpers"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func (h *Handler) AddComicBookmark(w http.ResponseWriter, r *http.Request) (int, error) {
	// From the AuthenticateMiddleware btw
	claims, ok := r.Context().Value(ClaimsKey).(*CustomClaims)
	if !ok {
		return http.StatusUnauthorized, errors.New("not authorized")
	}

	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return http.StatusBadRequest, errors.New("invalid user id")
	}

	comicSlug := r.PathValue("comic_slug")

	ctx := context.Background()

	comic, code, err := helpers.GetComicBySlugHelper(ctx, h.Queries, h.Logger, comicSlug)
	if err != nil {
		return code, err
	}

	AddBookmarkResp := database.BookmarkComicParams{
		UserID:  int32(userId),
		ComicID: comic.ID,
	}

	err = h.Queries.BookmarkComic(ctx, AddBookmarkResp)
	if err != nil {
		h.Logger.Printf("Error bookmarking comic: %s", err)
		return http.StatusInternalServerError, errors.New("error bookmarking comic")
	}

	message := fmt.Sprintf("%s comic is added to bookmark", comic.Title)
	helpers.RespondWithMessage(w, http.StatusCreated, message)
	return http.StatusCreated, nil
}
