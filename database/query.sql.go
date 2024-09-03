// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const addGenre = `-- name: AddGenre :exec
INSERT INTO genres (
    genre_name
) VALUES (
    $1
) RETURNING id, genre_name
`

// Add a new genre
func (q *Queries) AddGenre(ctx context.Context, genreName string) error {
	_, err := q.db.Exec(ctx, addGenre, genreName)
	return err
}

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

const addSession = `-- name: AddSession :one
INSERT INTO sessions (
    user_id, session_token, created_at, expires_at, is_active
) VALUES (
    $1, $2, now(), $3, TRUE
) RETURNING session_id, user_id, session_token, created_at, expires_at, is_active
`

type AddSessionParams struct {
	UserID       *int32           `json:"user_id"`
	SessionToken string           `json:"session_token"`
	ExpiresAt    pgtype.Timestamp `json:"expires_at"`
}

// Add a session
func (q *Queries) AddSession(ctx context.Context, arg AddSessionParams) (*Session, error) {
	row := q.db.QueryRow(ctx, addSession, arg.UserID, arg.SessionToken, arg.ExpiresAt)
	var i Session
	err := row.Scan(
		&i.SessionID,
		&i.UserID,
		&i.SessionToken,
		&i.CreatedAt,
		&i.ExpiresAt,
		&i.IsActive,
	)
	return &i, err
}

const authorComic = `-- name: AuthorComic :exec
INSERT INTO author_comic (
    author_id,
    comic_id
) VALUES (
    $1, $2
)
`

type AuthorComicParams struct {
	AuthorID int32 `json:"author_id"`
	ComicID  int32 `json:"comic_id"`
}

// Insert an author to a comic
func (q *Queries) AuthorComic(ctx context.Context, arg AuthorComicParams) error {
	_, err := q.db.Exec(ctx, authorComic, arg.AuthorID, arg.ComicID)
	return err
}

const bookmarkComic = `-- name: BookmarkComic :exec
INSERT INTO bookmark (
    user_id,
    comic_id
) VALUES (
     $1, $2
)
`

type BookmarkComicParams struct {
	UserID  int32 `json:"user_id"`
	ComicID int32 `json:"comic_id"`
}

// Bookmark a comic
func (q *Queries) BookmarkComic(ctx context.Context, arg BookmarkComicParams) error {
	_, err := q.db.Exec(ctx, bookmarkComic, arg.UserID, arg.ComicID)
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

const getAuthorByName = `-- name: GetAuthorByName :one
SELECT
    id,
    name
FROM
    authors
WHERE
    name = $1
`

// Select author by name
func (q *Queries) GetAuthorByName(ctx context.Context, name string) (*Author, error) {
	row := q.db.QueryRow(ctx, getAuthorByName, name)
	var i Author
	err := row.Scan(&i.ID, &i.Name)
	return &i, err
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

const getSessionFromToken = `-- name: GetSessionFromToken :one
SELECT
    session_id, user_id, session_token, created_at, expires_at, is_active
FROM
    sessions
WHERE
    session_token = $1
`

// Get a session using token
func (q *Queries) GetSessionFromToken(ctx context.Context, sessionToken string) (*Session, error) {
	row := q.db.QueryRow(ctx, getSessionFromToken, sessionToken)
	var i Session
	err := row.Scan(
		&i.SessionID,
		&i.UserID,
		&i.SessionToken,
		&i.CreatedAt,
		&i.ExpiresAt,
		&i.IsActive,
	)
	return &i, err
}

const getSessionFromUserId = `-- name: GetSessionFromUserId :one
SELECT
    session_id, user_id, session_token, created_at, expires_at, is_active
FROM
    sessions
WHERE
    user_id = $1
