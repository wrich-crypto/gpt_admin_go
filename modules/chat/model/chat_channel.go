package model

import (
	"github.com/cool-team-official/cool-admin-go/cool"
)

type ChatChannel struct {
	*cool.Model
	ChannelUUID string `gorm:"column:channel_uuid;size:100"`
	Title       string `gorm:"column:title;type:text"`
	UserID      int    `gorm:"column:user_id;not null"`
	Version     string `gorm:"column:version;size:10;not null"`
	Status      int    `gorm:"column:status;not null"`
}

const TableNameChatChannel = "chat_channels"

func (*ChatChannel) TableName() string {
	return TableNameChatChannel
}

func (*ChatChannel) GroupName() string {
	return "default"
}

func NewChatChannel() *ChatChannel {
	return &ChatChannel{
		Model: cool.NewModel(),
	}
}

// init 创建表
func init() {
	cool.CreateTable(&ChatChannel{})
}
