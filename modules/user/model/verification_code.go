package model

import (
	"time"

	"github.com/cool-team-official/cool-admin-go/cool"
)

type VerificationCodes struct {
	*cool.Model
	Username  string    `gorm:"column:username;size:100"`
	Email     string    `gorm:"column:email;size:100"`
	Phone     string    `gorm:"column:phone;size:100"`
	CodeType  string    `gorm:"column:code_type;size:10;not null"`
	Code      string    `gorm:"column:code;size:10;not null"`
	ExpiredAt time.Time `gorm:"column:expired_at"`
}

const TableNameVerificationCodes = "verification_code"

func (*VerificationCodes) TableName() string {
	return TableNameVerificationCodes
}

func (*VerificationCodes) GroupName() string {
	return "default"
}

func NewVerificationCodes() *VerificationCodes {
	return &VerificationCodes{
		Model: cool.NewModel(),
	}
}

// init 创建表
func init() {
	cool.CreateTable(&VerificationCodes{})
}
