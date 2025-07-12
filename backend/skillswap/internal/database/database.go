package database

import (
	"log"

	models "github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Initialize establishes database connection and returns GORM DB instance
func Initialize(dbURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// Test the connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	log.Println("Database connected successfully")

	// Create PostgreSQL extensions and types if they don't exist
	if err := createPostgreSQLExtensions(db); err != nil {
		log.Printf("Warning: Could not create PostgreSQL extensions: %v", err)
	}

	return db, nil
}

// createPostgreSQLExtensions creates necessary PostgreSQL extensions and types
func createPostgreSQLExtensions(db *gorm.DB) error {
	// Enable UUID extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return err
	}

	// Create swap_status enum type
	if err := db.Exec("CREATE TYPE swap_status AS ENUM ('pending', 'accepted', 'rejected', 'cancelled')").Error; err != nil {
		// Ignore error if type already exists
		log.Printf("Info: swap_status type might already exist: %v", err)
	}

	log.Println("✓ PostgreSQL extensions and types created")
	return nil
}

// Migrate runs database migrations
func Migrate(db *gorm.DB) error {
	log.Println("Starting database migrations...")

	// Check if core tables exist first
	if hasAllTables(db) {
		log.Println("✓ All core tables already exist")

		// Still try to seed default skills
		if err := seedDefaultSkills(db); err != nil {
			log.Printf("Warning: Could not seed default skills: %v", err)
		}

		log.Println("✅ Database schema is ready")
		return nil
	}

	// Execute SQL migration script
	if err := runSQLMigration(db); err != nil {
		log.Printf("Error running SQL migration: %v", err)
		return err
	}

	// Insert default skills if they don't exist
	if err := seedDefaultSkills(db); err != nil {
		log.Printf("Warning: Could not seed default skills: %v", err)
	}

	log.Println("✅ Database migrations completed successfully")
	return nil
}

// runSQLMigration executes the SQL migration script
func runSQLMigration(db *gorm.DB) error {
	log.Println("Running SQL migration script...")

	// Read and execute the migration file
	sqlContent := `
-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create custom ENUM types
DO $$ BEGIN
    CREATE TYPE swap_status AS ENUM ('pending', 'accepted', 'rejected', 'cancelled');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Users table
CREATE TABLE IF NOT EXISTS users (
    user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    location TEXT,
    photo_url TEXT,
    is_public BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Skills table
CREATE TABLE IF NOT EXISTS skills (
    skill_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Availability slots table
CREATE TABLE IF NOT EXISTS availability_slots (
    slot_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    label TEXT NOT NULL,
    day_bitmask INTEGER NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_availability_slots_user_id 
        FOREIGN KEY (user_id) REFERENCES users(user_id) 
        ON UPDATE CASCADE ON DELETE CASCADE
);

-- User skills offered junction table
CREATE TABLE IF NOT EXISTS user_skills_offered (
    user_id UUID NOT NULL,
    skill_id UUID NOT NULL,
    
    PRIMARY KEY (user_id, skill_id),
    
    CONSTRAINT fk_user_skills_offered_user_id 
        FOREIGN KEY (user_id) REFERENCES users(user_id) 
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_user_skills_offered_skill_id 
        FOREIGN KEY (skill_id) REFERENCES skills(skill_id) 
        ON UPDATE CASCADE ON DELETE CASCADE
);

-- User skills wanted junction table
CREATE TABLE IF NOT EXISTS user_skills_wanted (
    user_id UUID NOT NULL,
    skill_id UUID NOT NULL,
    
    PRIMARY KEY (user_id, skill_id),
    
    CONSTRAINT fk_user_skills_wanted_user_id 
        FOREIGN KEY (user_id) REFERENCES users(user_id) 
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_user_skills_wanted_skill_id 
        FOREIGN KEY (skill_id) REFERENCES skills(skill_id) 
        ON UPDATE CASCADE ON DELETE CASCADE
);

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

-- Swap ratings table
CREATE TABLE IF NOT EXISTS swap_ratings (
    rating_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    swap_id UUID NOT NULL,
    rater_id UUID NOT NULL,
    ratee_id UUID NOT NULL,
    score SMALLINT NOT NULL CHECK (score >= 1 AND score <= 5),
    comment TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
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

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
CREATE INDEX IF NOT EXISTS idx_skills_name ON skills(name);
CREATE INDEX IF NOT EXISTS idx_availability_slots_user_id ON availability_slots(user_id);
CREATE INDEX IF NOT EXISTS idx_user_skills_offered_user_id ON user_skills_offered(user_id);
CREATE INDEX IF NOT EXISTS idx_user_skills_offered_skill_id ON user_skills_offered(skill_id);
CREATE INDEX IF NOT EXISTS idx_user_skills_wanted_user_id ON user_skills_wanted(user_id);
CREATE INDEX IF NOT EXISTS idx_user_skills_wanted_skill_id ON user_skills_wanted(skill_id);
CREATE INDEX IF NOT EXISTS idx_swap_requests_requester_id ON swap_requests(requester_id);
CREATE INDEX IF NOT EXISTS idx_swap_requests_responder_id ON swap_requests(responder_id);
CREATE INDEX IF NOT EXISTS idx_swap_requests_offered_skill_id ON swap_requests(offered_skill_id);
CREATE INDEX IF NOT EXISTS idx_swap_requests_wanted_skill_id ON swap_requests(wanted_skill_id);
CREATE INDEX IF NOT EXISTS idx_swap_requests_status ON swap_requests(status);
CREATE INDEX IF NOT EXISTS idx_swap_requests_deleted_at ON swap_requests(deleted_at);
CREATE INDEX IF NOT EXISTS idx_swap_ratings_swap_id ON swap_ratings(swap_id);
CREATE INDEX IF NOT EXISTS idx_swap_ratings_rater_id ON swap_ratings(rater_id);
CREATE INDEX IF NOT EXISTS idx_swap_ratings_ratee_id ON swap_ratings(ratee_id);
CREATE INDEX IF NOT EXISTS idx_swap_ratings_score ON swap_ratings(score);

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at columns
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_swap_requests_updated_at ON swap_requests;
CREATE TRIGGER update_swap_requests_updated_at 
    BEFORE UPDATE ON swap_requests 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
`

	if err := db.Exec(sqlContent).Error; err != nil {
		return err
	}

	log.Println("✓ SQL migration completed successfully")
	return nil
}

