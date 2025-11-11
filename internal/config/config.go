package config

import (
	"fmt"
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
	DBHost               string
	DBPort               int
	DBUser               string
	DBPassword           string
	DBName               string
	RedisURI             string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		GRPCPort:             mustGetInt("GRPC_PORT", 50051), //nolint:mnd // false-positive
		GRPCEnableReflection: mustGetBool("GRPC_ENABLE_REFLECTION", false),
		EnableHTTPHandler:    mustGetBool("HTTP_HANDLER_ENABLE", false),
		HTTPPort:             mustGetInt("HTTP_PORT", 8080), //nolint:mnd // false-positive
		LogLevel:             getEnv("LOG_LEVEL", "info"),
		DBHost:               getEnv("POSTGRES_HOST", "localhost"),
		DBPort:               mustGetInt("POSTGRES_PORT", 5432),
		DBUser:               getEnv("POSTGRES_USERNAME", "postgres"),
		DBPassword:           getEnv("POSTGRES_PASSWORD", "postgres"),
		DBName:               getEnv("POSTGRES_DATABASE", "postgres"),
		RedisURI:             getEnv("REDIS_URI", "redis://localhost:6379"),
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

func (c Config) BuildDsn() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}
