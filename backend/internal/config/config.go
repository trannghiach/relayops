package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort             string
	DatabaseURL         string
	ReadonlyDatabaseURL string
	NATSURL             string
	DemoEnabled         bool
	DemoAPIKey          string
	DemoMaxEvents       int
	DemoRandomFailure   bool
	DemoFailureRate     float64
	DemoFastBackoff     bool
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := Config{
		AppPort:             getEnv("APP_PORT", "8080"),
		DatabaseURL:         getEnv("DATABASE_URL", "postgresql://postgres:postgres@localhost/relayops?sslmode=disable"),
		ReadonlyDatabaseURL: getEnv("READONLY_DATABASE_URL", "postgresql://relayops_readonly:readonly_password_change_me@localhost/relayops?sslmode=disable"),
		NATSURL:             getEnv("NATS_URL", "nats://localhost:4222"),
		DemoEnabled:         getBoolEnv("DEMO_ENABLED", false),
		DemoAPIKey:          getEnv("DEMO_API_KEY", ""),
		DemoMaxEvents:       1,
		DemoRandomFailure:   getBoolEnv("DEMO_RANDOM_FAILURE", false),
		DemoFailureRate:     getFloatEnv("DEMO_FAILURE_RATE", 0.5),
		DemoFastBackoff:     getBoolEnv("DEMO_FAST_BACKOFF", false),
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}

func getBoolEnv(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value == "true"
}

func getIntEnv(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func getFloatEnv(key string, fallback float64) float64 {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fallback
	}

	return parsed
}
