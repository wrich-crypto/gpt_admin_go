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
