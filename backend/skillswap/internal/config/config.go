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
}

func Load() Config {
	_ = godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	port := os.Getenv("PORT")
	jwtSecret := os.Getenv("JWT_SECRET")

	if dbURL == "" {
		log.Fatal("DB_URL missing")
	}

	if jwtSecret == "" {
		jwtSecret = "default-secret-change-in-production"
		log.Println("Warning: Using default JWT secret. Set JWT_SECRET environment variable in production.")
	}

	return Config{
		DBUrl:     dbURL,
		Port:      port,
		JWTSecret: jwtSecret,
	}
}
