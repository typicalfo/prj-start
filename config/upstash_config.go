package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type UpstashConfig struct {
	URL               string `yaml:"url"`
	Token             string `yaml:"token"`
	IndexURL          string `yaml:"indexurl"`
	Email             string `yaml:"email"`
	APIKey            string `yaml:"apikey"`
	BatchSize         int    `yaml:"batchsize"`
	ProcessingTimeout int    `yaml:"processingtimeout"`
	LogLevel          string `yaml:"loglevel"`
}

func LoadUpstashConfig() *UpstashConfig {
	// First, check if .env file exists and read it directly to debug
	envPath := ".env"
	if _, err := os.Stat(envPath); err == nil {
		fmt.Printf("DEBUG: Reading .env file directly for debugging...\n")
		envContent, err := os.ReadFile(envPath)
		if err != nil {
			fmt.Printf("DEBUG: Failed to read .env file: %v\n", err)
		} else {
			fmt.Printf("DEBUG: Raw .env file content:\n%s\n", string(envContent))

			// Look for the URL line specifically
			lines := strings.Split(string(envContent), "\n")
			for i, line := range lines {
				if strings.Contains(line, "UPSTASH_VECTOR_REST_URL") {
					fmt.Printf("DEBUG: Line %d: '%s'\n", i, line)
					fmt.Printf("DEBUG: Line %d bytes: %v\n", i, []byte(line))
				}
			}
		}
	}

	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		fmt.Printf("DEBUG: Could not load .env file: %v\n", err)
	} else {
		fmt.Printf("DEBUG: Successfully loaded .env file with godotenv\n")
	}

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
		// Debug: print the exact bytes we got
		fmt.Printf("DEBUG getEnv(%s): bytes=%v, len=%d\n", key, []byte(value), len(value))
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
