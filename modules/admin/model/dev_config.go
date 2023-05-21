package model

import (
	"github.com/cool-team-official/cool-admin-go/cool"
)

type DevConfig struct {
	*cool.Model
	ModelType        string `gorm:"column:model_type;size:20;default:'4'" json:"model_type"`
	Supplier         string `gorm:"column:supplier;size:50;default:'openai'" json:"supplier"`
	FreeTokenCount   int    `gorm:"column:free_token_count;default:0" json:"free_token_count"`
	InviteTokenCount int    `gorm:"column:invite_token_count;default:0" json:"invite_token_count"`
}

const TableNameDevConfig = "dev_config"

func (*DevConfig) TableName() string {
	return TableNameDevConfig
}

func (*DevConfig) GroupName() string {
	return "default"
}

func NewDevConfig() *DevConfig {
	return &DevConfig{
		Model: cool.NewModel(),
	}
}

// init 创建表
func init() {
	cool.CreateTable(&DevConfig{})
}
