package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	MongoURI    string
	MongoDBName string
	PokeAPIURL  string
	RateLimitMs int
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading from environment variables.")
	}

	return &Config{
		Port:        getEnv("PORT", "4001"),
		MongoURI:    getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDBName: getEnv("MONGO_DB_NAME", "pokemondb"),
		PokeAPIURL:  getEnv("POKEAPI_URL", "https://pokeapi.co/api/v2"),
		RateLimitMs: getEnvAsInt("POKEAPI_RATE_LIMIT_MS", 1000),
	}
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("Warning: Invalid integer value for %s, using default %d\n", key, defaultValue)
			return defaultValue
		}
		return intValue
	}
	return defaultValue
}
