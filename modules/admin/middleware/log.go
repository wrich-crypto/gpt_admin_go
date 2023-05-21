package middleware

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func BaseLog(r *ghttp.Request) {
	var (
		ctx = r.GetCtx()
	)
	Record(ctx)

	r.Middleware.Next()
}

// Record 记录日志
func Record(ctx g.Ctx) {
	var (
		r = g.RequestFromCtx(ctx)
	)
	action := r.Method + ":" + r.URL.Path
	IP := r.GetClientIp()
	IPAddr := r.GetClientIp()
	Params := r.GetBodyString()

	g.Log().Infof(ctx, "Action: %s, IP: %s, IPAddr: %s, Params: %s", action, IP, IPAddr, Params)
}
