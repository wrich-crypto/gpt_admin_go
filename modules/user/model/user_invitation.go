package model

// user_invitations.go

import (
	"github.com/cool-team-official/cool-admin-go/cool"
)

type UserInvitations struct {
	*cool.Model
	InviterID     uint    `gorm:"column:inviter_id;not null"`
	InviterReward float64 `gorm:"column:inviter_reward;type:decimal(10,2);not null"`
	InviteeID     uint    `gorm:"column:invitee_id;not null"`
	InviteeReward float64 `gorm:"column:invitee_reward;type:decimal(10,2);not null"`
}

const TableNameUserInvitations = "user_invitations"

func (*UserInvitations) TableName() string {
	return TableNameUserInvitations
}

func (*UserInvitations) GroupName() string {
	return "default"
}

func NewUserInvitations() *UserInvitations {
	return &UserInvitations{
		Model: cool.NewModel(),
	}
}

// init 创建表
func init() {
	cool.CreateTable(&UserInvitations{})
}
