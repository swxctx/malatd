package api

import (
	td "github.com/swxctx/malatd"
	"github.com/swxctx/malatd/binding"
	"github.com/swxctx/malatd/example/apiserver/args"
	"github.com/swxctx/malatd/example/apiserver/logic"
)

// loginHandle
func loginHandle(ctx *td.Context) {
	// bind arg
	arg := new(args.LoginArgs)
	if err := binding.Binder(ctx, arg); err != nil {
		ctx.RenderRerr(td.RerrInternalServer.SetReason(err.Error()))
		return
	}

	// api逻辑调用
	result, rerr := logic.LoginLogic(ctx, arg)
	if rerr != nil {
		ctx.RenderRerr(rerr)
		return
	}
	ctx.Render(result)
}

// indexHandle
func indexHandle(ctx *td.Context) {
	logic.IndexLogic(ctx, nil)
}
