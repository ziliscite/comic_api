package handler

import (
	"bookstore/database"
	"bookstore/utils/helpers"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (h *Handler) CreateAuthors(w http.ResponseWriter, r *http.Request) (int, error) {
	type AddAuthorParams struct {
		AuthorName string `json:"author_name"`
	}

	authorName := AddAuthorParams{}
	err := json.NewDecoder(r.Body).Decode(&authorName)
	if err != nil {
		h.Middlewares.Printf("Error parsing request body: %s", err)
		return http.StatusBadRequest, errors.New("invalid request body")
	}

	author, err := h.Queries.CreateAuthor(h.Context, authorName.AuthorName)
	if err != nil {
		h.Middlewares.Printf("Error adding genre: %s", err)
		return http.StatusInternalServerError, errors.New("error adding genre")
	}

	helpers.RespondWithMessage(w, http.StatusCreated, fmt.Sprintf("%s has been successfully added", author.Name))
	return http.StatusCreated, nil
}

func (h *Handler) GetAuthors(w http.ResponseWriter, r *http.Request) (int, error) {
	authors, err := h.Queries.GetAuthors(h.Context)
	if err != nil {
		h.Middlewares.Printf("Error getting authors: %s", err)
		return http.StatusInternalServerError, errors.New("error getting authors")
	}

	helpers.RespondWithJSON(w, http.StatusOK, authors)
	return http.StatusOK, nil
}

func (h *Handler) AddComicAuthor(w http.ResponseWriter, r *http.Request) (int, error) {
	comicSlug := r.PathValue("comic_slug")
	authorName := r.PathValue("author_name")

	comic, code, err := h.GetComicBySlugHelper(comicSlug)
	if err != nil {
		return code, err
	}

	author, err := h.Queries.GetAuthorByName(h.Context, authorName)
	if err != nil {
		h.Middlewares.Printf("Error getting author by name: %s", err)
		return http.StatusNotFound, fmt.Errorf("%s author is not found", authorName)
	}

	AddAuthorResp := database.AuthorComicParams{
		AuthorID: author.ID,
		ComicID:  comic.ID,
	}

	err = h.Queries.AuthorComic(h.Context, AddAuthorResp)
	if err != nil {
		h.Middlewares.Printf("Error adding author's comic: %s", err)
		return http.StatusInternalServerError, errors.New("error adding author's comic")
	}

	message := fmt.Sprintf("%s's comic %s has been added", author.Name, comic.Title)
	helpers.RespondWithMessage(w, http.StatusCreated, message)
	return http.StatusCreated, nil
}
