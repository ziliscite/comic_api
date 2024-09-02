package handler

import (
	"bookstore/database"
	"bookstore/token_maker"
	"bookstore/utils/helpers"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func (h *Handler) AddComicBookmark(w http.ResponseWriter, r *http.Request) (int, error) {
	// From the AuthenticateMiddleware btw
	claims, ok := r.Context().Value(token_maker.ClaimsKey).(*token_maker.CustomClaims)
	if !ok {
		return http.StatusUnauthorized, errors.New("not authorized")
	}

	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return http.StatusBadRequest, errors.New("invalid user id")
	}

	comicSlug := r.PathValue("comic_slug")

	comic, code, err := h.GetComicBySlugHelper(comicSlug)
	if err != nil {
		return code, err
	}

	AddBookmarkResp := database.BookmarkComicParams{
		UserID:  int32(userId),
		ComicID: comic.ID,
	}

	err = h.Queries.BookmarkComic(h.Context, AddBookmarkResp)
	if err != nil {
		h.Middlewares.Printf("Error bookmarking comic: %s", err)
		return http.StatusInternalServerError, errors.New("error bookmarking comic")
	}

	message := fmt.Sprintf("%s comic is added to bookmark", comic.Title)
	helpers.RespondWithMessage(w, http.StatusCreated, message)
	return http.StatusCreated, nil
}

func (h *Handler) RemoveComicBookmark(w http.ResponseWriter, r *http.Request) (int, error) {
	claims, ok := r.Context().Value(token_maker.ClaimsKey).(*token_maker.CustomClaims)
	if !ok {
		return http.StatusUnauthorized, errors.New("not authorized")
	}

	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return http.StatusBadRequest, errors.New("invalid user id")
	}

	comicSlug := r.PathValue("comic_slug")
	comic, code, err := h.GetComicBySlugHelper(comicSlug)
	if err != nil {
		return code, err
	}

	RemoveBookmarkResp := database.RemoveComicFromBookmarkParams{
		UserID:  int32(userId),
		ComicID: comic.ID,
	}

	err = h.Queries.RemoveComicFromBookmark(h.Context, RemoveBookmarkResp)
	if err != nil {
		h.Middlewares.Printf("Error removing comic from bookmark: %s", err)
		return http.StatusInternalServerError, errors.New("error removing comic from bookmark")
	}

	// Or we use 204 without any json response ?
	message := fmt.Sprintf("%s comic is removed from bookmark", comic.Title)
	helpers.RespondWithMessage(w, http.StatusOK, message)
	return http.StatusOK, nil
}
