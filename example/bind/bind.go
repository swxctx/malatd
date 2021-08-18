package main

import (
	"github.com/swxctx/malatd/binding"

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
	A int    `json:"a"`
	B string `json:"b"`
}

var (
	binder = binding.JSON
)

// malatd
func malatdApi(ctx *malatd.Context) {
	// bind params
	params := new(Args)
	err := binder.Bind(ctx, params)
	if err != nil {
		panic(err)
	}
	xlog.Infof("Args-> %v", params)
	ctx.String(200, "malatd")
}
