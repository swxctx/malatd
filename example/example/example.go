package main

import (
	td "github.com/swxctx/malatd"
	"github.com/swxctx/malatd/binding"
)

/**
    @date: 2022/3/14
**/

/*
	http://127.0.0.1:8080/malatd
*/

func main() {
	// new server
	srv := td.NewServer(td.NewSrvConfig())

	login := srv.Group("/user")
	login.Post("/login", loginHandle)
	login.Get("/login1", loginHandle)
	srv.Post("/logi2", loginHandle)
	srv.Run()
}

type (
	LoginArgs struct {
		AppVer   string `query:"app_ver" json:"app_ver"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	LoginResult struct {
		AppVer   string `json:"app_ver"`
		Username string `json:"username"`
	}
)

// loginHandle
func loginHandle(ctx *td.Context) {
	// bind arg
	arg := new(LoginArgs)
	err := binding.Binder(ctx, arg)
	if err != nil {
		ctx.RenderRerr(td.RerrInternalServer.SetReason(err.Error()))
		return
	}

	// api逻辑调用
	result, rerr := loginLogic(ctx, arg)
	if rerr != nil {
		ctx.RenderRerr(rerr)
		return
	}
	ctx.Render(result)
}

// malatdApi1Logic
func loginLogic(ctx *td.Context, arg *LoginArgs) (*LoginResult, *td.Rerror) {
	result := &LoginResult{
		AppVer:   arg.AppVer,
		Username: arg.Username,
	}
	return result, nil
}
