package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfiguration(t *testing.T) {
	t.Setenv("PORT", "8080")
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/dbname")

	conf := newConfiguration()

	assert.Equal(t, "8080", conf.Port)
	assert.Equal(t, "postgres://user:pass@localhost:5432/dbname", conf.DatabaseURL)
}

func TestNewConfiguration_Defaults(t *testing.T) {
	conf := newConfiguration()

	assert.Equal(t, "3000", conf.Port)
	assert.Equal(t, "", conf.DatabaseURL)
}

func TestGetEnvOrDefault(t *testing.T) {
	t.Setenv("TEST_ENV", "value")

	value := getEnvOrDefault("TEST_ENV", "default")
	assert.Equal(t, "value", value)

	value = getEnvOrDefault("NON_EXISTENT_ENV", "default")
	assert.Equal(t, "default", value)
}
