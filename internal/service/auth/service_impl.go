package auth

import (
	"context"
	"errors"

	"github.com/um3ra/auth-microservice/internal/model"
	"github.com/um3ra/auth-microservice/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepository repository.UserRepository
}

func (a *AuthService) Login(ctx context.Context, email, password string) (int64, error) {
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



func (a *AuthService) Register(ctx context.Context, user *model.User) (int64, error){
	exsUser, err := a.userRepository.GetByEmail(ctx, user.Email)
	if err != nil {
		return 0, err
	}

	if exsUser != nil {
		return 0, errors.New(userExistsError)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	newUser := &model.User{
		Name: user.Name,
		Password: string(hashedPassword),
		Email: user.Email,
	}


	id, err := a.userRepository.Create(ctx, newUser)
	if err != nil {
		return 0, err
	}
	return id, nil
}