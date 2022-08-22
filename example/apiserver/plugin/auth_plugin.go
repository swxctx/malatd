package plugin

import (
	td "github.com/swxctx/malatd"
	"github.com/swxctx/xlog"
)

// AuthPlugin
func AuthPlugin(ctx *td.Context) {
	xlog.Infof("AuthPlugin: current authPlugin")
	xlog.Infof("AuthPlugin: AppVer-> %s", ctx.Request.URL.Query().Get("app_ver"))
	ctx.Next()
}
