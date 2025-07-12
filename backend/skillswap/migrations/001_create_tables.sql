-- Skill Swap Platform Database Schema
-- PostgreSQL Schema Creation Script

-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create custom ENUM types
CREATE TYPE swap_status AS ENUM ('pending', 'accepted', 'rejected', 'cancelled');

-- =====================================================
-- 1. Independent Tables (no foreign keys)
-- =====================================================

-- Users table
CREATE TABLE IF NOT EXISTS users (
    user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR NOT NULL,
    email VARCHAR UNIQUE NOT NULL,
    password_hash VARCHAR NOT NULL,
    location VARCHAR,
    photo_url VARCHAR,
    is_public BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for users table
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- Skills table
CREATE TABLE IF NOT EXISTS skills (
    skill_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for skills table
CREATE INDEX IF NOT EXISTS idx_skills_name ON skills(name);

-- =====================================================
-- 2. Tables with foreign keys to independent tables
-- =====================================================

-- Availability slots table
CREATE TABLE IF NOT EXISTS availability_slots (
    slot_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    label VARCHAR NOT NULL,
    day_bitmask INTEGER NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key constraints
    CONSTRAINT fk_availability_slots_user_id 
        FOREIGN KEY (user_id) REFERENCES users(user_id) 
        ON UPDATE CASCADE ON DELETE CASCADE
);

-- Create indexes for availability_slots table
CREATE INDEX IF NOT EXISTS idx_availability_slots_user_id ON availability_slots(user_id);

-- User skills offered junction table
CREATE TABLE IF NOT EXISTS user_skills_offered (
    user_id UUID NOT NULL,
    skill_id UUID NOT NULL,
    
    -- Composite primary key
    PRIMARY KEY (user_id, skill_id),
    
    -- Foreign key constraints
    CONSTRAINT fk_user_skills_offered_user_id 
        FOREIGN KEY (user_id) REFERENCES users(user_id) 
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_user_skills_offered_skill_id 
        FOREIGN KEY (skill_id) REFERENCES skills(skill_id) 
        ON UPDATE CASCADE ON DELETE CASCADE
);

-- Create indexes for user_skills_offered table
CREATE INDEX IF NOT EXISTS idx_user_skills_offered_user_id ON user_skills_offered(user_id);
CREATE INDEX IF NOT EXISTS idx_user_skills_offered_skill_id ON user_skills_offered(skill_id);

-- User skills wanted junction table
CREATE TABLE IF NOT EXISTS user_skills_wanted (
    user_id UUID NOT NULL,
    skill_id UUID NOT NULL,
    
    -- Composite primary key
    PRIMARY KEY (user_id, skill_id),
    
    -- Foreign key constraints
    CONSTRAINT fk_user_skills_wanted_user_id 
        FOREIGN KEY (user_id) REFERENCES users(user_id) 
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_user_skills_wanted_skill_id 
        FOREIGN KEY (skill_id) REFERENCES skills(skill_id) 
        ON UPDATE CASCADE ON DELETE CASCADE
);

-- Create indexes for user_skills_wanted table
CREATE INDEX IF NOT EXISTS idx_user_skills_wanted_user_id ON user_skills_wanted(user_id);
CREATE INDEX IF NOT EXISTS idx_user_skills_wanted_skill_id ON user_skills_wanted(skill_id);

-- =====================================================
-- 3. Tables with foreign keys to multiple tables
-- =====================================================

-- Swap requests table
CREATE TABLE IF NOT EXISTS swap_requests (
    swap_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    requester_id UUID NOT NULL,
    responder_id UUID NOT NULL,
    offered_skill_id UUID NOT NULL,
    wanted_skill_id UUID NOT NULL,
    status swap_status DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    -- Foreign key constraints
    CONSTRAINT fk_swap_requests_requester_id 
        FOREIGN KEY (requester_id) REFERENCES users(user_id) 
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_swap_requests_responder_id 
        FOREIGN KEY (responder_id) REFERENCES users(user_id) 
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_swap_requests_offered_skill_id 
        FOREIGN KEY (offered_skill_id) REFERENCES skills(skill_id) 
        ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_swap_requests_wanted_skill_id 
        FOREIGN KEY (wanted_skill_id) REFERENCES skills(skill_id) 
        ON UPDATE CASCADE ON DELETE RESTRICT
);

-- Create indexes for swap_requests table
CREATE INDEX IF NOT EXISTS idx_swap_requests_requester_id ON swap_requests(requester_id);
CREATE INDEX IF NOT EXISTS idx_swap_requests_responder_id ON swap_requests(responder_id);
CREATE INDEX IF NOT EXISTS idx_swap_requests_offered_skill_id ON swap_requests(offered_skill_id);
CREATE INDEX IF NOT EXISTS idx_swap_requests_wanted_skill_id ON swap_requests(wanted_skill_id);
CREATE INDEX IF NOT EXISTS idx_swap_requests_status ON swap_requests(status);
CREATE INDEX IF NOT EXISTS idx_swap_requests_deleted_at ON swap_requests(deleted_at);

-- =====================================================
-- 4. Tables with foreign keys to dependent tables
-- =====================================================

-- Swap ratings table
CREATE TABLE IF NOT EXISTS swap_ratings (
    rating_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    swap_id UUID NOT NULL,
    rater_id UUID NOT NULL,
    ratee_id UUID NOT NULL,
    score SMALLINT NOT NULL CHECK (score >= 1 AND score <= 5),
    comment TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key constraints
    CONSTRAINT fk_swap_ratings_swap_id 
        FOREIGN KEY (swap_id) REFERENCES swap_requests(swap_id) 
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_swap_ratings_rater_id 
        FOREIGN KEY (rater_id) REFERENCES users(user_id) 
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_swap_ratings_ratee_id 
        FOREIGN KEY (ratee_id) REFERENCES users(user_id) 
        ON UPDATE CASCADE ON DELETE CASCADE
);

-- Create indexes for swap_ratings table
CREATE INDEX IF NOT EXISTS idx_swap_ratings_swap_id ON swap_ratings(swap_id);
CREATE INDEX IF NOT EXISTS idx_swap_ratings_rater_id ON swap_ratings(rater_id);
CREATE INDEX IF NOT EXISTS idx_swap_ratings_ratee_id ON swap_ratings(ratee_id);
CREATE INDEX IF NOT EXISTS idx_swap_ratings_score ON swap_ratings(score);

-- =====================================================
-- 5. Additional constraints and triggers
-- =====================================================

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at columns
CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_swap_requests_updated_at 
    BEFORE UPDATE ON swap_requests 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- 6. Insert some sample skills (optional)
-- =====================================================

INSERT INTO skills (name) VALUES 
    ('JavaScript'),
    ('Python'),
    ('React'),
    ('Node.js'),
    ('PostgreSQL'),
    ('Docker'),
    ('AWS'),
    ('Machine Learning'),
    ('UI/UX Design'),
    ('Data Analysis'),
    ('Go'),
    ('TypeScript'),
    ('GraphQL'),
    ('Kubernetes'),
    ('DevOps')
ON CONFLICT (name) DO NOTHING;

-- =====================================================
-- Success message
-- =====================================================

DO $$
BEGIN
    RAISE NOTICE 'Skill Swap Platform database schema created successfully!';
    RAISE NOTICE 'Tables created: users, skills, availability_slots, user_skills_offered, user_skills_wanted, swap_requests, swap_ratings';
    RAISE NOTICE 'Sample skills inserted. You can now start the backend server.';
END $$;
