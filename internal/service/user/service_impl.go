package user

import (
	"context"
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

func (s *userService)  GetById(ctx context.Context, id int64) (*model.User, error){
	res, err := s.userRepository.GetById(ctx, id)
	return res, err
}