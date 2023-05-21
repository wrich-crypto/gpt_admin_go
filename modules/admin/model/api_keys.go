package model

import (
	"time"

	"github.com/cool-team-official/cool-admin-go/cool"
)

type ApiKeys struct {
	*cool.Model
	ApiKey      string    `gorm:"column:api_key;size:255;not null;comment:api_key" json:"api_key"`
	ModelType   string    `gorm:"column:model_type;size:20;not null;comment:3.5或者4" json:"model_type"`
	Status      int       `gorm:"column:status;size:20;not null;default:1;comment:1正常 2禁用" json:"status"`
	TotalTokens int       `gorm:"column:total_tokens;comment:总消耗" json:"total_tokens"`
	DailyTokens int       `gorm:"column:daily_tokens;comment:当天消耗" json:"daily_tokens"`
	ExpiryDate  time.Time `gorm:"column:expiry_date;comment:过期时间" json:"expiry_date"`
	Supplier    string    `gorm:"column:supplier;size:50;default:'openai';comment:供应商默认openai 目前支持openai和uchat" json:"supplier"`
}

const TableNameApiKeys = "api_keys"

func (*ApiKeys) TableName() string {
	return TableNameApiKeys
}

func (*ApiKeys) GroupName() string {
	return "default"
}

func NewApiKeys() *ApiKeys {
	return &ApiKeys{
		Model: cool.NewModel(),
	}
}

// init 创建表
func init() {
	cool.CreateTable(&ApiKeys{})
}
