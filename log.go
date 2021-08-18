package malatd

import (
	"github.com/swxctx/xlog"
)

// runLogPlugin
func runLogPlugin(ctx *Context) {
	xlog.Infof("From %s, %s", ctx.CallCtx.RemoteAddr().String(), ctx.CallCtx.Request.String())
	ctx.Next()
}
