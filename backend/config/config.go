package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Auth     AuthConfig
	AI       AIConfig
}

type DatabaseConfig struct {
	Type     string // sqlite, postgres, mysql
	Path     string // for sqlite
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
	CloudProvider string // aws, gcp, azure, or local
}

type ServerConfig struct {
	Port string
	Host string
}

type AuthConfig struct {
	OTPExpiry int // seconds
	JWTSecret string
	OTPLength int
}

type AIConfig struct {
	Provider string // openai, google, huggingface, etc.
	APIKey   string
}

func LoadConfig() *Config {
	godotenv.Load()

	return &Config{
		Database: DatabaseConfig{
			Type:          getEnv("DB_TYPE", "sqlite"),
			Path:          getEnv("DB_PATH", "./clarity.db"),
			CloudProvider: getEnv("CLOUD_PROVIDER", "local"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "50051"),
			Host: getEnv("SERVER_HOST", "localhost"),
		},
		Auth: AuthConfig{
			OTPExpiry: 600, // 10 minutes
			JWTSecret: getEnv("JWT_SECRET", "your-secret-key"),
			OTPLength: 6,
		},
		AI: AIConfig{
			Provider: getEnv("AI_PROVIDER", "openai"),
			APIKey:   getEnv("AI_API_KEY", ""),
		},
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
