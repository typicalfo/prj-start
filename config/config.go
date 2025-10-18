package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Upstash          UpstashConfig `yaml:"upstash"`
	DefaultNamespace string        `yaml:"default_namespace"`
	BatchSize        int           `yaml:"batch_size"`
	LogLevel         string        `yaml:"log_level"`
	ConfigFile       string        `yaml:"-"`
}

// GetConfigPaths returns possible config file paths in order of preference
func GetConfigPaths() []string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = ""
	}

	var paths []string

	// Current directory
	paths = append(paths, "./prj-start.yaml")
	paths = append(paths, "./.prj-start.yaml")

	// User config directory
	if homeDir != "" {
		switch runtime.GOOS {
		case "windows":
			paths = append(paths, filepath.Join(homeDir, "AppData", "Local", "prj-start", "config.yaml"))
		case "darwin":
			paths = append(paths, filepath.Join(homeDir, ".config", "prj-start", "config.yaml"))
		default: // linux and others
			paths = append(paths, filepath.Join(homeDir, ".config", "prj-start", "config.yaml"))
			paths = append(paths, filepath.Join(homeDir, ".prj-start.yaml"))
		}
	}

	return paths
}

// FindConfigFile searches for existing config files in standard locations
func FindConfigFile() (string, error) {
	paths := GetConfigPaths()

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("no config file found")
}

// LoadConfig loads configuration from file, environment variables, and defaults
func LoadConfig(configFile string) (*Config, error) {
	cfg := &Config{
		DefaultNamespace: "default",
		BatchSize:        10,
		LogLevel:         "info",
	}

	// Try to load from specified file or find one
	if configFile != "" {
		cfg.ConfigFile = configFile
	} else {
		found, err := FindConfigFile()
		if err == nil {
			cfg.ConfigFile = found
		}
	}

	// Load from file if found
	if cfg.ConfigFile != "" {
		if err := loadFromFile(cfg, cfg.ConfigFile); err != nil {
			return nil, fmt.Errorf("failed to load config from %s: %w", cfg.ConfigFile, err)
		}
	}

	// Load environment variables (including .env file)
	_ = godotenv.Load()
	loadFromEnv(cfg)

	return cfg, nil
}

func loadFromFile(cfg *Config, filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, cfg)
}

func loadFromEnv(cfg *Config) {
	// Override with environment variables
	if url := os.Getenv("UPSTASH_VECTOR_REST_URL"); url != "" {
		cfg.Upstash.URL = url
	}
	if token := os.Getenv("UPSTASH_VECTOR_REST_TOKEN"); token != "" {
		cfg.Upstash.Token = token
	}
	if indexURL := os.Getenv("UPSTASH_VECTOR_INDEX_URL"); indexURL != "" {
		cfg.Upstash.IndexURL = indexURL
	}
	if email := os.Getenv("UPSTASH_EMAIL"); email != "" {
		cfg.Upstash.Email = email
	}
	if apiKey := os.Getenv("UPSTASH_API_KEY"); apiKey != "" {
		cfg.Upstash.APIKey = apiKey
	}
	if namespace := os.Getenv("DEFAULT_NAMESPACE"); namespace != "" {
		cfg.DefaultNamespace = namespace
	}
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		cfg.LogLevel = logLevel
	}
}

// SaveConfig saves configuration to the specified file
func SaveConfig(cfg *Config, filename string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func (c *Config) HasUpstashConfig() bool {
	return c.Upstash.URL != "" && c.Upstash.Token != ""
}

func (c *Config) Validate() error {
	if !c.HasUpstashConfig() {
		return fmt.Errorf("upstash configuration is incomplete (URL and Token are required)")
	}
	return nil
}

// GetDefaultConfigPath returns the recommended config file location
func GetDefaultConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	if homeDir == "" {
		return "./prj-start.yaml"
	}

	switch runtime.GOOS {
	case "windows":
		return filepath.Join(homeDir, "AppData", "Local", "prj-start", "config.yaml")
	case "darwin":
		return filepath.Join(homeDir, ".config", "prj-start", "config.yaml")
	default: // linux and others
		return filepath.Join(homeDir, ".config", "prj-start", "config.yaml")
	}
}
