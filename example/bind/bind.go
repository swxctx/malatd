package main

import (
	"encoding/json"
	"github.com/swxctx/malatd"
	"github.com/swxctx/malatd/binding"
	"github.com/swxctx/xlog"
)

var (
	myApiware = binding.New(nil, nil)
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
	A int `param:"in(query)" json:"a"`
	B int `param:"in(query)" json:"b"`
	// query
	C int `param:"in(query)" json:"c"`
}

// malatd
func malatdApi(ctx *malatd.Context) {
	err := myApiware.Register(
		new(Args),
	)
	if err != nil {
		panic(err)
	}

	// bind params
	params := new(Args)
	err = myApiware.Bind(params, ctx.CallCtx)
	if err != nil{
		panic(err)
	}
	json.MarshalIndent(params, "", " ")
	xlog.Infof("Args-> %v",params)

	ctx.String(200, "malatd")
}
