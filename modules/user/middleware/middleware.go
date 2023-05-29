package middleware

import (
	"gpt_admin_go/modules/user/config"

	"github.com/gogf/gf/v2/frame/g"
)

func init() {
	if config.Config.Middleware.Authority.Enable {
		g.Server().BindMiddleware("/admin/base/open/login/*", BaseAuthorityMiddlewareOpen)
		g.Server().BindMiddleware("/admin/user/*", TokenAuthMiddleware)
		g.Server().BindMiddleware("/admin/recharge_card/*", TokenAuthMiddleware)
	}
	if config.Config.Middleware.Log.Enable {
		g.Server().BindMiddleware("/admin/*", BaseLog)
	}

}
