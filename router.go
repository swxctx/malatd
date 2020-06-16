package malatd

import (
	"github.com/swxctx/xlog"
	"github.com/valyala/fasthttp"
)

// RouterGroup
type RouterGroup struct {
	Handlers HandlersChain
	basePath string
	server   *Server
	root     bool
}

// Use
func (r *RouterGroup) Use(middleware ...HandlerFunc) {
	r.Handlers = append(r.Handlers, middleware...)
}

// Get
func (r *RouterGroup) Get(relativePath string, handlers ...HandlerFunc) *urlp {
	url := GetUri(r.basePath, relativePath)
	handle := append(r.Handlers, handlers...)
	r.handle("GET", url, handle)
	return &urlp{url: url}
}

// handle
func (r *RouterGroup) handle(httpMethod, relativePath string, handlers HandlersChain) {
	ctx := Context{
		i:    0,
		server:   r.server,
		handlers: handlers,
	}

	var (
		err error
	)

	r.server.router.GET(relativePath, func(ctxF *fasthttp.RequestCtx) {
		defer func() {
			if re := recover(); re != nil {
				ctx.Ctx.SetStatusCode(500)
				xlog.Errorf("[GET] handle err-> %v",re)
				_, err = ctx.Ctx.WriteString("server error")
			}
		}()
		ctx.Ctx = ctxF
		ctx.Next()
	})
	if err != nil {
		xlog.Errorf("[HANDLE] err-> %v",err)
	}
}

// GetUri
func GetUri(h1, h2 string) string {
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