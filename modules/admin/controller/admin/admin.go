package admin

import (
	"context"
	"fmt"
	admin_model "gpt_admin_go/modules/admin/model"
	"gpt_admin_go/modules/admin/service"
	"gpt_admin_go/modules/user/model"

	"github.com/cool-team-official/cool-admin-go/cool"
	"github.com/gogf/gf/v2/frame/g"
)

type ApiKeysController struct {
	*cool.Controller
}

type DevConfigController struct {
	*cool.Controller
}

type AgentController struct {
	*cool.Controller
}

func init() {
	var api_key_controller = &ApiKeysController{
		&cool.Controller{
			Perfix:  "/admin/api_key",
			Api:     []string{"Add", "Delete", "Update", "Info", "List", "Page"},
			Service: service.NewApiKeysService(),
		},
	}

	var dev_config_controller = &DevConfigController{
		&cool.Controller{
			Perfix:  "/admin/dev_config",
			Api:     []string{"Add", "Delete", "Update", "Info", "List", "Page"},
			Service: service.NewDevConfigService(),
		},
	}

	var agent_controller = &AgentController{
		&cool.Controller{
			Perfix:  "/admin/agent",
			Api:     []string{"Add", "Delete", "Update", "Info", "List", "Page"},
			Service: service.NewAgentService(),
		},
	}

	cool.RegisterController(api_key_controller)
	cool.RegisterController(dev_config_controller)
	cool.RegisterController(agent_controller)
}

type AddBalanceReq struct {
	g.Meta  `path:"/add_balance" method:"POST"`
	AgentID int     `p:"agent_id"`
	Amount  float64 `p:"amount"`
}

func (c *AgentController) HandleAddBalance(ctx context.Context, req *AddBalanceReq) (res *cool.BaseRes, err error) {
	user, ok := ctx.Value("user").(*model.Users)
	if !ok {
		return cool.Fail("Invalid user context"), fmt.Errorf("Invalid user context")
	}

	// Check if the user is a superadmin
	if user.Role != 3 {
		g.Log().Error(ctx, "HandleAddBalance user role not authorized")
		return cool.Fail("User role not authorized"), nil
	}

	err = admin_model.AddAgentBalanceByAgentID(req.AgentID, req.Amount)
	if err != nil {
		g.Log().Error(ctx, "HandleAddBalance error updating agent balance", err)
		return cool.Fail("Error updating agent balance"), err
	}

	return cool.Ok("Agent balance updated successfully"), nil
}

type UpdateAgentReq struct {
	g.Meta               `path:"/update_agent" method:"POST"`
	ID                   int     `p:"id"`
	AIName               string  `p:"ai_name"`
	AIAvatar             string  `p:"ai_avatar"`
	ReferralReward       int     `p:"referral_reward"`
	CustomerServicePhone string  `p:"customer_service_phone"`
	Domain               string  `p:"domain"`
	Slogan               string  `p:"slogan"`
	Title                string  `p:"title"`
	WeChatQRCode         string  `p:"wechat_qr_code"`
	ReferralCode         string  `p:"referral_code"`
	Balance              float64 `p:"balance"`
}

func (c *AgentController) UpdateAgent(ctx context.Context, req *UpdateAgentReq) (res *cool.BaseRes, err error) {
	// 在这里处理更新逻辑
	user, ok := ctx.Value("user").(*model.Users)
	if !ok {
		return cool.Fail("Invalid user context"), fmt.Errorf("Invalid user context")
	}

	// Check if the user is a superadmin
	if user.Role != model.ROLE_ADMIN && user.Role != model.ROLE_AGENT {
		g.Log().Error(ctx, "UpdateAgent user role not authorized")
		return cool.Fail("User role not authorized"), nil
	}

	agent := admin_model.NewAgent()

	// For super admin
	if user.Role == model.ROLE_ADMIN {
		// Fetch the agent based on the provided id
		err := g.DB().Model(agent.TableName()).Where("id", req.ID).Scan(agent)
		if err != nil {
			g.Log().Error(ctx, "UpdateAgent error fetching agent", err)
			return cool.Fail("Error fetching agent"), err
		}
		agent.Balance = req.Balance
	} else if user.Role == model.ROLE_AGENT {
		// Fetch the agent based on the user's referral code
		err := g.DB().Model(agent.TableName()).Where("referral_code", user.ReferralCode).Scan(agent)
		if err != nil {
			g.Log().Error(ctx, "UpdateAgent error fetching agent", err)
			return cool.Fail("Error fetching agent"), err
		}
	}

	// Update the fields
	agent.AIName = req.AIName
	agent.AIAvatar = req.AIAvatar
	agent.ReferralReward = req.ReferralReward
	agent.CustomerServicePhone = req.CustomerServicePhone
	agent.Domain = req.Domain
	agent.Slogan = req.Slogan
	agent.Title = req.Title
	agent.WeChatQRCode = req.WeChatQRCode
	agent.ReferralCode = req.ReferralCode

	_, err = g.DB().Model(agent.TableName()).Where("id", agent.ID).Data(agent).Update()
	if err != nil {
		g.Log().Error(ctx, "UpdateAgent error updating agent", err)
		return cool.Fail("Error updating agent"), err
	}

	return cool.Ok("Agent updated successfully"), nil
}

type GetAgentReq struct {
	g.Meta `path:"/get_agent" method:"GET"`
}

func (c *AgentController) GetAgentByReferralCode(ctx context.Context, req *GetAgentReq) (res *cool.BaseRes, err error) {
	user, ok := ctx.Value("user").(*model.Users)
	if !ok {
		return cool.Fail("Invalid user context"), fmt.Errorf("Invalid user context")
	}

	agent := admin_model.NewAgent()

	// Fetch the agent based on the user's referral code
	err = g.DB().Model(agent.TableName()).Where("referral_code", user.ReferralCode).Scan(agent)
	if err != nil {
		g.Log().Error(ctx, "GetAgentByReferralCode error fetching agent", err)
		return cool.Fail("Error fetching agent"), err
	}

	responseData := g.Map{"agent": agent}
	return cool.Ok(responseData), nil
}
