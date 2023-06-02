package model

import (
	"time"

	"github.com/cool-team-official/cool-admin-go/cool"
)

type RechargeCards struct {
	*cool.Model
	CardAccount    string    `gorm:"column:card_account;size:20;unique;not null" json:"card_account"`
	CardPassword   string    `gorm:"column:card_password;size:20;not null" json:"card_password"`
	CreateTime     time.Time `gorm:"column:create_time" json:"create_time"`
	ExpireTime     time.Time `gorm:"column:expire_time" json:"expire_time"`
	RechargeAmount float64   `gorm:"column:recharge_amount;type:decimal(20,2);not null"`
	UsedPoints     int       `gorm:"column:used_points" json:"used_points"`
	BoundUser      string    `gorm:"column:bound_user;size:100" json:"bound_user"`
	CreateUser     string    `gorm:"column:create_user;size:100" json:"create_user"`
	CreateUserId   int       `gorm:"column:create_user_id" json:"create_user_id"`
	RechargeType   string    `gorm:"column:recharge_type;size:50" json:"recharge_type"`
	Remark         string    `gorm:"column:remark;size:255" json:"remark"`
	Status         int       `gorm:"column:status;default:1;not null" json:"status"`
	RechargeTime   time.Time `gorm:"column:recharge_time" json:"recharge_time"`
}

const TableNameRechargeCards = "recharge_cards"

func (*RechargeCards) TableName() string {
	return TableNameRechargeCards
}

func (*RechargeCards) GroupName() string {
	return "default"
}

func NewRechargeCards() *RechargeCards {
	return &RechargeCards{
		Model: cool.NewModel(),
	}
}

// init 创建表
func init() {
	cool.CreateTable(&RechargeCards{})
}
