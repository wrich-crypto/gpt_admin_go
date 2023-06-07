package middleware

import (
	"gpt_admin_go/modules/user/config"

	"github.com/gogf/gf/v2/frame/g"
)

func init() {
	if config.Config.Middleware.Authority.Enable {
		g.Server().BindMiddleware("/admin/api_key/*", AdminTokenAuthMiddleware)
		g.Server().BindMiddleware("/admin/dev_config/*", AdminTokenAuthMiddleware)
		g.Server().BindMiddleware("/admin/agent/add", AdminTokenAuthMiddleware)
		g.Server().BindMiddleware("/admin/agent/delete", AdminTokenAuthMiddleware)
		g.Server().BindMiddleware("/admin/agent/update", AdminTokenAuthMiddleware)
		g.Server().BindMiddleware("/admin/agent/info", AdminTokenAuthMiddleware)
		g.Server().BindMiddleware("/admin/agent/list", AdminTokenAuthMiddleware)
		g.Server().BindMiddleware("/admin/agent/page", AdminTokenAuthMiddleware)
		g.Server().BindMiddleware("/admin/agent/add_balance", AdminTokenAuthMiddleware)
		g.Server().BindMiddleware("/admin/agent/update_agent", TokenAuthMiddleware)
		g.Server().BindMiddleware("/admin/agent/get_agent", TokenAuthMiddleware)
	}
	if config.Config.Middleware.Log.Enable {
		g.Server().BindMiddleware("/admin/*", BaseLog)
	}

}
