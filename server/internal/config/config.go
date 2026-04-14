package config

import (
	"os"
	"strings"
)

type Config struct {
	Port               string
	DatabaseURL        string
	CORSAllowedOrigins []string
}

func Load() Config {
	port := getEnv("PORT", "8080")
	databaseURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@postgres:5432/project_tracker?sslmode=disable")
	originsRaw := getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")
	origins := strings.Split(originsRaw, ",")

	return Config{
		Port:               port,
		DatabaseURL:        databaseURL,
		CORSAllowedOrigins: origins,
	}
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
