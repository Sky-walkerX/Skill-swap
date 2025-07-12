-- Migration to add admin and ban fields to users table
-- Run this after the initial migration

ALTER TABLE users 
ADD COLUMN IF NOT EXISTS is_admin BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS is_banned BOOLEAN DEFAULT FALSE;

-- Create index for admin lookups
CREATE INDEX IF NOT EXISTS idx_users_is_admin ON users(is_admin);
CREATE INDEX IF NOT EXISTS idx_users_is_banned ON users(is_banned);

-- Create a default admin user (optional - for development)
-- Password is 'admin123' hashed with bcrypt
-- You should change this in production
INSERT INTO users (name, email, password_hash, is_admin) 
VALUES ('Admin User', 'admin@skillswap.com', '$2a$10$YourHashedPasswordHere', TRUE)
ON CONFLICT (email) DO NOTHING;
