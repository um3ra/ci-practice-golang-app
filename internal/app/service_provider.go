package app

import (
	"context"

	"github.com/um3ra/auth-microservice/config"
	authApi "github.com/um3ra/auth-microservice/internal/api/auth"
	userApi "github.com/um3ra/auth-microservice/internal/api/user"
	"github.com/um3ra/auth-microservice/internal/client/db"
	"github.com/um3ra/auth-microservice/internal/repository"
	"github.com/um3ra/auth-microservice/internal/service"
	"github.com/um3ra/auth-microservice/pkg/jwt"
)

type ServiceProvider interface {
	DbConfig() config.DbConfig
	GrpcConfig() config.GrpcConfig
	AuthConfig() config.AuthConfig
	DbClient(ctx context.Context) db.Client
	UserRepository(ctx context.Context) repository.UserRepository
	UserService(ctx context.Context) service.UserService
	JwtService() jwt.JwtService
	AuthService(ctx context.Context) service.AuthService
	UserHandler(ctx context.Context) *userApi.UserHandler
	AuthHandler(ctx context.Context) *authApi.AuthHandler
}
