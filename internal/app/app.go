package app

import (
	"context"
	"log"

	"net"

	"github.com/um3ra/auth-microservice/config"
	"github.com/um3ra/auth-microservice/pkg/closer"
	userGrpc "github.com/um3ra/auth-microservice/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func Bootstrap() {
	if err := config.NewConfig(); err != nil {
		log.Fatal(err)
	}

	defer func () {
		log.Println("Close functions")
		closer.CloseAll()
		closer.Wait()
	}()

	context := context.Background()
	provider := newServiceProvider()
	grpcConf := provider.GrpcConfig(context)

	l, err := net.Listen("tcp", grpcConf.Address())

	if err != nil {
		log.Fatalf("Failed to listen: %s", err.Error())
	}
	s := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	userGrpc.RegisterUserV1Server(s, provider.UserHandler(context))
	reflection.Register(s)
	if err := s.Serve(l); err != nil {
		log.Fatalf("Failed to serve: %s", err.Error())
	}
}