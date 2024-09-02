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

-- Gotta know how to create an index somehow
-- CREATE INDEX idx_comics_slug ON comics(slug);

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

CREATE TYPE user_role AS ENUM ('user', 'admin', 'moderator');

CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,  -- Unique identifier for each user
    username VARCHAR(50) UNIQUE NOT NULL,  -- Username, must be unique
    email VARCHAR(100) UNIQUE NOT NULL,  -- Email, must be unique
    password VARCHAR(255) NOT NULL,  -- Password (hashed)
    first_name VARCHAR(50),  -- First name of the user
    last_name VARCHAR(50),  -- Last name of the user
    date_of_birth DATE,  -- Date of birth
    role user_role DEFAULT 'user' NOT NULL,  -- Role of the user (using enum type)
    created_at TIMESTAMP DEFAULT NOW(),  -- Timestamp when the user was created
    updated_at TIMESTAMP DEFAULT NOW() -- Timestamp for the last update
);

-- Join table for comics that is being bookmarked by users
CREATE TABLE bookmark (
    user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
    comic_id INT REFERENCES comics(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, comic_id)
);

CREATE TABLE sessions (
    session_id SERIAL PRIMARY KEY,  -- Unique identifier for each session
    user_id INT REFERENCES users(user_id) ON DELETE CASCADE, -- Foreign key referencing the user_id in the users table
    session_token VARCHAR(512) UNIQUE NOT NULL,  -- Unique session token_maker
    created_at TIMESTAMP DEFAULT NOW(),  -- Timestamp when the session was created
    expires_at TIMESTAMP NOT NULL,  -- Timestamp when the session expires
    is_active BOOLEAN DEFAULT TRUE  -- Status to indicate if the session is currently active
);