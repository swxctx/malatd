package main

import (
	"github.com/swxctx/malatd/binding"

	td "github.com/swxctx/malatd"
)

/*
	http://127.0.0.1:8080/malatd
*/
func main() {
	// new server
	srv := td.NewServer(td.NewSrvConfig())

	// api router
	srv.Get("/malatd1", malatdApi1)
	srv.Post("/malatd2", malatdApi2)
	srv.Run()
}

type Args struct {
	A int    `form:"a" json:"a"`
	B string `form:"b" json:"b"`
	C string `json:"c"`
}

type Result struct {
	A int    `json:"a"`
	B string `json:"b"`
}

var (
	binder      = binding.JSON
	binderQuery = binding.QUERY
)

// malatd
func malatdApi1(ctx *td.Context) {
	// bind params
	params := new(Args)
	err := binderQuery.Bind(ctx, params)
	if err != nil {
		panic(err)
	}
	td.Infof("Args-> %v", params)

	result := &Result{
		A: params.A,
		B: params.B,
	}
	ctx.Json(result)
}

// malatd
func malatdApi2(ctx *td.Context) {
	// bind params
	params := new(Args)
	err := binder.Bind(ctx, params)
	if err != nil {
		panic(err)
	}

	err = binderQuery.Bind(ctx, params)
	if err != nil {
		panic(err)
	}
	td.Infof("Args-> %v", params)

	ctx.String(200, "malatd")
}
