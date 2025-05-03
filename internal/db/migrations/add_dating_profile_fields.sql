-- Add new dating profile fields to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS unix_timestamp BIGINT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS gender TEXT CHECK (gender IN ('Man', 'Woman', 'Non-binary'));
ALTER TABLE users ADD COLUMN IF NOT EXISTS dating_preference TEXT CHECK (dating_preference IN ('Man', 'Woman', 'Everyone'));

-- Add comment to explain the fields
COMMENT ON COLUMN users.unix_timestamp IS 'User''s birthday as Unix timestamp';
COMMENT ON COLUMN users.gender IS 'User gender (Man, Woman, or Non-binary)';
COMMENT ON COLUMN users.dating_preference IS 'Dating preference (Man, Woman, or Everyone)'; 