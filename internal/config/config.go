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
	EnableHTTPHandler    bool
	HTTPPort             int
	LogLevel             string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		GRPCPort:             mustGetInt("GRPC_PORT", 50051), //nolint:mnd // false-positive
		GRPCEnableReflection: mustGetBool("GRPC_ENABLE_REFLECTION", false),
		EnableHTTPHandler:    mustGetBool("HTTP_HANDLER_ENABLE", false),
		HTTPPort:             mustGetInt("HTTP_PORT", 8080), //nolint:mnd // false-positive
		LogLevel:             getEnv("LOG_LEVEL", "info"),
	}, nil
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

func mustGetInt(key string, def int) int {
	val := getEnv(key, strconv.Itoa(def))
	n, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("invalid int for %s: %v", key, err)
	}
	return n
}

func mustGetBool(key string, def bool) bool {
	val := getEnv(key, strconv.FormatBool(def))
	b, err := strconv.ParseBool(val)
	if err != nil {
		log.Fatalf("invalid bool for %s: %v", key, err)
	}
	return b
}
