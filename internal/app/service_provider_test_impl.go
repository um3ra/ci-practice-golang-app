package app

import (
	"context"
	"log"
	"os"

	"github.com/um3ra/auth-microservice/config"
	authApi "github.com/um3ra/auth-microservice/internal/api/auth"
	userApi "github.com/um3ra/auth-microservice/internal/api/user"
	"github.com/um3ra/auth-microservice/internal/client/db"
	"github.com/um3ra/auth-microservice/internal/client/db/pg"
	"github.com/um3ra/auth-microservice/internal/repository"
	userRepo "github.com/um3ra/auth-microservice/internal/repository/user"
	"github.com/um3ra/auth-microservice/internal/service"
	authSrv "github.com/um3ra/auth-microservice/internal/service/auth"
	userSrv "github.com/um3ra/auth-microservice/internal/service/user"
	"github.com/um3ra/auth-microservice/pkg/closer"
	jwtSrv "github.com/um3ra/auth-microservice/pkg/jwt"
)

var (
	TEST_DB_DSN = os.Getenv("TEST_DB_DSN")
)

type serviceProviderTest struct {
	db             db.Client
	grpcConfig     config.GrpcConfig
	authConfig     config.AuthConfig
	userRepository repository.UserRepository
	userService    service.UserService
	userHandler    *userApi.UserHandler
	jwtService     jwtSrv.JwtService
	authService    service.AuthService
	authHandler    *authApi.AuthHandler
}

func NewServiceProviderTest() *serviceProviderTest {
	return &serviceProviderTest{}
}

func (s *serviceProviderTest) GrpcConfig() config.GrpcConfig {
	if s.grpcConfig == nil {
		grpc, err := config.NewGrpcConfig()
		if err != nil {
			log.Fatalf("Failed to get grpc config: %s", err.Error())
		}
		s.grpcConfig = grpc
	}
	return s.grpcConfig
}

func (s *serviceProviderTest) AuthConfig() config.AuthConfig {
	if s.authConfig == nil {
		authCfg, err := config.NewAuthConfig()
		if err != nil {
			log.Fatalf("Failed to get auth config: %s", err.Error())
		}
		s.authConfig = authCfg
	}
	return s.authConfig
}

func (s *serviceProviderTest) DbClient(ctx context.Context) db.Client {
	if s.db == nil {
		if len(TEST_DB_DSN) == 0 {
			log.Fatalf("test db dsn must be provided!")
		}
		client, err := pg.NewClient(ctx, TEST_DB_DSN)
		if err != nil {
			log.Fatalf("Failed to database connect: %s\n", err.Error())
		}
		err = client.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("Failed to ping database: %s\n", err.Error())
		}
		s.db = client
		closer.Add(s.db.Close)
	}

	return s.db
}

func (s *serviceProviderTest) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepo.NewUserRepository(s.DbClient(ctx))
	}
	return s.userRepository
}

func (s *serviceProviderTest) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userSrv.NewUserService(s.UserRepository(ctx))
	}
	return s.userService
}

func (s *serviceProviderTest) JwtService() jwtSrv.JwtService {
	if s.jwtService == nil {
		s.jwtService = jwtSrv.NewJwtService("test")
	}
	return s.jwtService
}

func (s *serviceProviderTest) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authSrv.NewAuthService(s.UserRepository(ctx), s.JwtService())
	}
	return s.authService
}

func (s *serviceProviderTest) UserHandler(ctx context.Context) *userApi.UserHandler {
	if s.userHandler == nil {
		s.userHandler = userApi.NewUserHandler(s.UserService(ctx))
	}
	return s.userHandler
}

func (s *serviceProviderTest) AuthHandler(ctx context.Context) *authApi.AuthHandler {
	if s.authHandler == nil {
		s.authHandler = authApi.NewAuthHandler(s.AuthService(ctx))
	}
	return s.authHandler
}
