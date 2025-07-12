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
	return db, nil
}

// Migrate runs database migrations
func Migrate(db *gorm.DB) error {
	// Check if tables already exist (created via SQL script)
	if db.Migrator().HasTable(&models.User{}) {
		log.Println("Tables already exist. Skipping migrations.")
		log.Println("✅ Database schema is ready")
		return nil
	}

	// For PostgreSQL, we need to handle the migration order carefully
	// Let's use AutoMigrate but in the correct dependency order

	log.Println("Starting database migrations...")

	// Step 1: Create independent tables first (User, Skill)
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Printf("Error migrating User table: %v", err)
		return err
	}
	log.Println("✓ Users table migrated")

	if err := db.AutoMigrate(&models.Skill{}); err != nil {
		log.Printf("Error migrating Skill table: %v", err)
		return err
	}
	log.Println("✓ Skills table migrated")

	// Step 2: Create tables that depend on User and/or Skill
	if err := db.AutoMigrate(&models.AvailabilitySlot{}); err != nil {
		log.Printf("Error migrating AvailabilitySlot table: %v", err)
		return err
	}
	log.Println("✓ Availability slots table migrated")

	if err := db.AutoMigrate(&models.UserSkillOffered{}); err != nil {
		log.Printf("Error migrating UserSkillOffered table: %v", err)
		return err
	}
	log.Println("✓ User skills offered table migrated")

	if err := db.AutoMigrate(&models.UserSkillWanted{}); err != nil {
		log.Printf("Error migrating UserSkillWanted table: %v", err)
		return err
	}
	log.Println("✓ User skills wanted table migrated")

	// Step 3: Create tables that depend on previous tables
	if err := db.AutoMigrate(&models.SwapRequest{}); err != nil {
		log.Printf("Error migrating SwapRequest table: %v", err)
		return err
	}
	log.Println("✓ Swap requests table migrated")

	// Step 4: Create tables that depend on swap requests
	if err := db.AutoMigrate(&models.SwapRating{}); err != nil {
		log.Printf("Error migrating SwapRating table: %v", err)
		return err
	}
	log.Println("✓ Swap ratings table migrated")

	log.Println("✅ Database migrations completed successfully")
	return nil
}
