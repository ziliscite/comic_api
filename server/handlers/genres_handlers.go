package handler

import (
	"bookstore/database"
	"bookstore/utils/helpers"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (h *Handler) GetGenres(w http.ResponseWriter, r *http.Request) (int, error) {
	comics, err := h.Queries.GetGenres(h.Context)
	if err != nil {
		h.Middlewares.Printf("Error getting genres: %s", err)
		return http.StatusInternalServerError, err
	}

	helpers.RespondWithJSON(w, http.StatusOK, comics)
	return http.StatusOK, nil
}

func (h *Handler) CreateGenres(w http.ResponseWriter, r *http.Request) (int, error) {
	type AddGenreParams struct {
		GenreName string `json:"genre_name"`
	}

	genreName := AddGenreParams{}
	err := json.NewDecoder(r.Body).Decode(&genreName)
	if err != nil {
		h.Middlewares.Printf("Error parsing request body: %s", err)
		return http.StatusBadRequest, errors.New("invalid request body")
	}

	err = h.Queries.AddGenre(h.Context, genreName.GenreName)
	if err != nil {
		h.Middlewares.Printf("Error adding genre: %s", err)
		return http.StatusInternalServerError, errors.New("error adding genre")
	}

	helpers.RespondWithMessage(w, http.StatusCreated, "Genre has been successfully added")
	return http.StatusCreated, nil
}

func (h *Handler) AddGenreToComic(w http.ResponseWriter, r *http.Request) (int, error) {
	comicSlug := r.PathValue("comic_slug")
	genreName := r.PathValue("genre_name")

	comic, code, err := h.GetComicBySlugHelper(comicSlug)
	if err != nil {
		return code, err
	}

	genre, err := h.Queries.GetGenreByName(h.Context, genreName)
	if err != nil {
		h.Middlewares.Printf("Error getting genre by name: %s", err)
		return http.StatusNotFound, fmt.Errorf("%s genre is not found", genreName)
	}

	AddGenreResp := database.AddGenreToComicParams{
		ComicID: comic.ID,
		GenreID: genre.ID,
	}

	err = h.Queries.AddGenreToComic(h.Context, AddGenreResp)
	if err != nil {
		h.Middlewares.Printf("Error adding genre to comic: %s", err)
		return http.StatusInternalServerError, errors.New("error adding genre to comic")
	}

	message := fmt.Sprintf("%s genre is added to %s", genre.GenreName, comic.Title)
	helpers.RespondWithMessage(w, http.StatusCreated, message)
	return http.StatusCreated, nil
}
