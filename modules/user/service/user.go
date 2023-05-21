package service

import (
	"gpt_admin_go/modules/user/model"

	"github.com/cool-team-official/cool-admin-go/cool"
)

type UserService struct {
	*cool.Service
}

func NewUserService() *UserService {
	return &UserService{
		&cool.Service{
			Model: model.NewUsers(),
		},
	}
}
