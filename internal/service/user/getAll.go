package user

import (
	"context"

	"github.com/um3ra/auth-microservice/internal/model"
)

func (s *userService) GetAll(ctx context.Context) ([]model.User, error) {
	res, err := s.userRepository.GetAll(ctx)
	return res, err
}
