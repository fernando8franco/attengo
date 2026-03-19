package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DatabaseURL    string
	Env            string
	SecretJWT      string
	IssuerJWT      string
	ExpirationTime time.Duration
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, reading from enviroment")
	}

	return &Config{
		Port:           getEnv("PORT", ":8080"),
		DatabaseURL:    getEnv("DATABASE_URL", "./data/assistance.db"),
		Env:            getEnv("ENV", "development"),
		SecretJWT:      getEnv("SecretJWT", "fernando"),
		IssuerJWT:      getEnv("IssuerJWT", "franco"),
		ExpirationTime: time.Hour,
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
