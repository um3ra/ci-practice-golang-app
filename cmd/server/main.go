package main

import (
	"context"
	"fmt"
	"log"
	"net"

	userGrpc "github.com/um3ra/auth-microservice/pkg/user_v1"
	r "github.com/um3ra/auth-microservice/cmd/db/repositories"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pg "github.com/jackc/pgx/v5"
	"github.com/um3ra/auth-microservice/config"


)

const PORT = 8888

type userGrpcServer struct {
	userGrpc.UnimplementedUserV1Server

	UserRepository *r.Repository
}


// func (s *userGrpcServer) Get(ctx context.Context, method *userGrpc.GetRequest) (*userGrpc.GetResponse, error){


// 	id := method.GetId()
// 	if id == 0 {
// 		return nil, nil
// 	}

// 	// if id != 2 {
// 	// 	return nil, errors.New(fmt.Sprintf("User was not found with id: %d", id))
// 	// }
// 	user, err := s.UserRepository.GetById(int(id))
// 	if err != nil {
// 		fmt.Println(err.Error())

// 		return nil, err
// 	}


// 	if user == nil {
// 		return nil, errors.New("User was not found!")
// 	}

// 	// newUser := userGrpc.GetResponse{
// 	// 	Id: int64(user.Id),
// 	// 	Name: user.Name,
// 	// 	Email: user.Email,
// 	// 	Role: userGrpc.Role_Regular,
// 	// }
// 	return &newUser, nil
// }

// func (s * userGrpcServer) Create(ctx context.Context, method *userGrpc.CreateRequest)(*userGrpc.CreateResponse, error){

// }

func main() {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", PORT))

	if err != nil {
		log.Fatalf("Failed to listen %v", PORT)
	}
	s := grpc.NewServer()
	cfg := config.NewConfig()
	ctx := context.Background()
	conn, err := pg.Connect(ctx, cfg.Db.DSN)

	repo := r.NewRepository(r.RepositoryDeps{Db: conn, Context: ctx})
	reflection.Register(s)
	userGrpc.RegisterUserV1Server(s, &userGrpcServer{UserRepository: repo})
	log.Printf("Server listening on %s", l.Addr())

	if err = s.Serve(l); err != nil {
		log.Fatalf("Failed to serve: %s", err.Error())
	}
}
