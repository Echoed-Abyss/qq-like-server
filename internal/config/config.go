package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Token    TokenConfig
	TCP      TCPConfig
}

type ServerConfig struct {
	Host string
	Port int
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type TokenConfig struct {
	Secret     string
	ExpireHour int
}

type TCPConfig struct {
	Host string
	Port int
}

var AppConfig *Config

func LoadConfig() {
	AppConfig = &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnvInt("SERVER_PORT", 8080),
			Mode: getEnv("SERVER_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "127.0.0.1"),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "qq_like"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Token: TokenConfig{
			Secret:     getEnv("TOKEN_SECRET", "qq-like-server-secret-key-2024"),
			ExpireHour: getEnvInt("TOKEN_EXPIRE_HOUR", 24*7),
		},
		TCP: TCPConfig{
			Host: getEnv("TCP_HOST", "0.0.0.0"),
			Port: getEnvInt("TCP_PORT", 9090),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
