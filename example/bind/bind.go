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
	srv.Get("/malatd1", malatdApi1Handle)
	srv.Post("/malatd2", malatdApi2Handle)
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

// malatdApi1Handle
func malatdApi1Handle(ctx *td.Context) {
	// bind params
	params := new(Args)
	err := binderQuery.Bind(ctx, params)
	if err != nil {
		panic(err)
	}

	// api逻辑调用
	result, rerr := malatdApi1Logic(ctx, params)
	if rerr != nil {
		ctx.RenderRerr(rerr)
		return
	}
	ctx.Render(result)
}

// malatdApi1Logic
func malatdApi1Logic(ctx *td.Context, arg *Args)(*Result, *td.Rerror) {
	td.Infof("Args-> %v", arg)
	result := &Result{
		A: arg.A,
		B: arg.B,
	}
	return result, nil
}

// malatdApi2Handle
func malatdApi2Handle(ctx *td.Context) {
	// bind params
	params := new(Args)
	err := binder.Bind(ctx, params)
	if err != nil {
		//ctx.RenderRerr(td.RerrInternalServer.SetReason(err.Error()))
		//return
		panic(td.RerrInternalServer.SetReason(err.Error()))
	}

	err = binderQuery.Bind(ctx, params)
	if err != nil {
		panic(err)
	}

	// api逻辑调用
	result, rerr := malatdApi2Logic(ctx, params)
	if rerr != nil {
		ctx.RenderRerr(rerr)
		return
	}
	ctx.Render(result)
}

func malatdApi2Logic(ctx *td.Context, arg *Args)(*Result, *td.Rerror) {
	td.Infof("Args-> %v", arg)
	result := &Result{
		A: arg.A,
		B: arg.B,
	}
	return result, nil
}
