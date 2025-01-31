package tests

import (
	"context"
	"testing"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"github.com/um3ra/auth-microservice/internal/api/auth"
	"github.com/um3ra/auth-microservice/internal/model"
	authMockSrv "github.com/um3ra/auth-microservice/internal/service/mocks"
	authGrpc "github.com/um3ra/auth-microservice/pkg/auth_v1"
	"errors"
)

func TestHandler_register(t *testing.T) {

	type args struct {
		ctx context.Context
		req *authGrpc.RegisterRequest
	}
	type mockBehavior func (service *authMockSrv.AuthService)

	var (
		ctx      = context.Background()
		name = fake.Name()
		email    = fake.Email()
		password = fake.Animal()
		token    = fake.BeerName()
	)
	user := model.User{
		Name: name,
		Email: email,
		Password: password,
	}

	testTable := []struct{
		name string
		want *authGrpc.RegisterResponse
		args args
		err error
		mockBehavior mockBehavior
	}{
		{
			name: "handler register success case",
			want: &authGrpc.RegisterResponse{
				Token: token,
			},
			args: args{
				ctx: ctx,
				req: &authGrpc.RegisterRequest{
					Name: name,
					Email: email,
					Password: password,
					ConfirmPassword: password,
				},
			},
			err: nil,
			mockBehavior: func(service *authMockSrv.AuthService) {
				service.EXPECT().Register(ctx, &user).Return(token, nil).Once()
			},
		},


		{
			name: "handler register fail case (password mismatch)",
			want: nil,
			args: args{
				ctx: ctx,
				req: &authGrpc.RegisterRequest{
					Name: name,
					Email: email,
					Password: password,
					ConfirmPassword: "random",
				},
			},
			err: errors.New(auth.PasswordMismatch),
			mockBehavior: func(service *authMockSrv.AuthService) {
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			authService := authMockSrv.NewAuthService(t)
			tt.mockBehavior(authService)
			handler := auth.NewAuthHandler(authService)
			res, err := handler.Register(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.want, res)
			require.Equal(t, tt.err, err)
		})
	}
}