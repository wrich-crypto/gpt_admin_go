package middleware

import (
	"gpt_admin_go/modules/user/config"

	"github.com/gogf/gf/v2/frame/g"
)

func init() {
	if config.Config.Middleware.Authority.Enable {
		g.Server().BindMiddleware("/admin/api_key/*", AdminTokenAuthMiddleware)
		g.Server().BindMiddleware("/admin/dev_config/*", AdminTokenAuthMiddleware)
		g.Server().BindMiddleware("/admin/agent/*", TokenAuthMiddleware)
	}
	if config.Config.Middleware.Log.Enable {
		g.Server().BindMiddleware("/admin/*", BaseLog)
	}

}
