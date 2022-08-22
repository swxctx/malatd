package td

import (
	"github.com/swxctx/xlog"
)

// runLogPlugin
func runLogPlugin(ctx *Context) {
	xlog.Infof("From %s, %s", ctx.Request.RemoteAddr, ctx.Request.RequestURI)
	ctx.Next()
}
