package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl     string
	RedisUrl  string
	JWTSecret string
	Port      string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := &Config{
		DBUrl:     getEnv("DB_URL", "postgres://postgres:postgres@localhost:5432/crowdfunding?sslmode=disable"),
		RedisUrl:  getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret: getEnv("JWT_SECRET", "change-me-in-production"),
		Port:      getEnv("PORT", "8080"),
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
