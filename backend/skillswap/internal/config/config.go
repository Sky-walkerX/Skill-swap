package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl string
}

func Load() Config {
	_ = godotenv.Load()           
	dbURL := os.Getenv("DB_URL")  

	if dbURL == "" {
		log.Fatal("DB_URL missing")
	}
	return Config{DBUrl: dbURL}
}