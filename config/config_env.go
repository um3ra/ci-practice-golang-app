package config

import (
	env "github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Db DbConfig
}

type DbConfig struct {
	DSN string
}


func NewConfig() *Config {
	if err := env.Load(); err != nil {
		log.Println("Error loading data from .env file, using default config")
	}
	return &Config{
		Db: DbConfig{
			DSN: os.Getenv("DSN"),
		},
	}
}