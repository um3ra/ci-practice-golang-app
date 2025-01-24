package main

import (
	"context"
	"fmt"
	userGrpc "github.com/um3ra/auth-microservice/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:8888", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		log.Fatalf("Error description: %v", err)
	}

	client := userGrpc.NewUserV1Client(conn)
	ctx := context.Background()

	resp, err := client.Get(ctx, &userGrpc.GetRequest{Id: 1})

	if err != nil {
		log.Fatalf("failed to get user by id: %v", err)
	}
	fmt.Printf("Response: %#v\n", resp.User.GetName())
}
