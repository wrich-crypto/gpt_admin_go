package service

import (
	"gpt_admin_go/modules/user/model"

	"github.com/cool-team-official/cool-admin-go/cool"
)

type RechargeCardService struct {
	*cool.Service
}

func NewRechargeCardService() *RechargeCardService {
	return &RechargeCardService{
		&cool.Service{
			Model: model.NewRechargeCards(),
		},
	}
}
