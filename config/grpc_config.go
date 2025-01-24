package config

import (
	"errors"
	"net"
	"os"
)

const (
	grpcHostEnv = "GRPC_HOST"
	grpcPortEnv = "GRPC_PORT"
)

type grpcConfig struct {
	host string
	port string
}

type GrpcConfig interface {
	Address() string
}

func NewGrpcConfig() (GrpcConfig, error) {
	host := os.Getenv(grpcHostEnv)
	if len(host) == 0 {
		return nil, errors.New("ENV: grpc host not found")
	}
	port := os.Getenv(grpcPortEnv)
	if len(port) == 0 {
		return nil, errors.New("ENV: grpc port not found")
	}
	return &grpcConfig{
		host: host,
		port: port,
	}, nil
}

func (c *grpcConfig) Address() string {
	return net.JoinHostPort(c.host, c.port)
}
