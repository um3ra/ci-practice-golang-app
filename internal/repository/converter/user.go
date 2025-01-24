package converter

import (
	repoMod "github.com/um3ra/auth-microservice/internal/repository/model"
	model "github.com/um3ra/auth-microservice/internal/model"
)

func ToUserFromRepo(user *repoMod.User) *model.User {
	return &model.User{
		Id: user.Id,
		Name: user.Name,
		Password: user.Password,
		Email: user.Email,
		UpdatedAt: user.UpdatedAt,
		CreatedAt: user.CreatedAt,
	}
}