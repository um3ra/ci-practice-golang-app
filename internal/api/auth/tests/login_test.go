package tests

import (
	"context"
	"errors"
	"testing"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"github.com/um3ra/auth-microservice/internal/api/auth"
	authMockSrv "github.com/um3ra/auth-microservice/internal/service/mocks"
	authGrpc "github.com/um3ra/auth-microservice/pkg/auth_v1"
)

func TestHandler_login(t *testing.T) {
	type mockBehavior func(authService *authMockSrv.AuthService)
	type args struct {
		ctx context.Context
		req *authGrpc.LoginRequest
	}
	var (
		mockErr = errors.New("mock error")
		ctx      = context.Background()
		email    = fake.Email()
		password = fake.Animal()
		token    = fake.BeerName()
	)

	testTable := []struct {
		name         string
		args         args
		want         *authGrpc.LoginResponse
		err          error
		mockBehavior mockBehavior
	}{
		{
			name: "handler login success",
			args: args{
				ctx: ctx,
				req: &authGrpc.LoginRequest{
					Email:    email,
					Password: password,
				},
			},
			want: &authGrpc.LoginResponse{
				Token: token,
			},
			err: nil,
			mockBehavior: func(authService *authMockSrv.AuthService) {
				authService.EXPECT().Login(ctx, email, password).Return(token, nil).Once()
			},
		},

		{
			name: "handler login fail case",
			args: args{
				ctx: ctx,
				req: &authGrpc.LoginRequest{
					Email:    email,
					Password: password,
				},
			},
			want: nil,
			err: mockErr,
			mockBehavior: func(authService *authMockSrv.AuthService) {
				authService.EXPECT().Login(ctx, email, password).Return("", mockErr).Once()
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			authService := authMockSrv.NewAuthService(t)
			tt.mockBehavior(authService)
			handler := auth.NewAuthHandler(authService)
			res, err := handler.Login(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.want, res)
			require.Equal(t, tt.err, err)
		})
	}
}
