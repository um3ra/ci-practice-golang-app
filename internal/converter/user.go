package converter

import (
	"github.com/um3ra/auth-microservice/internal/model"
	pb "github.com/um3ra/auth-microservice/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToUserPbFromService(user *model.User) *pb.User {
	var updatedAt *timestamppb.Timestamp
	createdAt := timestamppb.New(user.CreatedAt)
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return &pb.User{
		Id: user.Id,
		Name: user.Name,
		Email: user.Email,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func ToUserFromPb(user *pb.User) *model.User {
	return &model.User{
		Name: user.Name,
		Email: user.Email,
	}
}