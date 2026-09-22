package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Env         string
	DatabaseURL string
}

func Load() *Config {

	if err := godotenv.Load(); err != nil {
		log.Printf("Notice: .env file not loaded (%v), falling back to system environment variables", err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("FATAL:DATABASE_URL environment variable is required")
	}

	return &Config{
		Port:        getEnv("PORT", "8080"),
		Env:         getEnv("ENV", "development"),
		DatabaseURL: dbURL,
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
