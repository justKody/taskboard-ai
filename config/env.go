package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseUrl   string
	Port          string
	JWTKey        string
	JWTExpiration int

}

var Envs = initConfig()

func initConfig() Config {
	// Optional locally; in Docker/K8s env vars are injected directly.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return Config{
		DatabaseUrl:   getEnv("DATABASE_URL", "postgresql://test"),
		Port:          getEnv("PORT", "8080"),
		JWTExpiration: getEnvInt("JWT_EXPIRATION", 15),
		JWTKey:        getEnv("JWT_KEY", "secret"),
	}
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return fallback
		}
		return intValue
	}
	return fallback
}
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
