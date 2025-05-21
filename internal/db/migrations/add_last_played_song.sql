-- Add last_played_song and user_last_active_at columns to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_played_song JSONB;
ALTER TABLE users ADD COLUMN IF NOT EXISTS user_last_active_at BIGINT;

-- Update the schema version (if tracked)
COMMENT ON TABLE users IS 'Updated to include last_played_song and user_last_active_at'; 