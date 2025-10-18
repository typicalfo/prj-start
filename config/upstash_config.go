package config

import (
	"os"
)

type UpstashConfig struct {
	URL      string
	Token    string
	IndexURL string
}

func LoadUpstashConfig() *UpstashConfig {
	return &UpstashConfig{
		URL:      getEnv("UPSTASH_VECTOR_URL", ""),
		Token:    getEnv("UPSTASH_VECTOR_TOKEN", ""),
		IndexURL: getEnv("UPSTASH_VECTOR_INDEX_URL", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *UpstashConfig) Validate() bool {
	return c.URL != "" && c.Token != ""
}
