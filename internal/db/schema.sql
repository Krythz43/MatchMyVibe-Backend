-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    spotify_uri TEXT UNIQUE NOT NULL,
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    token_expiry TIMESTAMP NOT NULL,
    name TEXT,
    university_name TEXT,
    work JSONB,
    home_town TEXT,
    height TEXT,
    age TEXT,
    zodiac TEXT,
    currently_playing TEXT,
    unix_timestamp BIGINT,
    gender TEXT CHECK (gender IN ('Man', 'Woman', 'Non-binary')),
    dating_preference TEXT CHECK (dating_preference IN ('Men', 'Women', 'Everyone')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create images table
CREATE TABLE IF NOT EXISTS images (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    data BYTEA NOT NULL
);

-- Create interests table
CREATE TABLE IF NOT EXISTS interests (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL
);

-- Create interest_ratings table
CREATE TABLE IF NOT EXISTS interest_ratings (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    rating INTEGER NOT NULL
);

-- Create prompts table
CREATE TABLE IF NOT EXISTS prompts (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    question TEXT NOT NULL,
    answer TEXT NOT NULL
);

-- Create artists table
CREATE TABLE IF NOT EXISTS artists (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    uri TEXT NOT NULL,
    image_url TEXT
);

-- Create songs table
CREATE TABLE IF NOT EXISTS songs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    artist TEXT NOT NULL,
    uri TEXT NOT NULL,
    image_url TEXT
);

-- Create playlists table
CREATE TABLE IF NOT EXISTS playlists (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    uri TEXT NOT NULL,
    image_url TEXT
);

-- Create indexes for faster queries
CREATE INDEX IF NOT EXISTS idx_images_user_id ON images(user_id);
CREATE INDEX IF NOT EXISTS idx_interests_user_id ON interests(user_id);
CREATE INDEX IF NOT EXISTS idx_interest_ratings_user_id ON interest_ratings(user_id);
CREATE INDEX IF NOT EXISTS idx_prompts_user_id ON prompts(user_id);
CREATE INDEX IF NOT EXISTS idx_artists_user_id ON artists(user_id);
CREATE INDEX IF NOT EXISTS idx_songs_user_id ON songs(user_id);
CREATE INDEX IF NOT EXISTS idx_playlists_user_id ON playlists(user_id); 