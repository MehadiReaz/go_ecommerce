package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Payment  PaymentConfig
	Email    EmailConfig
}

type ServerConfig struct {
	Port         string
	Environment  string
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret     string
	ExpiryHours int
}

type PaymentConfig struct {
	StripeSecretKey string
	StripePublicKey string
	BkashAppKey     string
	BkashAppSecret  string
	BkashUsername   string
	BkashPassword   string
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			Environment:  getEnv("ENVIRONMENT", "development"),
			ReadTimeout:  getEnvAsInt("SERVER_READ_TIMEOUT", 10),
			WriteTimeout: getEnvAsInt("SERVER_WRITE_TIMEOUT", 10),
			IdleTimeout:  getEnvAsInt("SERVER_IDLE_TIMEOUT", 120),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "ecommerce"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "your-secret-key"),
			ExpiryHours: getEnvAsInt("JWT_EXPIRY_HOURS", 24),
		},
		Payment: PaymentConfig{
			StripeSecretKey: getEnv("STRIPE_SECRET_KEY", ""),
			StripePublicKey: getEnv("STRIPE_PUBLIC_KEY", ""),
			BkashAppKey:     getEnv("BKASH_APP_KEY", ""),
			BkashAppSecret:  getEnv("BKASH_APP_SECRET", ""),
			BkashUsername:   getEnv("BKASH_USERNAME", ""),
			BkashPassword:   getEnv("BKASH_PASSWORD", ""),
		},
		Email: EmailConfig{
			SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort:     getEnv("SMTP_PORT", "587"),
			SMTPUsername: getEnv("SMTP_USERNAME", ""),
			SMTPPassword: getEnv("SMTP_PASSWORD", ""),
			FromEmail:    getEnv("FROM_EMAIL", "noreply@ecommerce.com"),
			FromName:     getEnv("FROM_NAME", "E-Commerce"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks if all required configurations are set
func (c *Config) Validate() error {
	if c.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	if c.JWT.Secret == "" || c.JWT.Secret == "your-secret-key" {
		return fmt.Errorf("JWT_SECRET must be set to a secure value")
	}
	return nil
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
