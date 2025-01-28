package app

import (
	"context"
	"log"

	"github.com/um3ra/auth-microservice/config"
	cfg "github.com/um3ra/auth-microservice/config"
	"github.com/um3ra/auth-microservice/internal/api/user"
	userHandler "github.com/um3ra/auth-microservice/internal/api/user"
	"github.com/um3ra/auth-microservice/internal/client/db"
	"github.com/um3ra/auth-microservice/internal/client/db/pg"
	"github.com/um3ra/auth-microservice/internal/repository"
	userRepo "github.com/um3ra/auth-microservice/internal/repository/user"
	"github.com/um3ra/auth-microservice/internal/service"
	userSrv "github.com/um3ra/auth-microservice/internal/service/user"
	"github.com/um3ra/auth-microservice/pkg/closer"
)

type serviceProvider struct {
	db             db.Client
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

func (s *serviceProvider) GrpcConfig() config.GrpcConfig {
	if s.grpcConfig == nil {
		grpc, err := cfg.NewGrpcConfig()
		if err != nil {
			log.Fatalf("Failed to get grpc config: %s", err.Error())
		}
		s.grpcConfig = grpc
	}
	return s.grpcConfig
}

func (s *serviceProvider) DbClient(ctx context.Context) db.Client {
	if s.db == nil {
		client, err := pg.NewClient(ctx, s.DbConfig().DSN())
		if err != nil {
			log.Fatalf("Failed to database connect: %s\n", err.Error())
		}
		// err = client.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("Failed to ping database: %s\n", err.Error())
		}
		s.db = client
		closer.Add(s.db.Close)
	}

	return s.db
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepo.NewUserRepository(s.DbClient(ctx))
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
