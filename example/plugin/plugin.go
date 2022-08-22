package main

import (
	td "github.com/swxctx/malatd"
	"github.com/swxctx/xlog"
)

/*
	http://127.0.0.1:8080/malatd?id=123
*/
func main() {
	// new server
	srv := td.NewServer(
		td.NewSrvConfig(),
		tokenPlugin,
	)

	// api router
	srv.Get("/malatd/v1", authPlugin, malatdApi)
	srv.Run()
}

// tokenPlugin
func tokenPlugin(ctx *td.Context) {
	xlog.Infof("tokenPlugin: current tokenPlugin")
	xlog.Infof("tokenPlugin: Params-> %v", ctx.Request.URL.Query().Get("id"))
	ctx.Next()
}

// authHandle
func authPlugin(ctx *td.Context) {
	xlog.Infof("authPlugin: current authPlugin")
	xlog.Infof("authPlugin: Params-> %v", ctx.Request.URL.Query().Get("id"))
	ctx.Next()
}

// malatd
func malatdApi(ctx *td.Context) {
	xlog.Infof("malatdApi: current malatdApi")
	xlog.Infof("malatdApi: Params-> %v", ctx.Request.URL.Query().Get("id"))
	ctx.RenderString("hello malatd")
}
