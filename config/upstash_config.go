package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type UpstashConfig struct {
	URL               string
	Token             string
	IndexURL          string
	Email             string
	APIKey            string
	BatchSize         int
	ProcessingTimeout int
	LogLevel          string
}

func LoadUpstashConfig() *UpstashConfig {
	// Load .env file if it exists
	_ = godotenv.Load()

	return &UpstashConfig{
		URL:               getEnv("UPSTASH_VECTOR_REST_URL", ""),
		Token:             getEnv("UPSTASH_VECTOR_REST_TOKEN", ""),
		IndexURL:          getEnv("UPSTASH_VECTOR_INDEX_URL", ""),
		Email:             getEnv("UPSTASH_EMAIL", ""),
		APIKey:            getEnv("UPSTASH_API_KEY", ""),
		BatchSize:         getIntEnv("BATCH_SIZE", 10),
		ProcessingTimeout: getIntEnv("PROCESSING_TIMEOUT_MINUTES", 30),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func (c *UpstashConfig) Validate() bool {
	return c.URL != "" && c.Token != ""
}

func (c *UpstashConfig) HasMCPConfig() bool {
	return c.Email != "" && c.APIKey != ""
}
