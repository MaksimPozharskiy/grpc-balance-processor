package config

import (
	"os"
)

type Config struct {
	DatabaseDSN      string
	GRPCPort         string
	CancelPeriodMin  int
	LogLevel         string
}

func Load() *Config {
	return &Config{
		DatabaseDSN:     getEnv("DATABASE_DSN", "postgres://user:password@localhost:5432/balance?sslmode=disable"),
		GRPCPort:        getEnv("GRPC_PORT", "8080"),
		CancelPeriodMin: 5, // TODO  потом сделать из енва
		LogLevel:        getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
