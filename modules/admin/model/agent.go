package model

import (
	"github.com/cool-team-official/cool-admin-go/cool"
	"github.com/gogf/gf/v2/frame/g"
)

type Agent struct {
	*cool.Model
	Title                string  `gorm:"column:title;size:100;comment:标题" json:"title"`
	Slogan               string  `gorm:"column:slogan;size:100;comment:广告语" json:"slogan"`
	AIName               string  `gorm:"column:ai_name;size:50;comment:AI名字" json:"ai_name"`
	AIAvatar             string  `gorm:"column:ai_avatar;size:50;comment:AI头像" json:"ai_avatar"`
	ReferralReward       int     `gorm:"column:referral_reward;comment:邀请奖励开关1打开 2关闭" json:"referral_reward"`
	WeChatQRCode         string  `gorm:"column:wechat_qr_code;size:255;comment:微信二维码" json:"wechat_qr_code"`
	CustomerServicePhone string  `gorm:"column:customer_service_phone;size:20;comment:客服电话" json:"customer_service_phone"`
	Domain               string  `gorm:"column:domain;size:100;comment:域名多个用逗号隔开" json:"domain"`
	ReferralCode         string  `gorm:"column:referral_code;size:100;comment:邀请码" json:"referral_code"`
	Balance              float64 `gorm:"column:balance;type:decimal(20,2);default:0"`
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

func GetAgentBalance(referralCode string) (float64, error) {
	agent := new(Agent)
	err := g.DB().Model(agent.TableName()).Where("referral_code", referralCode).Scan(agent)
	if err != nil {
		return 0, err
	}
	return agent.Balance, nil
}

func UpdateAgentBalance(referralCode string, newBalance float64) error {
	agent := new(Agent)
	err := g.DB().Model(agent.TableName()).Where("referral_code", referralCode).Scan(agent)
	if err != nil {
		return err
	}

	agent.Balance = newBalance
	_, err = g.DB().Model(agent.TableName()).Where("referral_code", referralCode).Data(agent).Update()
	if err != nil {
		return err
	}

	return nil
}

func UpdateAgentBalanceByAgentID(id int, newBalance float64) error {
	agent := new(Agent)
	err := g.DB().Model(agent.TableName()).Where("id", id).Scan(agent)
	if err != nil {
		return err
	}

	agent.Balance = newBalance
	_, err = g.DB().Model(agent.TableName()).Where("id", id).Data(agent).Update()
	if err != nil {
		return err
	}

	return nil
}

func AddAgentBalanceByAgentID(id int, amount float64) error {
	agent := new(Agent)
	err := g.DB().Model(agent.TableName()).Where("id", id).Scan(agent)
	if err != nil {
		return err
	}

	newBalance := amount + agent.Balance
	return UpdateAgentBalanceByAgentID(id, newBalance)
}
