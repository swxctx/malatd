# malatd
简洁+易用的Golang Web框架

### 使用
[查看详细使用说明](https://github.com/swxctx/malatd/tree/master/cmd)

```shell
NAME:
   Malatd project command - A deployment tools of malatd frameware

USAGE:
   malatd [global options] command [command options] [arguments...]

VERSION:
   1.0.0

AUTHOR:
   swxctx

COMMANDS:
   gen      Generate a malatd project
   run      Compile and run go project
   doc      Generate a project README.md(malatd doc || malatd doc -r ${root_group})
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

### 普通使用
```
package main

import (
	td "github.com/swxctx/malatd"
)

/*
	http://127.0.0.1:8080/malatd
*/
func main() {
	// new server
	srv := td.NewServer(td.NewSrvConfig())

	// api router
	srv.Get("/malatd", malatdApi)
	srv.Post("/malatd", malatdApi)
	srv.Run()
}

// malatd
func malatdApi(ctx *td.Context) {
	ctx.RenderString("malatd")
}
```

### 使用插件
```
package main

import (
	td "github.com/swxctx/malatd"
)

/*
	http://127.0.0.1:8080/malatd?id=123
*/
func main() {
	// new server
	srv := td.NewServer(
		td.NewSrvConfig(),
		tokenPlugin,
	)

	// api router
	srv.Get("/malatd/v1", authPlugin, malatdApi)
	srv.Run()
}

// tokenPlugin
func tokenPlugin(ctx *td.Context) {
	xlog.Infof("tokenPlugin: current tokenPlugin")
	xlog.Infof("tokenPlugin: Params-> %v", ctx.Request.URL.Query().Get("id"))
	ctx.Next()
}

// authHandle
func authPlugin(ctx *td.Context) {
	xlog.Infof("authPlugin: current authPlugin")
	xlog.Infof("authPlugin: Params-> %v", ctx.Request.URL.Query().Get("id"))
	ctx.Next()
}

// malatd
func malatdApi(ctx *td.Context) {
	xlog.Infof("malatdApi: current malatdApi")
	xlog.Infof("malatdApi: Params-> %v", ctx.Request.URL.Query().Get("id"))
	ctx.RenderString("hello malatd")
}
```

### 参数绑定
```
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
	A int    `query:"a" json:"a"`
	B string `query:"b" json:"b"`
	C string `json:"c"`
}

type Result struct {
	A int    `json:"a"`
	B string `json:"b"`
	C string `json:"c"`
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
		ctx.RenderRerr(td.RerrInternalServer.SetReason(err.Error()))
		return
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
	xlog.Infof("Args-> %v", arg)
	result := &Result{
		A: arg.A,
		B: arg.B,
		C: arg.C,
	}
	return result, nil
}

// malatdApi2Handle
func malatdApi2Handle(ctx *td.Context) {
	// bind params
	params := new(Args)
	err := binder.Bind(ctx, params)
	if err != nil {
		ctx.RenderRerr(td.RerrInternalServer.SetReason(err.Error()))
		return
	}

	err = binderQuery.Bind(ctx, params)
	if err != nil {
		ctx.RenderRerr(td.RerrInternalServer.SetReason(err.Error()))
		return
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
	xlog.Infof("Args-> %v", arg)
	result := &Result{
		A: arg.A,
		B: arg.B,
		C: arg.C,
	}
	return result, nil
}
```