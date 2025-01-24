package config

import (
	"errors"
	"os"
)

const (
	dsnPgEnv = "PG_DSN"
)

type dbConfig struct {
	dsn string
}

type DbConfig interface {
	DSN() string
}

func NewDbConfig() (DbConfig, error) {
	dsn := os.Getenv(dsnPgEnv)
	if len(dsn) == 0 {
		return nil, errors.New("ENV: db not found")
	}
	return &dbConfig{
		dsn: dsn,
	}, nil
}

func (c *dbConfig) DSN() string {
	return c.dsn
}
