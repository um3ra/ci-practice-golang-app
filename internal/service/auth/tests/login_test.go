package tests

import (
	"context"
	"errors"
	"log"
	"testing"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"github.com/um3ra/auth-microservice/internal/model"
	"github.com/um3ra/auth-microservice/internal/repository/mocks"
	authService "github.com/um3ra/auth-microservice/internal/service/auth"
	"github.com/um3ra/auth-microservice/pkg/jwt"
	jwtMock "github.com/um3ra/auth-microservice/pkg/jwt/mocks"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	type mockBehavior func(userRepoMock *mocks.UserRepository, jwtServiceMock *jwtMock.JwtService)
	type args struct {
		Ctx      context.Context
		Email    string
		Password string
	}
	password := fake.Animal()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	repoErr := errors.New("repo mock error")

	var (
		ctx      = context.Background()
		email    = fake.Email()
		secret   = fake.BeerName()
		userMock = &model.User{
			Password: string(hashedPassword),
			Email:    email,
		}
	)

	tests := []struct {
		testName     string
		args         args
		want         string
		err          error
		mockBehavior mockBehavior
	}{
		{
			testName: "login success case",
			args:     args{Ctx: ctx, Email: email, Password: password},
			want:     secret,
			err:      nil,
			mockBehavior: func(userRepoMock *mocks.UserRepository, jwtServiceMock *jwtMock.JwtService) {
				userRepoMock.On("GetByEmail", ctx, email).Return(userMock, nil).Once()
				jwtServiceMock.On("Create", jwt.JwtPayload{Email: email}).Return(secret, nil).Once()
			},
		},
		{
			testName: "login fail case (incorrect email)",
			args:     args{Ctx: ctx, Email: email, Password: password},
			want:     "",
			err:      errors.New(authService.IncorrectEmailOrPassword),
			mockBehavior: func(userRepoMock *mocks.UserRepository, jwtServiceMock *jwtMock.JwtService) {
				userRepoMock.On("GetByEmail", ctx, email).Return(nil, repoErr).Once()
			},
		},

		{
			testName: "login fail case (incorrect password)",
			args:     args{Ctx: ctx, Email: email, Password: "incorrect password"},
			want:     "",
			err:      errors.New(authService.IncorrectEmailOrPassword),
			mockBehavior: func(userRepoMock *mocks.UserRepository, jwtServiceMock *jwtMock.JwtService) {
				userRepoMock.On("GetByEmail", ctx, email).Return(nil, repoErr).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			userRepoMock := mocks.NewUserRepository(t)
			jwtSrvMock := jwtMock.NewJwtService(t)
			tt.mockBehavior(userRepoMock, jwtSrvMock)
			authSrv := authService.NewAuthService(userRepoMock, jwtSrvMock)
			res, err := authSrv.Login(tt.args.Ctx, tt.args.Email, tt.args.Password)
			require.Equal(t, tt.want, res)
			require.Equal(t, tt.err, err)
		})
	}
}
