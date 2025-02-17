package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	DB     DBConfig
	Server ServerConfig
	JWT    JWTConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type ServerConfig struct {
	Port string
}

type JWTConfig struct {
	Secret       string
	ExpiresIn    time.Duration
	RefreshToken time.Duration
}

func New() *Config {
	return &Config{
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "expense_user"),
			Password: getEnv("DB_PASSWORD", "expense_password"),
			DBName:   getEnv("DB_NAME", "expense_db"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Server: ServerConfig{
			Port: getEnv("PORT", "8081"),
		},
		JWT: JWTConfig{
			Secret:       getEnv("JWT_SECRET", "seu_secret_muito_secreto"),
			ExpiresIn:    parseDuration(getEnv("JWT_EXPIRES_IN", "24h")),
			RefreshToken: parseDuration(getEnv("JWT_REFRESH_TOKEN_EXPIRES", "168h")),
		},
	}
}

func (c *DBConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseDuration(value string) time.Duration {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 24 * time.Hour
	}
	return duration
}
