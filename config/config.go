package config

import (
	env "github.com/joho/godotenv"
)

func NewConfig() error {
	if err := env.Load(); err != nil {
		return err
	}
	return nil
}
