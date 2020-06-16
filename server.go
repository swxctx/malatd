package malatd

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/swxctx/xlog"
	"github.com/valyala/fasthttp"
)

// Server
type Server struct {
	RouterGroup
	srvConfig *SrvConfig
	router *fasthttprouter.Router
}

// NewServer
func NewServer(srvCfg *SrvConfig) *Server{
	// server config is nil
	if srvCfg == nil{
		panic("Malatd: srv config is nil.")
	}

	// :8080 -> 0.0.0.0:8080
	if string(srvCfg.ListenAddress[0]) == ":"{
		srvCfg.ListenAddress = "0.0.0.0" + srvCfg.ListenAddress
	}

	// new server
	srv := &Server{
		RouterGroup: RouterGroup{
			Handlers: nil,
			root:     true,
			basePath: "/",
		},
		router: fasthttprouter.New(),
		srvConfig: srvCfg,
	}
	srv.RouterGroup.server = srv
	return srv
}

// ListenAndServe fast http listen
func (srv *Server)ListenAndServe()error{
	xlog.Infof("Server Run %s",srv.srvConfig.ListenAddress)

	// start listen
	if err := fasthttp.ListenAndServe(srv.srvConfig.ListenAddress,srv.router.Handler);err != nil{
		return err
	}
	return nil
}