package model

import (
	"time"

	"github.com/cool-team-official/cool-admin-go/cool"
)

type UserOrders struct {
	*cool.Model
	TradeNo   string    `gorm:"column:trade_no;type:varchar(100);not null"`
	UserID    int       `gorm:"column:user_id;not null"`
	OrderType int       `gorm:"column:order_type;not null;default:1"`
	Amount    float64   `gorm:"column:amount;type:decimal(10,2);not null"`
	Status    int       `gorm:"column:status;not null;default:1"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;not null"`
}

const TableNameUserOrders = "user_orders"

func (*UserOrders) TableName() string {
	return TableNameUserOrders
}

func (*UserOrders) GroupName() string {
	return "default"
}

func NewUserOrders() *UserOrders {
	return &UserOrders{
		Model: cool.NewModel(),
	}
}

// init 创建表
func init() {
	cool.CreateTable(&UserOrders{})
}
