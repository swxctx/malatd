package td

import (
	"net/http"
	"path"
	"strings"

	"github.com/swxctx/gutil"
	"github.com/swxctx/xlog"

	"github.com/swxctx/malatd/httprouter"
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

// AddPlugin
func (r *Router) AddPlugin(plugins ...Plugin) {
	r.Plugins = append(r.Plugins, plugins...)
}

// Get
func (r *Router) Get(relativePath string, plugins ...Plugin) {
	path := getReqPath(r.basePath, relativePath)
	plugin := append(r.Plugins, plugins...)
	xlog.Infof("[ROUTE]: GET %s", path)
	r.handle("GET", path, plugin)
}

// Post
func (r *Router) Post(relativePath string, plugins ...Plugin) {
	path := getReqPath(r.basePath, relativePath)
	plugin := append(r.Plugins, plugins...)
	xlog.Infof("[ROUTE]: POST %s", path)
	r.handle("POST", path, plugin)
}

// Options
func (r *Router) Options(relativePath string, plugins ...Plugin) {
	path := getReqPath(r.basePath, relativePath)
	plugin := append(r.Plugins, plugins...)
	xlog.Infof("[ROUTE]: OPTIONS %s", path)
	r.handle("OPTIONS", path, plugin)
}

// handle
func (r *Router) handle(httpMethod, relativePath string, plugins Plugins) {
	switch httpMethod {
	case "GET":
		r.server.router.GET(relativePath, func(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
			ctx := Context{
				index:   0,
				server:  r.server,
				plugins: plugins,
			}
			defer func() {
				if re := recover(); re != nil {
					xlog.Errorf("[GET] err: %v", re)
					ctx.Render(RerrInternalServer)
				}
			}()
			ctx.Request = request
			ctx.ResponseWriter = response
			ctx.Next()
		})
	case "POST":
		r.server.router.POST(relativePath, func(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
			ctx := Context{
				index:   0,
				server:  r.server,
				plugins: plugins,
			}
			defer func() {
				if re := recover(); re != nil {
					xlog.Errorf("[POST] err: %v", re)
					ctx.Render(RerrInternalServer)
				}
			}()
			ctx.Request = request
			ctx.ResponseWriter = response
			ctx.Next()
		})
	case "OPTIONS":
		r.server.router.OPTIONS(relativePath, func(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
			ctx := Context{

				index:   0,
				server:  r.server,
				plugins: plugins,
			}
			defer func() {
				if re := recover(); re != nil {
					xlog.Errorf("[POST] err: %v", re)
					ctx.Render(RerrInternalServer)
				}
			}()
			ctx.Request = request
			ctx.ResponseWriter = response
			ctx.Next()
		})
	}
}

// ToUriPath maps struct(func) name to URI path.
func ToUriPath(name string) string {
	p := strings.Replace(name, "__", ".", -1)
	a := strings.Split(p, "_")
	for k, v := range a {
		a[k] = gutil.FieldSnakeString(v)
	}
	p = path.Join(a...)
	p = path.Join("/", p)
	return strings.Replace(p, ".", "_", -1)
}
