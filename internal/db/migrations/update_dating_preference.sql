-- First update existing values to match the new format
UPDATE users SET dating_preference = 'Men' WHERE dating_preference = 'Man';
UPDATE users SET dating_preference = 'Women' WHERE dating_preference = 'Woman';

-- Temporarily allow NULL values for dating_preference
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_dating_preference_check;

-- Add the new constraint with updated values
ALTER TABLE users ADD CONSTRAINT users_dating_preference_check 
    CHECK (dating_preference IN ('Men', 'Women', 'Everyone'));

-- Update column comment
COMMENT ON COLUMN users.dating_preference IS 'Dating preference (Men, Women, or Everyone)'; 