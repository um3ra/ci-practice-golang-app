package converter

import (
	model "github.com/um3ra/auth-microservice/internal/model"
	repoMod "github.com/um3ra/auth-microservice/internal/repository/user/model"
)

func ToUserFromRepo(user *repoMod.User) *model.User {
	return &model.User{
		ID:        user.Id,
		Name:      user.Name,
		Password:  user.Password,
		Email:     user.Email,
		UpdatedAt: user.UpdatedAt,
		CreatedAt: user.CreatedAt,
	}
}
