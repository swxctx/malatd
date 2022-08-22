package logic

import (
	td "github.com/swxctx/malatd"
	"github.com/swxctx/malatd/example/apiserver/args"
)

// LoginLogic
func LoginLogic(ctx *td.Context, arg *args.LoginArgs) (*args.LoginResult, *td.Rerror) {
	result := &args.LoginResult{
		AppVer:   arg.AppVer,
		Username: arg.Username,
	}
	return result, nil
}
