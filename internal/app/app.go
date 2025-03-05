package app

import (
	"context"
	"log"
	"net"

	"github.com/um3ra/auth-microservice/config"
	authGrpc "github.com/um3ra/auth-microservice/pkg/auth_v1"
	"github.com/um3ra/auth-microservice/pkg/closer"
	"github.com/um3ra/auth-microservice/pkg/interceptor"
	userGrpc "github.com/um3ra/auth-microservice/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const (
	envPath = ".env"
)

func NewApp(ctx context.Context) (*App, error) {
	app := &App{}
	err := app.initDeps(ctx)

	if err != nil {
		return nil, err
	}
	return app, nil
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()
	return a.runGrpcServer()
}

type App struct {
	provider   ServiceProvider
	grpcServer *grpc.Server
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initConfig,
		a.initProvider,
		a.initGrpcServer,
	}
	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initProvider(_ context.Context) error {
	a.provider = NewServiceProvider()
	return nil
}

func (a *App) initConfig(_ context.Context) error {
	if err := config.NewConfig(envPath); err != nil {
		return err
	}
	return nil
}

func (a *App) initGrpcServer(ctx context.Context) error {
	s := grpc.NewServer(grpc.Creds(insecure.NewCredentials()), grpc.UnaryInterceptor(interceptor.ValidateInterceptor))
	a.grpcServer = s
	reflection.Register(a.grpcServer)
	handler := a.provider.UserHandler(ctx)
	authHandler := a.provider.AuthHandler(ctx)
	userGrpc.RegisterUserV1Server(a.grpcServer, handler)
	authGrpc.RegisterAuthV1Server(a.grpcServer, authHandler)
	return nil
}

func (a *App) runGrpcServer() error {
	l, err := net.Listen("tcp", a.provider.GrpcConfig().Address())
	log.Printf("Server is listening on : %s", l.Addr())
	if err != nil {
		return err
	}
	err = a.grpcServer.Serve(l)
	if err != nil {
		return err
	}
	return nil
}
