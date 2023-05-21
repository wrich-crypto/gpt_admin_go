package model

import (
	"github.com/cool-team-official/cool-admin-go/cool"
)

type UserRecharges struct {
	*cool.Model
	UserID         uint    `gorm:"column:user_id;not null"`
	PayAmount      float64 `gorm:"column:pay_amount;type:decimal(20,2)"`
	Amount         float64 `gorm:"column:amount;type:decimal(20,2);not null"`
	RechargeMethod int     `gorm:"column:recharge_method;not null"`
	Status         int     `gorm:"column:status;not null"`
}

const TableNameUserRecharges = "user_recharges"

func (*UserRecharges) TableName() string {
	return TableNameUserRecharges
}

func (*UserRecharges) GroupName() string {
	return "default"
}

func NewUserRecharges() *UserRecharges {
	return &UserRecharges{
		Model: cool.NewModel(),
	}
}

// init 创建表
func init() {
	cool.CreateTable(&UserRecharges{})
}
