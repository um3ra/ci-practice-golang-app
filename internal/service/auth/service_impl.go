package auth

import (
	"context"
	"github.com/um3ra/auth-microservice/internal/model"
	"github.com/um3ra/auth-microservice/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"errors"
)

type authService struct {
	userRepository repository.UserRepository
}

func NewAuthService(userRepository repository.UserRepository) *authService {
	return &authService{
		userRepository: userRepository,
	}
}

func (a *authService) Login(ctx context.Context, email, password string) (int64, error) {
	exsUser, err := a.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(exsUser.Password), []byte(password))
	if err != nil {
		return 0, err
	}

	return exsUser.Id, err
}

func (a *authService) Register(ctx context.Context, user *model.User) (int64, error) {
	exsUser, _ := a.userRepository.GetByEmail(ctx, user.Email)
	if exsUser != nil {
		return 0, errors.New(UserExistsError)
	}
	_, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	
	newUser := &model.User{
		Name:     user.Name,
		Password: user.Password,
		Email:    user.Email,
	}
	id, err := a.userRepository.Create(ctx, newUser)
	if err != nil {
		return 0, err
	}

	return id, nil
}
