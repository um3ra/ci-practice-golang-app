package auth

import (
	authGrpc "github.com/um3ra/auth-microservice/pkg/auth_v1"
	"github.com/um3ra/auth-microservice/internal/service"
)


type AuthHandler struct {
	authService service.AuthService
	authGrpc.UnimplementedAuthV1Server
}

func NewAuthHandler(authService service.AuthService) *AuthHandler{
	return &AuthHandler{
		authService: authService,
	}
}