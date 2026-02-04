package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB       DBConfig
	HTTP     HTTPConfig
	LogLevel string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type HTTPConfig struct {
	Port string
}

func Load() (*Config, error) {
	// Load .env, if exists
	godotenv.Load()

	return &Config{
		DB: DBConfig{
			Host:     getEnv("POSTGRES_HOST", "postgres"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", "postgres"),
			Name:     getEnv("POSTGRES_DB", "subscriptions"),
			SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		},
		HTTP: HTTPConfig{
			Port: getEnv("HTTP_PORT", "8080"),
		},
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
