// user_balances.go
package model

import (
	"github.com/cool-team-official/cool-admin-go/cool"
)

type UserBalances struct {
	*cool.Model
	UserID         uint    `gorm:"column:user_id;not null"`
	TotalRecharge  float64 `gorm:"column:total_recharge;type:decimal(20,2);not null"`
	ConsumedAmount float64 `gorm:"column:consumed_amount;type:decimal(20,2);not null"`
}

const TableNameUserBalances = "user_balances"

func (*UserBalances) TableName() string {
	return TableNameUserBalances
}

func (*UserBalances) GroupName() string {
	return "default"
}

func NewUserBalances() *UserBalances {
	return &UserBalances{
		Model: cool.NewModel(),
	}
}

// init 创建表
func init() {
	cool.CreateTable(&UserBalances{})
}
