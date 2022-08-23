package main

import (
	td "github.com/swxctx/malatd"
	"github.com/swxctx/malatd/example/apiserver/api"
)

func main() {
	// new server
	srv := td.NewServer(cfg.SrvConfig)
	api.Route(srv)
	srv.Run()
}
