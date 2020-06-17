package main

import (
	"github.com/swxctx/malatd"
	"github.com/swxctx/xlog"
)

/*
	http://127.0.0.1:8080/malatd?id=123
*/
func main() {
	// new server
	srv := malatd.NewServer(
		malatd.NewSrvConfig(),
		tokenPlugin,
	)

	// api router
	srv.Get("/malatd", authPlugin, malatdApi)
	srv.Run()
}

// tokenPlugin
func tokenPlugin(ctx *malatd.Context) {
	xlog.Infof("tokenPlugin: current tokenPlugin")
	xlog.Infof("tokenPlugin: Params-> %v", string(ctx.Ctx.QueryArgs().Peek("id")))
	ctx.Next()
}

// authHandle
func authPlugin(ctx *malatd.Context) {
	xlog.Infof("authPlugin: current authPlugin")
	xlog.Infof("authPlugin: Params-> %v", string(ctx.Ctx.QueryArgs().Peek("id")))
	ctx.Next()
}

// malatd
func malatdApi(ctx *malatd.Context) {
	xlog.Infof("malatdApi: current malatdApi")
	xlog.Infof("malatdApi: Params-> %v", string(ctx.Ctx.QueryArgs().Peek("id")))
	ctx.String(200, "hello malatd")
}