`

// Get a session using user id
func (q *Queries) GetSessionFromUserId(ctx context.Context, userID *int32) (*Session, error) {
	row := q.db.QueryRow(ctx, getSessionFromUserId, userID)
	var i Session
	err := row.Scan(
		&i.SessionID,
		&i.UserID,
		&i.SessionToken,
		&i.CreatedAt,
		&i.ExpiresAt,
		&i.IsActive,
	)
	return &i, err
}

const getUserRole = `-- name: GetUserRole :one
SELECT
    user_id, username, email, role
FROM
    users
WHERE
    users.user_id = $1
`

type GetUserRoleRow struct {
	UserID   int32    `json:"user_id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Role     UserRole `json:"role"`
}

// Get User Role
func (q *Queries) GetUserRole(ctx context.Context, userID int32) (*GetUserRoleRow, error) {
	row := q.db.QueryRow(ctx, getUserRole, userID)
	var i GetUserRoleRow
	err := row.Scan(
		&i.UserID,
		&i.Username,
		&i.Email,
		&i.Role,
	)
	return &i, err
}

const loginWithEmail = `-- name: LoginWithEmail :one
SELECT
    user_id, username, email, password, first_name, last_name, date_of_birth, role
FROM
    users
WHERE
    users.email = $1
`

type LoginWithEmailRow struct {
	UserID      int32       `json:"user_id"`
	Username    string      `json:"username"`
	Email       string      `json:"email"`
	Password    string      `json:"password"`
	FirstName   *string     `json:"first_name"`
	LastName    *string     `json:"last_name"`
	DateOfBirth pgtype.Date `json:"date_of_birth"`
	Role        UserRole    `json:"role"`
}

// Login User
func (q *Queries) LoginWithEmail(ctx context.Context, email string) (*LoginWithEmailRow, error) {
	row := q.db.QueryRow(ctx, loginWithEmail, email)
	var i LoginWithEmailRow
	err := row.Scan(
		&i.UserID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.FirstName,
		&i.LastName,
		&i.DateOfBirth,
		&i.Role,
	)
	return &i, err
}

const registerUser = `-- name: RegisterUser :one
INSERT INTO users (
    username, email, password, first_name, last_name, date_of_birth, role, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, now(), now()
) RETURNING user_id, username, email, password, first_name, last_name, date_of_birth, role, created_at, updated_at
`

type RegisterUserParams struct {
	Username    string      `json:"username"`
	Email       string      `json:"email"`
	Password    string      `json:"password"`
	FirstName   *string     `json:"first_name"`
	LastName    *string     `json:"last_name"`
	DateOfBirth pgtype.Date `json:"date_of_birth"`
	Role        UserRole    `json:"role"`
}

// Register User
func (q *Queries) RegisterUser(ctx context.Context, arg RegisterUserParams) (*User, error) {
	row := q.db.QueryRow(ctx, registerUser,
		arg.Username,
		arg.Email,
		arg.Password,
		arg.FirstName,
		arg.LastName,
		arg.DateOfBirth,
		arg.Role,
	)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.FirstName,
		&i.LastName,
		&i.DateOfBirth,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const removeComicFromBookmark = `-- name: RemoveComicFromBookmark :exec
DELETE FROM bookmark
WHERE user_id = $1 AND comic_id = $2
`

type RemoveComicFromBookmarkParams struct {
	UserID  int32 `json:"user_id"`
	ComicID int32 `json:"comic_id"`
}

// Remove a comic from user's bookmark
func (q *Queries) RemoveComicFromBookmark(ctx context.Context, arg RemoveComicFromBookmarkParams) error {
	_, err := q.db.Exec(ctx, removeComicFromBookmark, arg.UserID, arg.ComicID)
	return err
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

const revokeSession = `-- name: RevokeSession :exec
UPDATE sessions
SET is_active = FALSE
WHERE session_id = $1
`

// Revoke a session
func (q *Queries) RevokeSession(ctx context.Context, sessionID int32) error {
	_, err := q.db.Exec(ctx, revokeSession, sessionID)
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
