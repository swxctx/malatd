package td

import (
	"github.com/swxctx/xlog"
)

// runLogPlugin
func runLogPlugin(ctx *Context) {
	xlog.Infof("[%s] Request From: %s, %s, %s", ctx.Method(), ctx.GetRemoteIP(), ctx.RequestURI(), ctx.UserAgent())
	ctx.Next()
}
