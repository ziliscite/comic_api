package handlers

import (
	"bookstore/database"
	"bookstore/helpers"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (h *Handler) CreateChapter(w http.ResponseWriter, r *http.Request) (int, error) {
	comicSlug := r.PathValue("comic_slug")

	ctx := context.Background()

	comic, code, err := helpers.GetComicBySlugHelper(ctx, h.Queries, h.Logger, comicSlug)
	if err != nil {
		return code, err
	}

	chapterParams := database.CreateChapterParams{
		ComicID: &comic.ID,
	}

	err = json.NewDecoder(r.Body).Decode(&chapterParams)
	if err != nil {
		h.Logger.Printf("Error parsing request body: %s", err)
		return http.StatusBadRequest, fmt.Errorf("invalid request body")
	}

	chapter, err := h.Queries.CreateChapter(ctx, chapterParams)
	if err != nil {
		h.Logger.Printf("Error creating chapter: %s", err)
		return http.StatusBadRequest, fmt.Errorf("error creating chapter")
	}

	helpers.RespondWithJSON(w, http.StatusCreated, chapter)
	return http.StatusCreated, nil
}

func (h *Handler) GetChapterByNumber(w http.ResponseWriter, r *http.Request) (int, error) {
	comicSlug := r.PathValue("comic_slug")
	chapterNumber, err := strconv.Atoi(r.PathValue("chapter_number"))
	if err != nil {
		h.Logger.Printf("Error parsing request body: %s", err)
		return http.StatusBadRequest, fmt.Errorf("invalid request body")
	}

	ctx := context.Background()

	chapterReq := database.GetChapterByComicSlugAndNumberParams{
		Slug:          &comicSlug,
		ChapterNumber: int32(chapterNumber),
	}

	chapterResp, err := h.Queries.GetChapterByComicSlugAndNumber(ctx, chapterReq)
	if err != nil {
		h.Logger.Printf("Error getting chapter: %s", err)
		return http.StatusNotFound, fmt.Errorf("chapter is not found")
	}

	helpers.RespondWithJSON(w, http.StatusOK, chapterResp)
	return http.StatusOK, nil
}
