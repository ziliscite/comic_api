package handlers

import (
	"bookstore/database"
	"bookstore/helpers"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gosimple/slug"
	"net/http"
)

func (h *Handler) GetComics(w http.ResponseWriter, r *http.Request) (int, error) {
	ctx := context.Background()

	comics, err := h.Queries.GetComics(ctx)
	if err != nil {
		h.Logger.Printf("Error getting comics: %s", err)
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

	ctx := context.Background()

	comic, code, err := helpers.GetComicBySlugHelper(ctx, h.Queries, h.Logger, comicSlug)
	if err != nil {
		return code, err
	}

	// I figured that
	// I guess it's okay if it's not found? Just make it an empty list or something
	genres, _ := h.Queries.GetGenresByComicId(ctx, comic.ID)
	chapters, _ := h.Queries.GetChaptersByComicID(ctx, &comic.ID)

	comicResp := ComicData{
		Comic:    comic,
		Genres:   genres,
		Chapters: chapters,
	}

	helpers.RespondWithJSON(w, code, comicResp)
	return code, nil
}

func (h *Handler) CreateComic(w http.ResponseWriter, r *http.Request) (int, error) {
	ctx := context.Background()

	comicParams := database.CreateComicParams{}
	err := json.NewDecoder(r.Body).Decode(&comicParams)
	if err != nil {
		h.Logger.Printf("Error parsing request body: %s", err)
		return http.StatusBadRequest, fmt.Errorf("invalid request body")
	}

	comicSlug := slug.Make(comicParams.Title)
	comicSlugPtr := &comicSlug

	// Such length...
	comicParams.Slug = comicSlugPtr

	comic, err := h.Queries.CreateComic(ctx, comicParams)
	if err != nil {
		h.Logger.Printf("Error creating comic: %s", err)
		return http.StatusBadRequest, fmt.Errorf("error creating comic")
	}

	helpers.RespondWithJSON(w, http.StatusCreated, comic)
	return http.StatusCreated, nil
}
