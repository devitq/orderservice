package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	GRPCPort int
	LogLevel string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	config := &Config{}

	portStr := getEnv("GRPC_PORT", "50051")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}
	config.GRPCPort = port

	config.LogLevel = getEnv("LOG_LEVEL", "info")

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
