package malatd

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/swxctx/xlog"
	"github.com/valyala/fasthttp"
)

// Server
type Server struct {
	// 路由组
	RouterGroup
	// 服务配置
	srvConfig *SrvConfig
	// api路由
	router *fasthttprouter.Router
}

// NewServer
func NewServer(srvCfg *SrvConfig, plugins ...Plugin) *Server {
	// server config is nil
	if srvCfg == nil {
		panic("Malatd: srv config is nil.")
	}

	// :8080 -> 0.0.0.0:8080
	if string(srvCfg.Address[0]) == ":" {
		srvCfg.Address = "0.0.0.0" + srvCfg.Address
	}

	// add runlog plugin
	if srvCfg.RunLog {
		plugins = append(plugins, runLog)
	}

	// new server
	srv := &Server{
		RouterGroup: RouterGroup{
			Plugins:  plugins,
			root:     true,
			basePath: "/",
		},
		router:    fasthttprouter.New(),
		srvConfig: srvCfg,
	}
	srv.RouterGroup.server = srv
	return srv
}

// ListenAndServe fast http listen
func (srv *Server) Run() error {
	xlog.Infof("Server Run %s", srv.srvConfig.Address)

	// start listen
	if err := fasthttp.ListenAndServe(srv.srvConfig.Address, srv.router.Handler); err != nil {
		return err
	}

	return nil
}
