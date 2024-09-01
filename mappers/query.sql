-- query.sql

-- Insert a new author
-- name: CreateAuthor :one
INSERT INTO authors (name) VALUES ($1) RETURNING *;

-- Select all author
-- name: GetAuthors :many
SELECT
    id,
    name
FROM
    authors
ORDER BY
    name;

-- Select authors by comic id
-- name: GetAuthorsByComicId :many
SELECT
    id,
    name
FROM
    authors
        JOIN
    author_comic ON authors.id = author_comic.author_id
WHERE
    author_comic.comic_id = $1;

-- Select all comics
-- name: GetComics :many
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
DESC;

-- Select a comic by ID
-- name: GetComicBySlug :one
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
    comics.slug = $1;

-- Insert a new comic
-- name: CreateComic :one
INSERT INTO comics (
    title, slug, status, summary, release_date, cover_image_url
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- Update comic status
-- name: UpdateComicStatus :exec
UPDATE comics
SET status = $2
WHERE id = $1;

-- Delete a comic by ID
-- name: DeleteComic :exec
DELETE FROM comics
WHERE id = $1;

-- Select all chapters for a specific comic
-- name: GetChaptersByComicID :many
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
    chapter_number;

-- Select a chapter by id
-- name: GetChapterById :one
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
    chapters.id = $1;

-- Select a chapter by comic slug and chapter number
-- name: GetChapterByComicSlugAndNumber :one
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
    ch.chapter_number = $2;

-- Insert a new chapter
-- name: CreateChapter :one
INSERT INTO chapters (
    comic_id, chapter_number, title, release_date, content_url
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- Select all genres
-- name: GetGenres :many
SELECT
    id,
    genre_name
FROM
    genres
ORDER BY
    genre_name;

-- Select a genre by name
-- name: GetGenreByName :one
SELECT
    id,
    genre_name
FROM
    genres
WHERE
    genres.genre_name = $1;

-- Link a comic with a genre
-- name: AddGenreToComic :exec
INSERT INTO comic_genres (
    comic_id,
    genre_id
) VALUES (
    $1, $2
);

-- Remove a genre from a comic
-- name: RemoveGenreFromComic :exec
DELETE FROM comic_genres
WHERE comic_id = $1 AND genre_id = $2;

-- Select comics by genre
-- name: GetComicsByGenre :many
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
group by comics.id;

-- Select genres by comic id
-- name: GetGenresByComicId :many
SELECT
    genres.id,
    genres.genre_name
FROM
    genres
        JOIN
    comic_genres ON genres.id = comic_genres.genre_id
WHERE
    comic_genres.comic_id = $1;

-- Select comics by title
-- name: GetComicsByTitle :many
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
    comics.upload_date DESC;

-- Register User
-- name: RegisterUser :one
INSERT INTO users (
    username, email, password, first_name, last_name, date_of_birth, role
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- Login User
-- name: LoginWithEmail :one
SELECT
    user_id, username, email, password, first_name, last_name, date_of_birth, role
FROM
    users
WHERE
    users.email = $1;

-- Bookmark a comic
-- name: BookmarkComic :exec
INSERT INTO bookmark (
    user_id,
    comic_id
) VALUES (
     $1, $2
);

-- Remove a comic from user's bookmark
-- name: RemoveComicFromBookmark :exec
DELETE FROM bookmark
WHERE user_id = $1 AND comic_id = $2;