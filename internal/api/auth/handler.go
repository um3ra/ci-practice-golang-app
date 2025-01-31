package auth

import (
	"context"
	"errors"

	"github.com/um3ra/auth-microservice/internal/model"
	"github.com/um3ra/auth-microservice/internal/service"
	authGrpc "github.com/um3ra/auth-microservice/pkg/auth_v1"
)

type AuthHandler struct {
	authService service.AuthService
	authGrpc.UnimplementedAuthV1Server
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(ctx context.Context, req *authGrpc.LoginRequest) (*authGrpc.LoginResponse, error) {
	token, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil{
		return nil, err
	}
	return &authGrpc.LoginResponse{
		Token: token,
	}, nil
}

func (h *AuthHandler) Register(ctx context.Context, req *authGrpc.RegisterRequest) (*authGrpc.RegisterResponse, error) {
	if req.GetPassword() != req.ConfirmPassword {
		return nil, errors.New(PasswordMismatch)
	}
	newUser := model.User{
		Name: req.GetName(),
		Password: req.GetPassword(),
		Email: req.GetEmail(),
	}
	token, err := h.authService.Register(ctx, &newUser)

	if err != nil {
		return nil, err
	}
	return &authGrpc.RegisterResponse{
		Token: token,
	}, nil
}
