// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const addGenreToComic = `-- name: AddGenreToComic :exec
INSERT INTO comic_genres (
    comic_id,
    genre_id
) VALUES (
    $1, $2
)
`

type AddGenreToComicParams struct {
	ComicID int32 `json:"comic_id"`
	GenreID int32 `json:"genre_id"`
}

// Link a comic with a genre
func (q *Queries) AddGenreToComic(ctx context.Context, arg AddGenreToComicParams) error {
	_, err := q.db.Exec(ctx, addGenreToComic, arg.ComicID, arg.GenreID)
	return err
}

const createAuthor = `-- name: CreateAuthor :one

INSERT INTO authors (name) VALUES ($1) RETURNING id, name
`

// query.sql
// Insert a new author
func (q *Queries) CreateAuthor(ctx context.Context, name string) (*Author, error) {
	row := q.db.QueryRow(ctx, createAuthor, name)
	var i Author
	err := row.Scan(&i.ID, &i.Name)
	return &i, err
}

const createChapter = `-- name: CreateChapter :one
INSERT INTO chapters (
    comic_id, chapter_number, title, release_date, content_url
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING id, comic_id, chapter_number, title, content_url, release_date
`

type CreateChapterParams struct {
	ComicID       *int32           `json:"comic_id"`
	ChapterNumber int32            `json:"chapter_number"`
	Title         *string          `json:"title"`
	ReleaseDate   pgtype.Timestamp `json:"release_date"`
	ContentUrl    *string          `json:"content_url"`
}

// Insert a new chapter
func (q *Queries) CreateChapter(ctx context.Context, arg CreateChapterParams) (*Chapter, error) {
	row := q.db.QueryRow(ctx, createChapter,
		arg.ComicID,
		arg.ChapterNumber,
		arg.Title,
		arg.ReleaseDate,
		arg.ContentUrl,
	)
	var i Chapter
	err := row.Scan(
		&i.ID,
		&i.ComicID,
		&i.ChapterNumber,
		&i.Title,
		&i.ContentUrl,
		&i.ReleaseDate,
	)
	return &i, err
}

const createComic = `-- name: CreateComic :one
INSERT INTO comics (
    title, slug, status, summary, release_date, cover_image_url
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING id, slug, title, status, summary, release_date, cover_image_url, upload_date
`

type CreateComicParams struct {
	Title         string      `json:"title"`
	Slug          *string     `json:"slug"`
	Status        StatusType  `json:"status"`
	Summary       *string     `json:"summary"`
	ReleaseDate   pgtype.Date `json:"release_date"`
	CoverImageUrl *string     `json:"cover_image_url"`
}

// Insert a new comic
func (q *Queries) CreateComic(ctx context.Context, arg CreateComicParams) (*Comic, error) {
	row := q.db.QueryRow(ctx, createComic,
		arg.Title,
		arg.Slug,
		arg.Status,
		arg.Summary,
		arg.ReleaseDate,
		arg.CoverImageUrl,
	)
	var i Comic
	err := row.Scan(
		&i.ID,
		&i.Slug,
		&i.Title,
		&i.Status,
		&i.Summary,
		&i.ReleaseDate,
		&i.CoverImageUrl,
		&i.UploadDate,
	)
	return &i, err
}

const deleteComic = `-- name: DeleteComic :exec
DELETE FROM comics
WHERE id = $1
`

// Delete a comic by ID
func (q *Queries) DeleteComic(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteComic, id)
	return err
}

const getAuthors = `-- name: GetAuthors :many
SELECT
    id,
    name
FROM
    authors
ORDER BY
    name
`

// Select all author
func (q *Queries) GetAuthors(ctx context.Context) ([]*Author, error) {
	rows, err := q.db.Query(ctx, getAuthors)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*Author{}
	for rows.Next() {
		var i Author
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAuthorsByComicId = `-- name: GetAuthorsByComicId :many
SELECT
    id,
    name
FROM
    authors
        JOIN
    author_comic ON authors.id = author_comic.author_id
WHERE
    author_comic.comic_id = $1
`

// Select authors by comic id
func (q *Queries) GetAuthorsByComicId(ctx context.Context, comicID int32) ([]*Author, error) {
	rows, err := q.db.Query(ctx, getAuthorsByComicId, comicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*Author{}
	for rows.Next() {
		var i Author
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getChapterByComicSlugAndNumber = `-- name: GetChapterByComicSlugAndNumber :one
SELECT
    ch.id,
    ch.comic_id,
    ch.chapter_number,
    ch.title,
    ch.release_date,
    ch.content_url
FROM
    chapters ch
JOIN
    comics c ON ch.comic_id = c.id
WHERE
    c.slug = $1
AND
    ch.chapter_number = $2
`

type GetChapterByComicSlugAndNumberParams struct {
	Slug          *string `json:"slug"`
	ChapterNumber int32   `json:"chapter_number"`
}

type GetChapterByComicSlugAndNumberRow struct {
	ID            int32            `json:"id"`
	ComicID       *int32           `json:"comic_id"`
	ChapterNumber int32            `json:"chapter_number"`
	Title         *string          `json:"title"`
	ReleaseDate   pgtype.Timestamp `json:"release_date"`
	ContentUrl    *string          `json:"content_url"`
}

// Select a chapter by comic slug and chapter number
func (q *Queries) GetChapterByComicSlugAndNumber(ctx context.Context, arg GetChapterByComicSlugAndNumberParams) (*GetChapterByComicSlugAndNumberRow, error) {
	row := q.db.QueryRow(ctx, getChapterByComicSlugAndNumber, arg.Slug, arg.ChapterNumber)
	var i GetChapterByComicSlugAndNumberRow
	err := row.Scan(
		&i.ID,
		&i.ComicID,
		&i.ChapterNumber,
		&i.Title,
		&i.ReleaseDate,
		&i.ContentUrl,
	)
	return &i, err
}

const getChapterById = `-- name: GetChapterById :one
SELECT
    id,
    comic_id,
    chapter_number,
    title,
    release_date,
    content_url
FROM
    chapters
WHERE
    chapters.id = $1
`

type GetChapterByIdRow struct {
	ID            int32            `json:"id"`
	ComicID       *int32           `json:"comic_id"`
	ChapterNumber int32            `json:"chapter_number"`
	Title         *string          `json:"title"`
	ReleaseDate   pgtype.Timestamp `json:"release_date"`
	ContentUrl    *string          `json:"content_url"`
}

// Select a chapter by id
func (q *Queries) GetChapterById(ctx context.Context, id int32) (*GetChapterByIdRow, error) {
	row := q.db.QueryRow(ctx, getChapterById, id)
	var i GetChapterByIdRow
	err := row.Scan(
		&i.ID,
		&i.ComicID,
		&i.ChapterNumber,
		&i.Title,
		&i.ReleaseDate,
		&i.ContentUrl,
	)
	return &i, err
}

const getChaptersByComicID = `-- name: GetChaptersByComicID :many
SELECT
    content_url,
    chapter_number,
    title,
    release_date
FROM
    chapters
WHERE
    comic_id = $1
ORDER BY
    chapter_number
`

type GetChaptersByComicIDRow struct {
	ContentUrl    *string          `json:"content_url"`
	ChapterNumber int32            `json:"chapter_number"`
	Title         *string          `json:"title"`
	ReleaseDate   pgtype.Timestamp `json:"release_date"`
}

// Select all chapters for a specific comic
func (q *Queries) GetChaptersByComicID(ctx context.Context, comicID *int32) ([]*GetChaptersByComicIDRow, error) {
	rows, err := q.db.Query(ctx, getChaptersByComicID, comicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetChaptersByComicIDRow{}
	for rows.Next() {
		var i GetChaptersByComicIDRow
		if err := rows.Scan(
			&i.ContentUrl,
			&i.ChapterNumber,
			&i.Title,
			&i.ReleaseDate,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getComicBySlug = `-- name: GetComicBySlug :one
SELECT
    id,
    slug,
    title,
    status,
    summary,
    release_date,
    upload_date,
    cover_image_url
FROM
    comics
WHERE
    comics.slug = $1
`

type GetComicBySlugRow struct {
	ID            int32            `json:"id"`
	Slug          *string          `json:"slug"`
	Title         string           `json:"title"`
	Status        StatusType       `json:"status"`
	Summary       *string          `json:"summary"`
	ReleaseDate   pgtype.Date      `json:"release_date"`
	UploadDate    pgtype.Timestamp `json:"upload_date"`
	CoverImageUrl *string          `json:"cover_image_url"`
}

// Select a comic by ID
func (q *Queries) GetComicBySlug(ctx context.Context, slug *string) (*GetComicBySlugRow, error) {
	row := q.db.QueryRow(ctx, getComicBySlug, slug)
	var i GetComicBySlugRow
	err := row.Scan(
		&i.ID,
		&i.Slug,
		&i.Title,
		&i.Status,
		&i.Summary,
		&i.ReleaseDate,
		&i.UploadDate,
		&i.CoverImageUrl,
	)
	return &i, err
}

const getComics = `-- name: GetComics :many
SELECT
    comics.slug,
    comics.title,
    comics.status,
    comics.cover_image_url,
    comics.upload_date
FROM
    comics
ORDER BY
    comics.upload_date
DESC
`

type GetComicsRow struct {
	Slug          *string          `json:"slug"`
	Title         string           `json:"title"`
	Status        StatusType       `json:"status"`
	CoverImageUrl *string          `json:"cover_image_url"`
	UploadDate    pgtype.Timestamp `json:"upload_date"`
}

// Select all comics
func (q *Queries) GetComics(ctx context.Context) ([]*GetComicsRow, error) {
	rows, err := q.db.Query(ctx, getComics)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetComicsRow{}
	for rows.Next() {
		var i GetComicsRow
		if err := rows.Scan(
			&i.Slug,
			&i.Title,
			&i.Status,
			&i.CoverImageUrl,
			&i.UploadDate,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getComicsByGenre = `-- name: GetComicsByGenre :many
SELECT
    comics.slug,
    comics.title,
    comics.status,
    comics.cover_image_url,
    comics.upload_date
FROM
    comics
        JOIN
    comic_genres ON comics.id = comic_genres.comic_id
        JOIN
    genres ON comic_genres.genre_id = genres.id
WHERE
    genres.genre_name = $1
group by comics.id
`

type GetComicsByGenreRow struct {
	Slug          *string          `json:"slug"`
	Title         string           `json:"title"`
	Status        StatusType       `json:"status"`
	CoverImageUrl *string          `json:"cover_image_url"`
	UploadDate    pgtype.Timestamp `json:"upload_date"`
}

// Select comics by genre
func (q *Queries) GetComicsByGenre(ctx context.Context, genreName string) ([]*GetComicsByGenreRow, error) {
	rows, err := q.db.Query(ctx, getComicsByGenre, genreName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetComicsByGenreRow{}
	for rows.Next() {
		var i GetComicsByGenreRow
		if err := rows.Scan(
			&i.Slug,
			&i.Title,
			&i.Status,
			&i.CoverImageUrl,
			&i.UploadDate,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getComicsByTitle = `-- name: GetComicsByTitle :many
SELECT
    comics.slug,
    comics.title,
    comics.status,
    comics.cover_image_url,
    comics.upload_date
FROM
    comics
WHERE
    comics.title ILIKE '%' || $1 || '%'
ORDER BY
    comics.upload_date DESC
`

type GetComicsByTitleRow struct {
	Slug          *string          `json:"slug"`
	Title         string           `json:"title"`
	Status        StatusType       `json:"status"`
	CoverImageUrl *string          `json:"cover_image_url"`
	UploadDate    pgtype.Timestamp `json:"upload_date"`
}

// Select comics by title
func (q *Queries) GetComicsByTitle(ctx context.Context, dollar_1 *string) ([]*GetComicsByTitleRow, error) {
	rows, err := q.db.Query(ctx, getComicsByTitle, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetComicsByTitleRow{}
	for rows.Next() {
		var i GetComicsByTitleRow
		if err := rows.Scan(
			&i.Slug,
			&i.Title,
			&i.Status,
			&i.CoverImageUrl,
			&i.UploadDate,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getGenreByName = `-- name: GetGenreByName :one
SELECT
    id,
    genre_name
FROM
    genres
WHERE
    genres.genre_name = $1
`

// Select a genre by name
func (q *Queries) GetGenreByName(ctx context.Context, genreName string) (*Genre, error) {
	row := q.db.QueryRow(ctx, getGenreByName, genreName)
	var i Genre
	err := row.Scan(&i.ID, &i.GenreName)
	return &i, err
}

const getGenres = `-- name: GetGenres :many
SELECT
    id,
    genre_name
FROM
    genres
ORDER BY
    genre_name
`

// Select all genres
func (q *Queries) GetGenres(ctx context.Context) ([]*Genre, error) {
	rows, err := q.db.Query(ctx, getGenres)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*Genre{}
	for rows.Next() {
		var i Genre
		if err := rows.Scan(&i.ID, &i.GenreName); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getGenresByComicId = `-- name: GetGenresByComicId :many
SELECT
    genres.id,
    genres.genre_name
FROM
    genres
        JOIN
    comic_genres ON genres.id = comic_genres.genre_id
WHERE
    comic_genres.comic_id = $1
`

// Select genres by comic id
func (q *Queries) GetGenresByComicId(ctx context.Context, comicID int32) ([]*Genre, error) {
	rows, err := q.db.Query(ctx, getGenresByComicId, comicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*Genre{}
	for rows.Next() {
		var i Genre
		if err := rows.Scan(&i.ID, &i.GenreName); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const removeGenreFromComic = `-- name: RemoveGenreFromComic :exec
DELETE FROM comic_genres
WHERE comic_id = $1 AND genre_id = $2
`

type RemoveGenreFromComicParams struct {
	ComicID int32 `json:"comic_id"`
	GenreID int32 `json:"genre_id"`
}

// Remove a genre from a comic
func (q *Queries) RemoveGenreFromComic(ctx context.Context, arg RemoveGenreFromComicParams) error {
	_, err := q.db.Exec(ctx, removeGenreFromComic, arg.ComicID, arg.GenreID)
	return err
}

const updateComicStatus = `-- name: UpdateComicStatus :exec
UPDATE comics
SET status = $2
WHERE id = $1
`

type UpdateComicStatusParams struct {
	ID     int32      `json:"id"`
	Status StatusType `json:"status"`
}

// Update comic status
func (q *Queries) UpdateComicStatus(ctx context.Context, arg UpdateComicStatusParams) error {
	_, err := q.db.Exec(ctx, updateComicStatus, arg.ID, arg.Status)
	return err
}
