package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/um3ra/auth-microservice/internal/converter"
	"github.com/um3ra/auth-microservice/internal/repository/user"
	"github.com/um3ra/auth-microservice/internal/service"
	userService "github.com/um3ra/auth-microservice/internal/service/user"
	userGrpc "github.com/um3ra/auth-microservice/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/um3ra/auth-microservice/config"
)

const PORT = 8888

type userGrpcServer struct {
	userGrpc.UnimplementedUserV1Server
	userService service.UserService
}

func (s *userGrpcServer) Get(ctx context.Context, data *userGrpc.GetRequest) (*userGrpc.GetResponse, error) {
	users, err := s.userService.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	var usersPb []*userGrpc.User

	for _, u := range users {
		usersPb = append(usersPb, converter.ToUserPbFromService(&u))
	}
	return &userGrpc.GetResponse{
		Users: usersPb,
	}, nil
}

func (s *userGrpcServer) GetById(ctx context.Context, req *userGrpc.GetByIdRequest) (*userGrpc.GetByIdResponse, error) {
	user, err := s.userService.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	userPb := converter.ToUserPbFromService(user)
	return &userGrpc.GetByIdResponse{
		User: userPb,
	}, nil
}

func main() {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", PORT))

	if err != nil {
		log.Fatalf("Failed to listen %v", PORT)
	}
	s := grpc.NewServer()
	cfg := config.NewConfig()
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.Db.DSN)
	userRepo := user.NewUserRepository(pool)
	userSrv := userService.NewUserService(userRepo)
	reflection.Register(s)
	userGrpc.RegisterUserV1Server(s, &userGrpcServer{userService: userSrv})
	log.Printf("Server listening on %s", l.Addr())

	if err = s.Serve(l); err != nil {
		log.Fatalf("Failed to serve: %s", err.Error())
	}
}
