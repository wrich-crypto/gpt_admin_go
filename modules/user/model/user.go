package model

import (
	"github.com/cool-team-official/cool-admin-go/cool"
)

type Users struct {
	*cool.Model
	Username        string `gorm:"column:username;size:50;unique;not null" json:"username"`
	Password        string `gorm:"column:password;size:255;not null" json:"-"`
	Email           string `gorm:"column:email;size:100" json:"email"`
	Phone           string `gorm:"column:phone;size:20" json:"phone"`
	BindPhone       int    `gorm:"column:bind_phone" json:"bind_phone"`
	Token           string `gorm:"column:token;size:100;index" json:"-"`
	InvitationCode  string `gorm:"column:invitation_code;size:20" json:"-"`
	ReferralCode    string `gorm:"column:referral_code;size:20;unique;not null" json:"referral_code"`
	CardCount       int    `gorm:"column:card_count" json:"card_count"`
	Points          int    `gorm:"column:points" json:"points"`
	UsedPoints      int    `gorm:"column:used_points" json:"used_points"`
	RemainingPoints int    `gorm:"column:remaining_points" json:"remaining_points"`
	Source          string `gorm:"column:source" json:"source"`
	Remarks         string `gorm:"column:remarks" json:"remarks"`
	Role            int    `gorm:"column:role;default:1" json:"role"`
	Status          int    `gorm:"column:status;default:1" json:"status"`
}

const TableNameUsers = "users"

func (*Users) TableName() string {
	return TableNameUsers
}

func (*Users) GroupName() string {
	return "default"
}

func NewUsers() *Users {
	return &Users{
		Model: cool.NewModel(),
	}
}

// init 创建表
func init() {
	cool.CreateTable(&Users{})
}
