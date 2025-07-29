package config

import "github.com/joho/godotenv"

type Config struct {
	Environment string
}

func Load() *Config{
	// Load .env file
	godotenv.Load()


	return &Config{
		Environment : "developement",
		DatabaseURL
	}

}