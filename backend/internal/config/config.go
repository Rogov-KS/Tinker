package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
	LogLevel    string
	SeedDB      bool
}

func Load() (*Config, error) {
	viper.SetDefault("PORT", "3001")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("SEED_DB", "false")

	// Читаем из env переменных
	viper.AutomaticEnv()

	// Также можно читать из .env файла если он есть
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./backend")
	_ = viper.ReadInConfig() // Игнорируем ошибку если файла нет

	cfg := &Config{
		DatabaseURL: getEnvOrViper("DATABASE_URL", ""),
		JWTSecret:   getEnvOrViper("JWT_SECRET", "dev-secret"),
		Port:        getEnvOrViper("PORT", "3001"),
		LogLevel:    getEnvOrViper("LOG_LEVEL", "info"),
		SeedDB:      getEnvOrViper("SEED_DB", "false") == "true" || getEnvOrViper("SEED_DB", "false") == "1",
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

func getEnvOrViper(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	if viper.IsSet(key) {
		return viper.GetString(key)
	}
	return defaultValue
}

