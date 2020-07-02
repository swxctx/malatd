package main

import (
	"github.com/swxctx/malatd"
	"github.com/swxctx/xlog"
)

/*
	http://127.0.0.1:8080/malatd
*/
func main() {
	// new server
	srv := malatd.NewServer(malatd.NewSrvConfig())

	// api router
	srv.Get("/malatd", malatdApi)
	srv.Post("/malatd", malatdApi)
	srv.Run()
}

type Args struct {
	// body
	A int `json:"a"`
	B int `json:"b"`
	// query
	C int `json:"c"`
}

// malatd
func malatdApi(ctx *malatd.Context) {
	var (
		argData *Args
	)

	xlog.Infof("Content-> %s", ctx.ContentType())
	ctx.Bind(&argData)

	xlog.Infof("args-> %v", argData)
	ctx.String(200, "malatd")
}
