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
	srv := malatd.NewServer(malatd.NewSrvConfig())

	// add global handle
	srv.Use(tokenHandle)

	// api router
	srv.Get("/malatd", authHandle, malatdApi)
	srv.Run()
}

// tokenHandle
func tokenHandle(ctx *malatd.Context) {
	xlog.Infof("tokenHandle: current tokenHandle")
	xlog.Infof("tokenHandle: Params-> %v", string(ctx.Ctx.QueryArgs().Peek("id")))
	ctx.Next()
}

// authHandle
func authHandle(ctx *malatd.Context) {
	xlog.Infof("authHandle: current authHandle")
	xlog.Infof("authHandle: Params-> %v", string(ctx.Ctx.QueryArgs().Peek("id")))
	ctx.Next()
}

// malatd
func malatdApi(ctx *malatd.Context) {
	xlog.Infof("malatdApi: current malatdApi")
	xlog.Infof("malatdApi: Params-> %v", string(ctx.Ctx.QueryArgs().Peek("id")))
	ctx.String(200, "hello malatd")
}
