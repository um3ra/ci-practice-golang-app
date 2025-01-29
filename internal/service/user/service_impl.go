package user

import (
	"context"
	"fmt"

	"github.com/um3ra/auth-microservice/internal/repository"

	"github.com/um3ra/auth-microservice/internal/model"
)

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) *userService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) GetAll(ctx context.Context) ([]model.User, error) {
	res, err := s.userRepository.GetAll(ctx)
	return res, err
}


func (s *userService) Create(ctx context.Context, user *model.User) (int64, error){
	id, err := s.userRepository.Create(ctx, user)
	fmt.Println(id, err)
	return id, err
}

func (s *userService)  GetById(ctx context.Context, id int64) (*model.User, error){
	res, err := s.userRepository.GetById(ctx, id)
	return res, err
}