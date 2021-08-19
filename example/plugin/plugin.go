package main

import (
	td "github.com/swxctx/malatd"
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
	td.Infof("tokenPlugin: current tokenPlugin")
	td.Infof("tokenPlugin: Params-> %v", ctx.Request.URL.Query().Get("id"))
	ctx.Next()
}

// authHandle
func authPlugin(ctx *td.Context) {
	td.Infof("authPlugin: current authPlugin")
	td.Infof("authPlugin: Params-> %v", ctx.Request.URL.Query().Get("id"))
	ctx.Next()
}

// malatd
func malatdApi(ctx *td.Context) {
	td.Infof("malatdApi: current malatdApi")
	td.Infof("malatdApi: Params-> %v", ctx.Request.URL.Query().Get("id"))
	ctx.RspString(200, "hello malatd")
}
