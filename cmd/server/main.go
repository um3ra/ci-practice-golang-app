package main

import (
	"context"
	"log"

	"github.com/um3ra/auth-microservice/internal/app"
)

func main() {

	ctx := context.Background()

	app, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("Failed to init app: %s\n", err.Error())
	}
	err = app.Run()
	if err != nil {
		log.Fatalf("Failed to run app: %s\n", err.Error())
	}
}
