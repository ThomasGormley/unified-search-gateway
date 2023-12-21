package configuration

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	OmdbApiKey string
}

var config Config

// getRequiredEnv tries to get an environment variable and panics if it's not set
func getRequiredEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("required environment variable %s is not set", key)
	}
	return value
}

// getEnv reads an environment variable or returns a default value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func Load() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	config = Config{}

	config.OmdbApiKey = getRequiredEnv("OMDB_API_KEY")

	log.Printf("Config loaded: %+v", config)
}

func Get() *Config {
	return &config
}
