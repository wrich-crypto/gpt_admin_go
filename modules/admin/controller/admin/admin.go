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
	if user.HasRole(model.ROLE_ADMIN) {
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
	Logo                 string  `p:"logo"`
	ReferralReward       int     `p:"referral_reward"`
	UserId               int     `p:"user_id"`
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
	if !user.HasRole(model.ROLE_AGENT) && !user.HasRole(model.ROLE_ADMIN) {
		g.Log().Error(ctx, "UpdateAgent user role not authorized")
		return cool.Fail("User role not authorized"), nil
	}

	agent := admin_model.NewAgent()

	// For super admin
	setting_user_id := 0
	if user.HasRole(model.ROLE_ADMIN) {
		// Fetch the agent based on the provided id
		err := g.DB().Model(agent.TableName()).Where("id", req.ID).Scan(agent)
		if err != nil {
			g.Log().Error(ctx, "UpdateAgent error fetching agent", err)
			return cool.Fail("Error fetching agent"), err
		}
		agent.Balance = req.Balance
		setting_user_id = agent.UserId
	} else if user.HasRole(model.ROLE_AGENT) {
		// Fetch the agent based on the user's referral code
		err := g.DB().Model(agent.TableName()).Where("user_id", user.ID).Scan(agent)
		if err != nil {
			g.Log().Error(ctx, "UpdateAgent error fetching agent", err)
			return cool.Fail("Error fetching agent"), err
		}
		setting_user_id = int(user.ID)
	} else {
		g.Log().Error(ctx, "UpdateAgent user role not authorized")
		return cool.Fail("User role not authorized"), nil
	}

	// Update the fields
	agent.UserId = setting_user_id
	agent.AIName = req.AIName
	agent.AIAvatar = req.AIAvatar
	agent.Logo = req.Logo
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

func (c *AgentController) GetAgentInfo(ctx context.Context, req *GetAgentReq) (res *cool.BaseRes, err error) {
	user, ok := ctx.Value("user").(*model.Users)
	if !ok {
		return cool.Fail("Invalid user context"), fmt.Errorf("Invalid user context")
	}

	agent := admin_model.NewAgent()

	// Fetch the agent based on the user's referral code
	err = g.DB().Model(agent.TableName()).Where("user_id", user.ID).Scan(agent)
	if err != nil {
		g.Log().Error(ctx, "GetAgentInfo error fetching agent", err)
		return cool.Fail("Error fetching agent"), err
	}

	// Fetch the consumed points
	var consumedPoints float64
	rechargeCards := model.NewRechargeCards()
	sum, err := g.DB().Model(rechargeCards.TableName()).Where("create_user_id", user.ID).Sum("recharge_amount")
	if err != nil {
		g.Log().Error(ctx, "GetAgentInfo error fetching consumed points", err)
		return cool.Fail("Error fetching consumed points"), err
	}
	consumedPoints = sum

	responseData := g.Map{
		"agent":           agent,
		"consumed_points": consumedPoints, // include consumed points in the response
	}
	return cool.Ok(responseData), nil
}
