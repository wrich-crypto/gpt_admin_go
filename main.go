package main

import (
	_ "gpt_admin_go/internal/packed"

	// _ "github.com/cool-team-official/cool-admin-go/contrib/drivers/sqlite"

	_ "github.com/cool-team-official/cool-admin-go/contrib/drivers/mysql"

	_ "gpt_admin_go/modules"

	_ "gpt_admin_go/modules/admin"
	_ "gpt_admin_go/modules/chat"
	_ "gpt_admin_go/modules/user"

	// _ "gpt_admin_go/modules/demo"

	"github.com/gogf/gf/v2/os/gctx"

	"gpt_admin_go/internal/cmd"
)

func main() {
	// gres.Dump()
	cmd.Main.Run(gctx.New())
}
