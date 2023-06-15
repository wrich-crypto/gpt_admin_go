package middleware

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"gpt_admin_go/modules/user/model"
)

// 本类接口无需权限验证
func BaseAuthorityMiddlewareOpen(r *ghttp.Request) {
	print("BaseAuthorityMiddlewareOpen")
	r.SetCtxVar("AuthOpen", true)
	r.Middleware.Next()
}

func TokenAuthMiddleware(r *ghttp.Request) {
	print("TokenAuthMiddleware")
	var (
		statusCode = 200
		ctx        = r.GetCtx()
	)

	authHeader := r.GetHeader("Authorization")

	// 验证头部信息是否存在
	if authHeader == "" {
		g.Log().Error(ctx, "TokenAuthMiddleware", "authorization header not provided")
		statusCode = 401
		r.Response.WriteStatusExit(statusCode, g.Map{
			"code":    1001,
			"message": "Login required",
		})
		return
	}

	// 验证头部信息格式是否为"Bearer <token>"
	authHeaderParts := strings.Fields(authHeader)
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		g.Log().Error(ctx, "TokenAuthMiddleware", "invalid authorization header format")
		statusCode = 401
		r.Response.WriteStatusExit(statusCode, g.Map{
			"code":    1001,
			"message": "Invalid authorization format. Format is 'Bearer <token>'",
		})
		return
	}

	token := authHeaderParts[1]

	// 从数据库中验证令牌是否存在
	user := new(model.Users) // 假设你有一个用户表的模型叫做 User
	err := g.DB().Model(user.TableName()).Where("token", token).Scan(user)
	if err != nil {
		g.Log().Error(ctx, "TokenAuthMiddleware", "error fetching user", err)
		statusCode = 500
		r.Response.WriteStatusExit(statusCode, g.Map{
			"code":    1002,
			"message": "Internal Server Error",
		})
		return
	}

	if user.ID == 0 { // 假设你的 User 模型有一个 ID 字段
		g.Log().Error(ctx, "TokenAuthMiddleware", "invalid token")
		statusCode = 401
		r.Response.WriteStatusExit(statusCode, g.Map{
			"code":    1001,
			"message": "Invalid token",
		})
		return
	}

	// 验证用户角色是否为管理员和代理商
	if !user.HasRole(model.ROLE_AGENT) && !user.HasRole(model.ROLE_ADMIN) {
		g.Log().Error(ctx, "TokenAuthMiddleware", "user role not authorized")
		statusCode = 403
		r.Response.WriteStatusExit(statusCode, g.Map{
			"code":    1003,
			"message": "User role not authorized",
		})
		return
	}

	// 将用户信息添加到context中
	ctx = context.WithValue(ctx, "user", user)
	r.SetCtx(ctx) // 将新的context关联到请求

	// 在此处添加任何你想要的逻辑
	r.Middleware.Next()
}

func AdminTokenAuthMiddleware(r *ghttp.Request) {
	print("TokenAuthMiddleware")
	var (
		statusCode = 200
		ctx        = r.GetCtx()
	)

	authHeader := r.GetHeader("Authorization")

	// 验证头部信息是否存在
	if authHeader == "" {
		g.Log().Error(ctx, "TokenAuthMiddleware", "authorization header not provided")
		statusCode = 401
		r.Response.WriteStatusExit(statusCode, g.Map{
			"code":    1001,
			"message": "Login required",
		})
		return
	}

	// 验证头部信息格式是否为"Bearer <token>"
	authHeaderParts := strings.Fields(authHeader)
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		g.Log().Error(ctx, "TokenAuthMiddleware", "invalid authorization header format")
		statusCode = 401
		r.Response.WriteStatusExit(statusCode, g.Map{
			"code":    1001,
			"message": "Invalid authorization format. Format is 'Bearer <token>'",
		})
		return
	}

	token := authHeaderParts[1]

	// 从数据库中验证令牌是否存在
	user := new(model.Users) // 假设你有一个用户表的模型叫做 User
	err := g.DB().Model(user.TableName()).Where("token", token).Scan(user)
	if err != nil {
		g.Log().Error(ctx, "TokenAuthMiddleware", "error fetching user", err)
		statusCode = 500
		r.Response.WriteStatusExit(statusCode, g.Map{
			"code":    1002,
			"message": "Internal Server Error",
		})
		return
	}

	if user.ID == 0 { // 假设你的 User 模型有一个 ID 字段
		g.Log().Error(ctx, "TokenAuthMiddleware", "invalid token")
		statusCode = 401
		r.Response.WriteStatusExit(statusCode, g.Map{
			"code":    1001,
			"message": "Invalid token",
		})
		return
	}

	// 验证用户角色是否为管理员和代理商
	if !user.HasRole(model.ROLE_ADMIN) {
		g.Log().Error(ctx, "TokenAuthMiddleware", "user role not authorized")
		statusCode = 403
		r.Response.WriteStatusExit(statusCode, g.Map{
			"code":    1003,
			"message": "User role not authorized",
		})
		return
	}

	// 将用户信息添加到context中
	ctx = context.WithValue(ctx, "user", user)
	r.SetCtx(ctx) // 将新的context关联到请求

	// 在此处添加任何你想要的逻辑
	r.Middleware.Next()
}
