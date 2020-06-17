package malatd

import (
	"github.com/swxctx/xlog"
)

// runLog
func runLog(ctx *Context) {
	xlog.Infof("From %s, %s", ctx.Ctx.RemoteAddr().String(), string(ctx.Ctx.Method()), ctx.Ctx.Request.String())
	ctx.Next()
}
