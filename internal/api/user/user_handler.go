package user

import (
	"context"
	"github.com/um3ra/auth-microservice/internal/converter"
	"github.com/um3ra/auth-microservice/internal/service"
	userGrpc "github.com/um3ra/auth-microservice/pkg/user_v1"
)

type UserHandler struct {
	userService service.UserService
	userGrpc.UnimplementedUserV1Server
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (s *UserHandler) Get(ctx context.Context, data *userGrpc.GetRequest) (*userGrpc.GetResponse, error) {
	users, err := s.userService.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	var usersPb []*userGrpc.User

	for _, u := range users {
		usersPb = append(usersPb, converter.ToUserPbFromService(&u))
	}
	return &userGrpc.GetResponse{
		Users: usersPb,
	}, nil
}

func (s *UserHandler) GetById(ctx context.Context, req *userGrpc.GetByIdRequest) (*userGrpc.GetByIdResponse, error) {
	user, err := s.userService.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	userPb := converter.ToUserPbFromService(user)
	return &userGrpc.GetByIdResponse{
		User: userPb,
	}, nil
}
