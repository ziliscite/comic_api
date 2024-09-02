package handler

import (
	"bookstore/database"
	"bookstore/utils/helpers"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) CreateChapter(w http.ResponseWriter, r *http.Request) (int, error) {
	comicSlug := r.PathValue("comic_slug")

	comic, code, err := h.GetComicBySlugHelper(comicSlug)
	if err != nil {
		return code, err
	}

	chapterParams := database.CreateChapterParams{
		ComicID: &comic.ID,
	}

	err = json.NewDecoder(r.Body).Decode(&chapterParams)
	if err != nil {
		h.Middlewares.Printf("Error parsing request body: %s", err)
		return http.StatusBadRequest, errors.New("invalid request body")
	}

	if chapterParams.Title == nil {
		chapterParams.Title = new(string)
		*chapterParams.Title = fmt.Sprintf("Chapter %d", chapterParams.ChapterNumber)
	}

	if !chapterParams.ReleaseDate.Valid {
		chapterParams.ReleaseDate = pgtype.Timestamp{
			Valid: true,
			Time:  time.Now(),
		}
	}

	chapter, err := h.Queries.CreateChapter(h.Context, chapterParams)
	if err != nil {
		h.Middlewares.Printf("Error creating chapter: %s", err)
		return http.StatusBadRequest, errors.New("error creating chapter")
	}

	helpers.RespondWithJSON(w, http.StatusCreated, chapter)
	return http.StatusCreated, nil
}

func (h *Handler) GetChapterByNumber(w http.ResponseWriter, r *http.Request) (int, error) {
	comicSlug := r.PathValue("comic_slug")
	chapterNumber, err := strconv.Atoi(r.PathValue("chapter_number"))
	if err != nil {
		h.Middlewares.Printf("Error parsing request body: %s", err)
		return http.StatusBadRequest, errors.New("invalid request body")
	}

	chapterReq := database.GetChapterByComicSlugAndNumberParams{
		Slug:          &comicSlug,
		ChapterNumber: int32(chapterNumber),
	}

	chapterResp, err := h.Queries.GetChapterByComicSlugAndNumber(h.Context, chapterReq)
	if err != nil {
		h.Middlewares.Printf("Error getting chapter: %s", err)
		return http.StatusNotFound, errors.New("chapter is not found")
	}

	helpers.RespondWithJSON(w, http.StatusOK, chapterResp)
	return http.StatusOK, nil
}
