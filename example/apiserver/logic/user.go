package logic

import (
	td "github.com/swxctx/malatd"
	"github.com/swxctx/malatd/example/apiserver/args"
	"github.com/swxctx/malatd/example/apiserver/rerrs"
)

// LoginLogic
func LoginLogic(ctx *td.Context, arg *args.LoginArgs) (*args.LoginResult, *td.Rerror) {
	result := &args.LoginResult{
		AppVer:   arg.AppVer,
		Username: arg.Username,
	}
	if len(arg.Username) <= 0 {
		return nil, rerrs.RerrUserNotExists
	}
	return result, nil
}