// hasAllTables checks if all required tables exist
func hasAllTables(db *gorm.DB) bool {
	tables := []interface{}{
		&models.User{},
		&models.Skill{},
		&models.AvailabilitySlot{},
		&models.UserSkillOffered{},
		&models.UserSkillWanted{},
		&models.SwapRequest{},
		&models.SwapRating{},
	}

	for _, table := range tables {
		if !db.Migrator().HasTable(table) {
			return false
		}
	}
	return true
}

// seedDefaultSkills inserts default skills into the database if they don't exist
func seedDefaultSkills(db *gorm.DB) error {
	defaultSkills := []string{
		"JavaScript",
		"Python",
		"React",
		"Node.js",
		"PostgreSQL",
		"Docker",
		"AWS",
		"Machine Learning",
		"UI/UX Design",
		"Data Analysis",
		"Go",
		"TypeScript",
		"GraphQL",
		"Kubernetes",
		"DevOps",
		"HTML/CSS",
		"Java",
		"C++",
		"Swift",
		"Flutter",
		"Angular",
		"Vue.js",
		"MongoDB",
		"Redis",
		"Microservices",
		"System Design",
		"Cybersecurity",
		"Mobile Development",
		"Web Development",
		"Backend Development",
	}

	var existingSkills []models.Skill
	db.Find(&existingSkills)

	existingSkillsMap := make(map[string]bool)
	for _, skill := range existingSkills {
		existingSkillsMap[skill.Name] = true
	}

	var skillsToInsert []models.Skill
	for _, skillName := range defaultSkills {
		if !existingSkillsMap[skillName] {
			skillsToInsert = append(skillsToInsert, models.Skill{
				Name: skillName,
			})
		}
	}

	if len(skillsToInsert) > 0 {
		if err := db.Create(&skillsToInsert).Error; err != nil {
			return err
		}
		log.Printf("✓ Inserted %d default skills", len(skillsToInsert))
	} else {
		log.Println("✓ Default skills already exist")
	}

	return nil
}
