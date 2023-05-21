package admin

import (
	"gpt_admin_go/modules/admin/service"

	"github.com/cool-team-official/cool-admin-go/cool"
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
