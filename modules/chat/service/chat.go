package service

import (
	"gpt_admin_go/modules/chat/model"

	"github.com/cool-team-official/cool-admin-go/cool"
)

type ChatService struct {
	*cool.Service
}

func NewChatChannelService() *ChatService {
	return &ChatService{
		&cool.Service{
			Model: model.NewChatChannel(),
		},
	}
}

func NewChatMessageService() *ChatService {
	return &ChatService{
		&cool.Service{
			Model: model.NewChatMessage(),
		},
	}
}
