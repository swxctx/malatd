package td

import (
	"net/http"

	"github.com/swxctx/xlog"

	"github.com/swxctx/malatd/httprouter"
)

// Server
type Server struct {
	// 路由组
	Router
	// 服务配置
	srvConfig *SrvConfig
	// api路由
	router *httprouter.Router
}

// NewServer
func NewServer(srvCfg *SrvConfig, plugins ...Plugin) *Server {
	// server config is nil
	if srvCfg == nil {
		srvCfg = NewSrvConfig()
	}

	// :8080 -> 0.0.0.0:8080
	if len(srvCfg.Address) <= 0 {
		srvCfg.Address = "0.0.0.0:8080"
	}

	// 请求日志
	plugins = append(plugins, runLogPlugin)

	// new server
	srv := &Server{
		Router: Router{
			Plugins:  plugins,
			root:     true,
			basePath: "/",
		},
		router:    httprouter.New(),
		srvConfig: srvCfg,
	}
	srv.Router.server = srv
	srv.router.NotFound = http.HandlerFunc(renderNotFound)

	xlog.Infof("[SERVER] New server, Address: %s", srvCfg.Address)
	return srv
}

// ListenAndServe fast http listen
func (srv *Server) Run() error {
	xlog.Infof("[SERVER] %s Server Run", srv.srvConfig.Address)

	// start listen
	if err := http.ListenAndServe(srv.srvConfig.Address, srv.router); err != nil {
		xlog.Errorf("[SERVER] Server Listen err: %v", err)
		return err
	}

	return nil
}
