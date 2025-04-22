package config

import (
	"os"
)

type Config struct {
	Port     string
	MongoURI string
	MongoDB  string
}

func Load() Config {
	cfg := Config{
		Port:     getEnv("PORT", "8080"),
		MongoURI: getEnv("DB_URI", "mongodb://localhost:27017"),
		MongoDB:  getEnv("DB_NAME", "swiftdb"),
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}
