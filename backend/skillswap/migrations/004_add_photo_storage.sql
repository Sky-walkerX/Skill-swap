-- Migration to add photo storage fields to users table
-- This adds fields to store profile photos directly in the database

ALTER TABLE users 
ADD COLUMN IF NOT EXISTS photo_data BYTEA,
ADD COLUMN IF NOT EXISTS photo_mime_type VARCHAR(100);

-- Create index for photo lookups
CREATE INDEX IF NOT EXISTS idx_users_has_photo ON users(photo_data) WHERE photo_data IS NOT NULL;

-- Update existing photo_url entries to use the new API format
-- This is optional if you have existing photo URLs that need to be updated
UPDATE users 
SET photo_url = CONCAT('http://localhost:8080/api/v1/files/users/', user_id::text, '/photo')
WHERE photo_url IS NOT NULL AND photo_url != '';
