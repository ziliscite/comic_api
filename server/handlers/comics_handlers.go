package handler

import (
	"bookstore/database"
	"bookstore/utils/helpers"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gosimple/slug"
	"net/http"
)

func (h *Handler) GetComics(w http.ResponseWriter, r *http.Request) (int, error) {
	comics, err := h.Queries.GetComics(h.Context)
	if err != nil {
		h.Middlewares.Printf("Error getting comics: %s", err)
		return http.StatusInternalServerError, err
	}

	helpers.RespondWithJSON(w, http.StatusOK, comics)
	return http.StatusOK, nil
}

func (h *Handler) GetComicBySlug(w http.ResponseWriter, r *http.Request) (int, error) {
	type ComicData struct {
		Comic    *database.GetComicBySlugRow         `json:"comic"`
		Genres   []*database.Genre                   `json:"genres"`
		Chapters []*database.GetChaptersByComicIDRow `json:"chapters"`
	}

	comicSlug := r.PathValue("comic_slug")

	comic, code, err := h.GetComicBySlugHelper(comicSlug)
	if err != nil {
		return code, err
	}

	genres, _ := h.Queries.GetGenresByComicId(h.Context, comic.ID)
	chapters, _ := h.Queries.GetChaptersByComicID(h.Context, &comic.ID)

	comicResp := ComicData{
		Comic:    comic,
		Genres:   genres,
		Chapters: chapters,
	}

	helpers.RespondWithJSON(w, code, comicResp)
	return code, nil
}

func (h *Handler) CreateComic(w http.ResponseWriter, r *http.Request) (int, error) {
	comicParams := database.CreateComicParams{}
	err := json.NewDecoder(r.Body).Decode(&comicParams)
	if err != nil {
		h.Middlewares.Printf("Error parsing request body: %s", err)
		return http.StatusBadRequest, errors.New("invalid request body")
	}

	comicSlug := slug.Make(comicParams.Title)
	comicSlugPtr := &comicSlug

	// Such length...
	comicParams.Slug = comicSlugPtr

	comic, err := h.Queries.CreateComic(h.Context, comicParams)
	if err != nil {
		h.Middlewares.Printf("Error creating comic: %s", err)
		return http.StatusBadRequest, errors.New("error creating comic")
	}

	helpers.RespondWithJSON(w, http.StatusCreated, comic)
	return http.StatusCreated, nil
}

func (h *Handler) GetComicBySlugHelper(comicSlug string) (*database.GetComicBySlugRow, int, error) {
	if !slug.IsSlug(comicSlug) {
		h.Middlewares.Printf("Invalid comic slug: %s", comicSlug)
		return nil, http.StatusBadRequest, fmt.Errorf("invalid comic slug")
	}

	comic, err := h.Queries.GetComicBySlug(h.Context, &comicSlug)
	if err != nil {
		h.Middlewares.Printf("Comic with slug %s is not found: %s", comicSlug, err.Error())
		return nil, http.StatusNotFound, fmt.Errorf("comic is not found")
	}

	return comic, http.StatusOK, nil
}
