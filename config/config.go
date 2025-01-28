package config

import (
	env "github.com/joho/godotenv"
)

func NewConfig(path string) error {
	if err := env.Load(path); err != nil {
		return err
	}
	return nil
}
