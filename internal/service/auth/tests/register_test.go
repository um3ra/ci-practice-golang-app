package tests

import (
	"context"
	"fmt"
	"testing"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"github.com/um3ra/auth-microservice/internal/model"
	"github.com/um3ra/auth-microservice/internal/repository/mocks"
	"github.com/stretchr/testify/mock"
	
	authService "github.com/um3ra/auth-microservice/internal/service/auth"
)

func TestRegister(t *testing.T) {
	type mockFun func(mock *mocks.UserRepository)
	type args struct {
		Ctx  context.Context
		User *model.User
	}
	repoErr := fmt.Errorf("repo error")

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
			testName: "success register test",
			args: args{
				Ctx:  ctx,
				User: successUser,
			},
			want: id,
			err:  repoErr,
			mockFun: func(repoMock *mocks.UserRepository) {
				repoMock.On("GetByEmail", ctx, email).Return(nil, repoErr).Once()
				repoMock.On("Create", ctx, mock.AnythingOfType("*model.User")).Return(successUser.Id, nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			authMock := mocks.NewUserRepository(t)
			tt.mockFun(authMock)
			authSrv := authService.NewAuthService(authMock)
			newId, err := authSrv.Register(tt.args.Ctx, tt.args.User)
			require.NotEmpty(t, newId)
			require.Equal(t, nil, err)
		})
	}
}
