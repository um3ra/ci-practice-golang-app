package app

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/um3ra/auth-microservice/config"
	cfg "github.com/um3ra/auth-microservice/config"
	"github.com/um3ra/auth-microservice/internal/api/user"
	userHandler "github.com/um3ra/auth-microservice/internal/api/user"
	"github.com/um3ra/auth-microservice/internal/repository"
	userRepo "github.com/um3ra/auth-microservice/internal/repository/user"
	"github.com/um3ra/auth-microservice/internal/service"
	userSrv "github.com/um3ra/auth-microservice/internal/service/user"
)

type serviceProvider struct {
	pgPool         *pgxpool.Pool
	dbConfig       config.DbConfig
	grpcConfig     config.GrpcConfig
	userRepository repository.UserRepository
	userService    service.UserService
	userHandler    *user.UserHandler
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) DbConfig() config.DbConfig {
	if s.dbConfig == nil {
		dbcfg, err := cfg.NewDbConfig()
		if err != nil {
			log.Fatalf("Failed to get pg config: %s", err.Error())
		}
		s.dbConfig = dbcfg
	}

	return s.dbConfig
}

func (s *serviceProvider) GrpcConfig(ctx context.Context) config.GrpcConfig {
	if s.grpcConfig == nil {
		grpc, err := cfg.NewGrpcConfig()
		if err != nil {
			log.Fatalf("Failed to get grpc config: %s", err.Error())
		}
		s.grpcConfig = grpc
	}

	return s.grpcConfig
}

func (s *serviceProvider) PgPool(ctx context.Context) *pgxpool.Pool {
	if s.pgPool == nil {
		pool, err := pgxpool.New(ctx, s.DbConfig().DSN())
		if err != nil {
			log.Fatalf("Db connection error: %s", err.Error())
		}
		s.pgPool = pool
	}
	return s.pgPool
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepo.NewUserRepository(s.PgPool(ctx))
	}
	return s.userRepository
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userSrv.NewUserService(s.UserRepository(ctx))
	}
	return s.userService
}

func (s *serviceProvider) UserHandler(ctx context.Context) *user.UserHandler {
	if s.userHandler == nil {
		s.userHandler = userHandler.NewUserHandler(s.UserService(ctx))
	}
	return s.userHandler
}
