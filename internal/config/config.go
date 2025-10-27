package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	GRPCPort             int
	GRPCEnableReflection bool
	LogLevel             string
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

	enableReflectionStr := getEnv("GRPC_ENABLE_REFLECTION", "false")
	enableReflection, err := strconv.ParseBool(enableReflectionStr)
	if err != nil {
		return nil, err
	}
	config.GRPCEnableReflection = enableReflection

	config.LogLevel = getEnv("LOG_LEVEL", "info")

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
