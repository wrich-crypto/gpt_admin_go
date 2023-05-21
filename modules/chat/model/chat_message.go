package model

import (
	"github.com/cool-team-official/cool-admin-go/cool"
)

type ChatMessage struct {
	*cool.Model
	UserID         int     `gorm:"column:user_id;not null"`
	ChannelID      int     `gorm:"column:channel_id;not null"`
	MessageID      string  `gorm:"column:message_id;size:100"`
	StreamID       string  `gorm:"column:stream_id;size:100;not null"`
	Question       string  `gorm:"column:question;type:text;not null"`
	Answer         string  `gorm:"column:answer;type:text"`
	UsingContext   int     `gorm:"column:using_context;not null"`
	TokensConsumed float64 `gorm:"column:tokens_consumed;type:decimal(10,2)"`
	ApiKey         string  `gorm:"column:api_key;size:255"`
	Version        string  `gorm:"column:version;size:10;not null"`
}

const TableNameChatMessage = "chat_messages"

func (*ChatMessage) TableName() string {
	return TableNameChatMessage
}

func (*ChatMessage) GroupName() string {
	return "default"
}

func NewChatMessage() *ChatMessage {
	return &ChatMessage{
		Model: cool.NewModel(),
	}
}

// init 创建表
func init() {
	cool.CreateTable(&ChatMessage{})
}
