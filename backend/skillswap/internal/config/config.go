package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl     string
	Port      string
	JWTSecret string
	UploadDir string
	BaseURL   string
}

func Load() Config {
	_ = godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	port := os.Getenv("PORT")
	jwtSecret := os.Getenv("JWT_SECRET")
	uploadDir := os.Getenv("UPLOAD_DIR")
	baseURL := os.Getenv("BASE_URL")

	if dbURL == "" {
		log.Fatal("DB_URL missing")
	}

	if jwtSecret == "" {
		jwtSecret = "default-secret-change-in-production"
		log.Println("Warning: Using default JWT secret. Set JWT_SECRET environment variable in production.")
	}

	if uploadDir == "" {
		uploadDir = "./uploads"
	}

	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	if port == "" {
		port = "8080"
	}

	return Config{
		DBUrl:     dbURL,
		Port:      port,
		JWTSecret: jwtSecret,
		UploadDir: uploadDir,
		BaseURL:   baseURL,
	}
}
