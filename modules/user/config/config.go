package config

import (
	"github.com/cool-team-official/cool-admin-go/cool"
	"github.com/gogf/gf/v2/frame/g"
)

// sConfig 配置
type sConfig struct {
	Jwt        *Jwt
	Middleware *Middleware
	Oss        *Oss
}

type Oss struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
}

type Middleware struct {
	Authority *Authority
	Log       *Log
}

type Authority struct {
	Enable bool
}

type Log struct {
	Enable bool
}

type Token struct {
	Expire        uint `json:"expire"`
	RefreshExpire uint `json:"refreshExprire"`
}

type Jwt struct {
	Sso    bool   `json:"sso"`
	Secret string `json:"secret"`
	Token  *Token `json:"token"`
}

// NewConfig new config
func NewConfig() *sConfig {
	var (
		ctx g.Ctx
	)
	config := &sConfig{
		Jwt: &Jwt{
			Sso:    cool.GetCfgWithDefault(ctx, "modules.user.jwt.sso", g.NewVar(false)).Bool(),
			Secret: cool.GetCfgWithDefault(ctx, "modules.user.jwt.secret", g.NewVar(cool.ProcessFlag)).String(),
			Token: &Token{
				Expire:        cool.GetCfgWithDefault(ctx, "modules.user.jwt.token.expire", g.NewVar(2*3600)).Uint(),
				RefreshExpire: cool.GetCfgWithDefault(ctx, "modules.user.jwt.token.refreshExpire", g.NewVar(15*24*3600)).Uint(),
			},
		},
		Middleware: &Middleware{
			Authority: &Authority{
				Enable: cool.GetCfgWithDefault(ctx, "modules.user.middleware.authority.enable", g.NewVar(true)).Bool(),
			},
			Log: &Log{
				Enable: cool.GetCfgWithDefault(ctx, "modules.user.middleware.log.enable", g.NewVar(true)).Bool(),
			},
		},
		Oss: &Oss{
			Endpoint:        cool.GetCfgWithDefault(ctx, "modules.user.oss.endpoint", g.NewVar("")).String(),
			AccessKeyID:     cool.GetCfgWithDefault(ctx, "modules.user.oss.accessKeyID", g.NewVar("")).String(),
			AccessKeySecret: cool.GetCfgWithDefault(ctx, "modules.user.oss.accessKeySecret", g.NewVar("")).String(),
			BucketName:      cool.GetCfgWithDefault(ctx, "modules.user.oss.bucketName", g.NewVar("")).String(),
		},
	}

	return config
}

// Config config
var Config = NewConfig()
