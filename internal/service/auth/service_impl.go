package auth

import (
	"context"
	"errors"

	"github.com/um3ra/auth-microservice/internal/model"
	"github.com/um3ra/auth-microservice/internal/repository"
	"github.com/um3ra/auth-microservice/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepository repository.UserRepository
	jwtService     jwt.JwtService
}

func NewAuthService(userRepository repository.UserRepository, jwt jwt.JwtService) *authService {
	return &authService{
		userRepository: userRepository,
		jwtService:     jwt,
	}
}

func (a *authService) Login(ctx context.Context, email, password string) (string, error) {
	exsUser, err := a.userRepository.GetByEmail(ctx, email)
	if err != nil || exsUser == nil {
		return "", errors.New(IncorrectEmailOrPassword)
	}

	err = bcrypt.CompareHashAndPassword([]byte(exsUser.Password), []byte(password))
	if err != nil {
		return "", errors.New(IncorrectEmailOrPassword)
	}

	token, err := a.jwtService.Create(jwt.JwtPayload{Email: exsUser.Email})
	if err != nil {
		return "", err
	}
	return token, err
}

func (a *authService) Register(ctx context.Context, user *model.User) (string, error) {
	exsUser, _ := a.userRepository.GetByEmail(ctx, user.Email)
	if exsUser != nil {
		return "", errors.New(UserExistsError)
	}
	_, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	newUser := &model.User{
		Name:     user.Name,
		Password: user.Password,
		Email:    user.Email,
	}
	_, err = a.userRepository.Create(ctx, newUser)
	if err != nil {
		return "", err
	}
	token, err := a.jwtService.Create(jwt.JwtPayload{Email: newUser.Email})
	if err != nil {
		return "", err
	}
	return token, nil
}
