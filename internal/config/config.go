package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DB     DBConfig
	JWT    JWTConfig
	Server ServerConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret       string
	ExpiresIn    time.Duration
	RefreshToken time.Duration
}

type ServerConfig struct {
	Port string
}

func (c *DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

func New() *Config {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Aviso: arquivo .env não encontrado: %v\n", err)
	}

	expiresIn, _ := strconv.Atoi(getEnv("JWT_EXPIRES_IN", "3600"))
	refreshToken, _ := strconv.Atoi(getEnv("JWT_REFRESH_TOKEN", "604800"))

	return &Config{
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "expense_db"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:       getEnv("JWT_SECRET", "seu_segredo_jwt"),
			ExpiresIn:    time.Duration(expiresIn) * time.Second,
			RefreshToken: time.Duration(refreshToken) * time.Second,
		},
		Server: ServerConfig{
			Port: getEnv("PORT", "8081"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Load carrega as configurações do ambiente
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("erro ao carregar .env: %w", err)
	}

	config := &Config{
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "expense_user"),
			Password: getEnv("DB_PASSWORD", "expense_password"),
			Name:     getEnv("DB_NAME", "expense_db"),
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

	// Valores padrão
	if config.Server.Port == "" {
		config.Server.Port = "8080"
	}

	// Validações
	if config.DB.Host == "" {
		return nil, fmt.Errorf("DB_HOST não definida")
	}
	if config.DB.Port == "" {
		return nil, fmt.Errorf("DB_PORT não definida")
	}
	if config.DB.User == "" {
		return nil, fmt.Errorf("DB_USER não definida")
	}
	if config.DB.Password == "" {
		return nil, fmt.Errorf("DB_PASSWORD não definida")
	}
	if config.DB.Name == "" {
		return nil, fmt.Errorf("DB_NAME não definida")
	}
	if config.DB.SSLMode == "" {
		return nil, fmt.Errorf("DB_SSL_MODE não definida")
	}
	if config.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET não definida")
	}

	return config, nil
}

func parseDuration(value string) time.Duration {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 24 * time.Hour
	}
	return duration
}
