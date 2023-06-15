package middleware

import (
	"gpt_admin_go/modules/user/config"

	"github.com/gogf/gf/v2/frame/g"
)

func init() {
	if config.Config.Middleware.Authority.Enable {
		g.Server().BindMiddleware("/admin/base/open/login/*", BaseAuthorityMiddlewareOpen)
		g.Server().BindMiddleware("/admin/user/user_list", TokenAuthMiddleware)
		g.Server().BindMiddleware("/admin/user/add", TokenAuthMiddleware)
		g.Server().BindMiddleware("/admin/user/delete", TokenAuthMiddleware)
		g.Server().BindMiddleware("/admin/user/update", AdminTokenAuthMiddleware)
		g.Server().BindMiddleware("/admin/user/info", TokenAuthMiddleware)
		g.Server().BindMiddleware("/admin/user/list", TokenAuthMiddleware)
		g.Server().BindMiddleware("/admin/user/page", TokenAuthMiddleware)
		g.Server().BindMiddleware("/admin/user/*", TokenAuthMiddleware)
		g.Server().BindMiddleware("/admin/recharge_card/*", TokenAuthMiddleware)
	}
	if config.Config.Middleware.Log.Enable {
		g.Server().BindMiddleware("/admin/*", BaseLog)
	}

}
