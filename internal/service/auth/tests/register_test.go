package tests

import (
	"context"
	"errors"
	"fmt"
	"testing"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/um3ra/auth-microservice/internal/model"
	"github.com/um3ra/auth-microservice/internal/repository/mocks"

	authService "github.com/um3ra/auth-microservice/internal/service/auth"
)

func TestRegister(t *testing.T) {
	type mockFun func(mock *mocks.UserRepository)
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

		successUser = &model.User{
			Id:       id,
			Name:     name,
			Password: password,
			Email:    email,
		}
	)

	tests := []struct {
		testName string
		args     args
		want     int64
		err      error
		mockFun  mockFun
	}{
		{
			testName: "register test success case",
			args: args{
				Ctx:  ctx,
				User: successUser,
			},
			want: id,
			err:  nil,
			mockFun: func(repoMock *mocks.UserRepository) {
				repoMock.On("GetByEmail", ctx, email).Return(nil, repoErr).Once()
				repoMock.On("Create", ctx, mock.AnythingOfType("*model.User")).Return(id, nil).Once()
			},
		},
		{
			testName: "register test fail case (user exists)",
			args: args{
				Ctx:  ctx,
				User: successUser,
			},
			want: 0,
			err:  errors.New(authService.UserExistsError),
			mockFun: func(repoMock *mocks.UserRepository) {
				repoMock.On("GetByEmail", ctx, successUser.Email).Return(successUser, nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			authMock := mocks.NewUserRepository(t)
			tt.mockFun(authMock)
			authSrv := authService.NewAuthService(authMock)
			res, err := authSrv.Register(tt.args.Ctx, tt.args.User)
			fmt.Println(res)
			require.Equal(t, tt.want, res)
			require.Equal(t, tt.err, err)
		})
	}
}
