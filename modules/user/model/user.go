package model

import (
	"strconv"
	"strings"

	"github.com/cool-team-official/cool-admin-go/cool"
)

const (
	ROLE_NORMAL = 1
	ROLE_AGENT  = 2
	ROLE_ADMIN  = 3
)

type Users struct {
	*cool.Model
	Username           string `gorm:"column:username;size:50;unique;not null" json:"username"`
	Password           string `gorm:"column:password;size:255;not null" json:"-"`
	Email              string `gorm:"column:email;size:100" json:"email"`
	Phone              string `gorm:"column:phone;size:20" json:"phone"`
	BindPhone          int    `gorm:"column:bind_phone" json:"bind_phone"`
	Token              string `gorm:"column:token;size:100;index" json:"-"`
	InvitationCode     string `gorm:"column:invitation_code;size:20" json:"-"`
	InvitationUserId   string `gorm:"column:invitation_user_id;size:20" json:"invitation_user_id"`
	InvitationUserName string `gorm:"column:invitation_user_name;size:20" json:"invitation_user_name"`
	ReferralCode       string `gorm:"column:referral_code;size:20;unique;not null" json:"referral_code"`
	CardCount          int    `gorm:"column:card_count" json:"card_count"`
	Points             int    `gorm:"column:points" json:"points"`
	UsedPoints         int    `gorm:"column:used_points" json:"used_points"`
	RemainingPoints    int    `gorm:"column:remaining_points" json:"remaining_points"`
	Source             string `gorm:"column:source" json:"source"`
	Remarks            string `gorm:"column:remarks" json:"remarks"`
	Role               string `gorm:"column:role;default:1;comment:多个角色用逗号隔开" json:"role"`
	Status             int    `gorm:"column:status;default:1" json:"status"`
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

func (u *Users) HasRole(role int) bool {
	roles := strings.Split(u.Role, ",")
	for _, r := range roles {
		rInt, err := strconv.Atoi(r)
		if err != nil {
			// handle error here, possibly continue to next iteration
			continue
		}
		if rInt == role {
			return true
		}
	}
	return false
}

// init 创建表
func init() {
	cool.CreateTable(&Users{})
}

// 检查role是否存在
