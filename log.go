package malatd

import (
	"github.com/swxctx/xlog"
)

// runLog
func runLog(ctx *Context) {
	xlog.Infof("From %s, %s", ctx.CallCtx.RemoteAddr().String(), ctx.CallCtx.Request.String())
	ctx.Next()
}
