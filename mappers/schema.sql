-- Authors Table
CREATE TABLE authors (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

-- Status Table
CREATE TYPE STATUS_TYPE AS ENUM ('ongoing', 'completed', 'dropped');

-- Comics Table
CREATE TABLE comics (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(255) UNIQUE,
    title VARCHAR(255) NOT NULL UNIQUE,
    status STATUS_TYPE DEFAULT 'ongoing' NOT NULL ,
    summary TEXT,
    release_date DATE NOT NULL,
    cover_image_url VARCHAR(255),
    upload_date TIMESTAMP DEFAULT NOW()
);

-- Join table for authors and comics
CREATE TABLE author_comic (
    author_id INT REFERENCES authors(id) ON DELETE CASCADE,
    comic_id INT REFERENCES comics(id) ON DELETE CASCADE,
    PRIMARY KEY (author_id, comic_id)
);

-- Chapters Table
CREATE TABLE chapters (
    id SERIAL PRIMARY KEY,
    comic_id INT REFERENCES comics(id) ON DELETE CASCADE,
    chapter_number INT NOT NULL,
    title VARCHAR(255),
    content_url VARCHAR(255),
    release_date TIMESTAMP DEFAULT NOW(),
    CONSTRAINT unique_comic_chapter UNIQUE (comic_id, chapter_number)
);

-- Genres Table
CREATE TABLE genres (
    id SERIAL PRIMARY KEY,
    genre_name VARCHAR(50) UNIQUE NOT NULL
);

-- ComicGenres Table (Many-to-Many Relationship)
CREATE TABLE comic_genres (
    comic_id INT REFERENCES comics(id) ON DELETE CASCADE,
    genre_id INT REFERENCES genres(id) ON DELETE CASCADE,
    PRIMARY KEY (comic_id, genre_id)
);