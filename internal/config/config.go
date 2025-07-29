package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment     string
	Port           string
	DatabaseURL    string
	DatabaseName   string
	JWTSecret      string
	JWTExpiryHours int
	LogLevel       string
	RateLimitRPM   int
	BCryptCost     int
}

func Load() (*Config, error) {
	// Load .env file only in non-production environments
	if os.Getenv("ENVIRONMENT") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Warning: .env file not found, using system environment variables")
		}
	}

	// Parse integer values with error handling
	jwtExpiryHours, err := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "5"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRY_HOURS: %v", err)
	}

	rateLimitRPM, err := strconv.Atoi(getEnv("RATE_LIMIT_RPM", "50"))
	if err != nil {
		return nil, fmt.Errorf("invalid RATE_LIMIT_RPM: %v", err)
	}

	bcryptCost, err := strconv.Atoi(getEnv("BCRYPT_COST", "12"))
	if err != nil {
		return nil, fmt.Errorf("invalid BCRYPT_COST: %v", err)
	}

	config := &Config{
		Environment:     getEnv("ENVIRONMENT", "development"),
		Port:           getEnv("PORT", "8080"),
		DatabaseURL:    getEnv("DATABASE_URL", "mongodb://localhost:27017"),
		DatabaseName:   getEnv("DATABASE_NAME", "userapi"),
		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key-change-this"),
		JWTExpiryHours: jwtExpiryHours,
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		RateLimitRPM:   rateLimitRPM,
		BCryptCost:     bcryptCost,
	}

	// Validate critical configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %v", err)
	}

	return config, nil
}

// Validate ensures critical configuration values are valid
func (c *Config) Validate() error {
	// Check for default JWT secret in production
	if c.Environment == "production" && c.JWTSecret == "your-secret-key-change-this" {
		return fmt.Errorf("JWT_SECRET must be changed in production")
	}

	// Validate bcrypt cost range
	if c.BCryptCost < 10 || c.BCryptCost > 15 {
		return fmt.Errorf("BCRYPT_COST must be between 10 and 15, got %d", c.BCryptCost)
	}

	// Validate JWT expiry hours
	if c.JWTExpiryHours < 1 || c.JWTExpiryHours > 168 { // 1 hour to 1 week
		return fmt.Errorf("JWT_EXPIRY_HOURS must be between 1 and 168, got %d", c.JWTExpiryHours)
	}

	// Validate log level
	validLogLevels := []string{"debug", "info", "warn", "error"}
	if !contains(validLogLevels, strings.ToLower(c.LogLevel)) {
		return fmt.Errorf("LOG_LEVEL must be one of %v, got %s", validLogLevels, c.LogLevel)
	}

	// Validate rate limit
	if c.RateLimitRPM < 1 || c.RateLimitRPM > 10000 {
		return fmt.Errorf("RATE_LIMIT_RPM must be between 1 and 10000, got %d", c.RateLimitRPM)
	}

	return nil
}

// getEnv retrieves environment variable with fallback to default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// contains checks if slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// MustLoad loads configuration and panics on error (use with caution)
func MustLoad() *Config {
	config, err := Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	return config
}

// PrintConfig prints the configuration (excluding sensitive values)
func (c *Config) PrintConfig() {
	log.Printf("Configuration loaded:")
	log.Printf("  Environment: %s", c.Environment)
	log.Printf("  Port: %s", c.Port)
	log.Printf("  Database: %s", c.DatabaseName)
	log.Printf("  Log Level: %s", c.LogLevel)
	log.Printf("  Rate Limit: %d RPM", c.RateLimitRPM)
	log.Printf("  BCrypt Cost: %d", c.BCryptCost)
	log.Printf("  JWT Expiry: %d hours", c.JWTExpiryHours)
	// Note: Not printing sensitive values like JWT_SECRET, DATABASE_URL
}