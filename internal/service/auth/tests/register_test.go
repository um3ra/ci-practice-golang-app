package tests

import (
	"context"
	"errors"
	"testing"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/um3ra/auth-microservice/internal/model"
	"github.com/um3ra/auth-microservice/internal/repository/mocks"
	"github.com/um3ra/auth-microservice/pkg/jwt"
	jwtMock "github.com/um3ra/auth-microservice/pkg/jwt/mocks"

	authService "github.com/um3ra/auth-microservice/internal/service/auth"
)

func TestRegister(t *testing.T) {
	type mockBehavior func(userRepoMock *mocks.UserRepository, jwtServiceMock *jwtMock.JwtService)
	type args struct {
		Ctx  context.Context
		User *model.User
	}
	repoErr := errors.New("repo mock error")

	var (
		ctx      = context.Background()
		name     = fake.Name()
		email    = fake.Email()
		password = fake.Animal()
		id       = fake.Int64()
		secret   = fake.BeerName()

		userMock = &model.User{
			Id:       id,
			Name:     name,
			Password: password,
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
			testName: "register test success case",
			args: args{
				Ctx:  ctx,
				User: userMock,
			},
			want: secret,
			err:  nil,
			mockBehavior: func(userRepoMock *mocks.UserRepository, jwtServiceMock *jwtMock.JwtService) {
				userRepoMock.EXPECT().GetByEmail(ctx, email).Return(nil, repoErr).Once()
				userRepoMock.EXPECT().Create(ctx, mock.AnythingOfType("*model.User")).Return(id, nil).Once()
				jwtServiceMock.EXPECT().Create(jwt.JwtPayload{Email: userMock.Email}).Return(secret, nil).Once()
			},
		},

		{
			testName: "register test fail case (user exists)",
			args: args{
				Ctx:  ctx,
				User: userMock,
			},
			want: "",
			err:  errors.New(authService.UserExistsError),
			mockBehavior: func(userRepoMock *mocks.UserRepository, jwtServiceMocj *jwtMock.JwtService) {
				userRepoMock.EXPECT().GetByEmail(ctx, email).Return(userMock, nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			userRepoMock := mocks.NewUserRepository(t)
			jwtSrvMock := jwtMock.NewJwtService(t)
			tt.mockBehavior(userRepoMock, jwtSrvMock)
			authSrv := authService.NewAuthService(userRepoMock, jwtSrvMock)
			res, err := authSrv.Register(tt.args.Ctx, tt.args.User)
			require.Equal(t, tt.want, res)
			require.Equal(t, tt.err, err)
		})
	}
}
