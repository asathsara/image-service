package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	SecretKey         string
	RateLimitMax      int
	RateLimitWindowMs int
}

func LoadConfig() *Config {
	_ = godotenv.Load() // Ignore error if .env is missing

	return &Config{
		Port:              getEnv("PORT", "4000"),
		SecretKey:         getEnv("IMAGE_COMPRESSOR_SECRET", ""),
		RateLimitMax:      getEnvInt("RATE_LIMIT_MAX", 20),
		RateLimitWindowMs: getEnvInt("RATE_LIMIT_WINDOW_MS", 60000),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}
