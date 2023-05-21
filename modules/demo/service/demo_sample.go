package service

import (
	"gpt_admin_go/modules/demo/model"

	"github.com/cool-team-official/cool-admin-go/cool"
)

type DemoSampleService struct {
	*cool.Service
}

func NewDemoSampleService() *DemoSampleService {
	return &DemoSampleService{
		&cool.Service{
			Model: model.NewDemoSample(),
		},
	}
}
