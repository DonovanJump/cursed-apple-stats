package config

import (
    "os"
)

type Config struct {
    Port            string
    DatabaseURL     string
    DeadlockBaseURL string
}

func Load() Config {
    return Config{
        Port:            getOrDefault("API_PORT", "8080"),
        DatabaseURL:     getOrDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/cursed_apple_stats?sslmode=disable"),
        DeadlockBaseURL: getOrDefault("DEADLOCK_API_BASE_URL", "https://api.deadlock-api.com"),
    }
}

func getOrDefault(key, fallback string) string {
    value := os.Getenv(key)
    if value == "" {
        return fallback
    }
    return value
}
