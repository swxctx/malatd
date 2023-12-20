package logic

import (
	td "github.com/swxctx/malatd"
	"github.com/swxctx/malatd/example/apiserver/args"
)

// IndexLogic
func IndexLogic(ctx *td.Context, arg *args.IndexArgs) (*args.IndexResult, *td.Rerror) {
	targetUrl := "https://www.baidu.com"
	ctx.Redirect(targetUrl, 302)
	return nil, nil
}
