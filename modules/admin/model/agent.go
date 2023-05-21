package model

import (
	"github.com/cool-team-official/cool-admin-go/cool"
)

type Agent struct {
	*cool.Model
	Title                string `gorm:"column:title;size:100;comment:标题" json:"title"`
	Slogan               string `gorm:"column:slogan;size:100;comment:广告语" json:"slogan"`
	AIName               string `gorm:"column:ai_name;size:50;comment:AI名字" json:"ai_name"`
	WeChatQRCode         string `gorm:"column:wechat_qr_code;size:255;comment:微信二维码" json:"wechat_qr_code"`
	CustomerServicePhone string `gorm:"column:customer_service_phone;size:20;comment:客服电话" json:"customer_service_phone"`
	Domain               string `gorm:"column:domain;size:100;comment:域名多个用逗号隔开" json:"domain"`
}

const TableNameAgent = "agent"

func (*Agent) TableName() string {
	return TableNameAgent
}

func (*Agent) GroupName() string {
	return "default"
}

func NewAgent() *Agent {
	return &Agent{
		Model: cool.NewModel(),
	}
}

// init 创建表
func init() {
	cool.CreateTable(&Agent{})
}
