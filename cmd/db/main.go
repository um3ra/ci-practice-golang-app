package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"encoding/json"
	pg "github.com/jackc/pgx/v5"
	r "github.com/um3ra/auth-microservice/cmd/db/repositories"
	"github.com/um3ra/auth-microservice/config"
)

func main() {
	cfg := config.NewConfig()
	ctx := context.Background()
	conn, err := pg.Connect(ctx, cfg.Db.DSN)
	if err != nil {
		panic(err)
	}
	repo := r.NewRepository(r.RepositoryDeps{Db: conn, Context: ctx})
	mux := http.NewServeMux()

	_ = repo
	mux.HandleFunc("POST /user/create", func(w http.ResponseWriter, r *http.Request) {
		res, err := repo.Save()

		if err != nil {
			fmt.Println(err.Error())
		}

		if res {
			fmt.Fprintf(w, "User was created successfully\n")
		}

		// isCreated, err := repo.Save()
		// if err != nil {
		// 	fmt.Println(err.Error())
		// }

		// if isCreated {
		// 	w.WriteHeader(http.StatusCreated)
		// 	fmt.Fprintf(w, "User was created successfully\n")
		// }
	})

	mux.HandleFunc("GET /user", func(w http.ResponseWriter, r *http.Request) {
		users, err := repo.GetAll()

		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content Type", "application/json")
		json.NewEncoder(w).Encode(users)
	})
	if err := http.ListenAndServe(":8888", mux); err != nil {
		log.Fatal(err)
	}

}
