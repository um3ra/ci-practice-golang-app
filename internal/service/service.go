package service

import (
	"context"
	"github.com/um3ra/auth-microservice/internal/model"
)

type UserService interface {
	GetAll(ctx context.Context) ([]model.User, error)
	GetById(ctx context.Context, id int64) (*model.User, error)
}

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, error)
	Register(ctx context.Context, user *model.User) (string, error)
}