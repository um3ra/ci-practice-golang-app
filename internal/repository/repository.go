package repository

import (
	"context"
	"github.com/um3ra/auth-microservice/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) (int64, error)
	GetAll(ctx context.Context) ([]model.User, error)
	GetById(ctx context.Context, id int64) (*model.User, error)
}
