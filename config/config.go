package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Environment variables for logging
const (
	AppName      = "APP_NAME"
	LogLevel     = "LOG_LEVEL"
	LogstashHost = "LOGSTASH_HOST"
	LogstashPort = "LOGSTASH_PORT"
)

type EnvConfig struct {
	AppName      string
	LogLevel     string
	LogstashHost string
	LogstashPort int
}

// Load environment variables with godotenv and initialize configuration
func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func GetConfig() *EnvConfig {
	return &EnvConfig{
		AppName:      os.Getenv(AppName),
		LogLevel:     os.Getenv(LogLevel),
		LogstashHost: os.Getenv(LogstashHost),
		LogstashPort: getEnvAsInt(LogstashPort, 0),
	}
}

// Parse environment variable as int, with default value
func getEnvAsInt(key string, defaultVal int) int {
	if value, exists := os.LookupEnv(key); exists {
		val, err := strconv.Atoi(value)
		if err == nil {
			return val
		}
		log.Println("Invalid integer value for environment variable", key)
	}
	return defaultVal
}
