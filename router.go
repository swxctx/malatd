package malatd

import (
	"github.com/swxctx/xlog"
	"github.com/valyala/fasthttp"
)

// Router
type Router struct {
	Plugins  Plugins
	basePath string
	server   *Server
	root     bool
}

// 注册组路由
func (r *Router) Group(relativePath string, plugins ...Plugin) *Router {
	gPlugins := append(r.Plugins, plugins...)
	return &Router{
		Plugins:  gPlugins,
		basePath: getReqPath(r.basePath, relativePath),
		server:   r.server,
		root:     false,
	}
}

// Use
func (r *Router) Use(plugins ...Plugin) {
	r.Plugins = append(r.Plugins, plugins...)
}

// Get
func (r *Router) Get(relativePath string, plugins ...Plugin) {
	path := getReqPath(r.basePath, relativePath)
	plugin := append(r.Plugins, plugins...)
	r.handle("GET", path, plugin)
}

// Post
func (r *Router) Post(relativePath string, plugins ...Plugin) {
	path := getReqPath(r.basePath, relativePath)
	plugin := append(r.Plugins, plugins...)
	r.handle("POST", path, plugin)
}

// Options
func (r *Router) Options(relativePath string, plugins ...Plugin) {
	path := getReqPath(r.basePath, relativePath)
	plugin := append(r.Plugins, plugins...)
	r.handle("OPTIONS", path, plugin)
}

// handle
func (r *Router) handle(httpMethod, relativePath string, plugins Plugins) {
	ctx := Context{
		index:   0,
		server:  r.server,
		plugins: plugins,
	}

	var (
		err error
	)

	switch httpMethod {
	case "GET":
		r.server.router.GET(relativePath, func(ctxF *fasthttp.RequestCtx) {
			defer func() {
				if re := recover(); re != nil {
					ctx.CallCtx.SetStatusCode(500)
					xlog.Errorf("[GET] err-> %v", re)
					_, err = ctx.CallCtx.WriteString("server error")
				}
			}()
			ctx.CallCtx = ctxF
			ctx.Next()
		})
	case "POST":
		r.server.router.POST(relativePath, func(ctxF *fasthttp.RequestCtx) {
			defer func() {
				if re := recover(); re != nil {
					ctx.CallCtx.SetStatusCode(500)
					xlog.Errorf("[POST] err-> %v", re)
					_, err = ctx.CallCtx.WriteString("server error")
				}
			}()
			ctx.CallCtx = ctxF
			ctx.Next()
		})
	case "OPTIONS":
		r.server.router.OPTIONS(relativePath, func(ctxF *fasthttp.RequestCtx) {
			defer func() {
				if re := recover(); re != nil {
					ctx.CallCtx.SetStatusCode(500)
					xlog.Errorf("[OPTIONS] err-> %v", re)
					_, err = ctx.CallCtx.WriteString("server error")
				}
			}()
			ctx.CallCtx = ctxF
			ctx.Next()
		})
	}
	if err != nil {
		xlog.Errorf("[PLUGIN] err-> %v", err)
	}
}

// getReqPath
func getReqPath(h1, h2 string) string {
	u := string(h1[len(h1)-1])
	if u == "/" {
		u = h1[:len(h1)-1]
	} else {
		u = h1
	}

	u2 := string(h2[0])
	if u2 == "/" {
		u += h2
	} else {
		u += "/" + h2
	}

	return u
}
