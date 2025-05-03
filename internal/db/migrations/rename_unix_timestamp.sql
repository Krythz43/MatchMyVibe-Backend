-- Rename unix_timestamp column to birthdayInUnix
ALTER TABLE users RENAME COLUMN unix_timestamp TO "birthdayInUnix";

-- Update the comment
COMMENT ON COLUMN users."birthdayInUnix" IS 'User''s birthday as Unix timestamp'; 