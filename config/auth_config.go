package config

import (
	"errors"
	"os"
)

const (
	secretKey = "JWT_SECRET"
)

type AuthConfig interface {
	Secret() string
}

type authConfig struct {
	secret string
}

func NewAuthConfig() (AuthConfig, error) {
	secret := os.Getenv(secretKey)
	if len(secret) == 0 {
		return nil, errors.New("JWT Secret env must be provided!")
	}
	return &authConfig{
		secret: secret,
	}, nil
}

func (c *authConfig) Secret() string {
	return c.secret
}
