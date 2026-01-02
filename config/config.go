package config

import (
	"bufio"
	"os"
	"strings"
)

// Config holds application configuration
type Config struct {
	EmailFrom     string
	EmailPassword string
	Port          string
}

// Load reads configuration from environment variables and .env file
func Load() *Config {
	// Load .env file if it exists
	loadEnvFile(".env")

	cfg := &Config{
		EmailFrom:     getEnv("EMAIL_FROM", ""),
		EmailPassword: getEnv("EMAIL_PASSWORD", ""),
		Port:          getEnv("PORT", "8080"),
	}

	return cfg
}

// loadEnvFile reads a .env file and sets environment variables
func loadEnvFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		return // .env file is optional
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse KEY=VALUE
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Don't override existing environment variables
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
