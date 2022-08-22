package main

import (
	td "github.com/swxctx/malatd"
	"github.com/swxctx/malatd/example/apiserver/api"
)

func main() {
	// new server
	srv := td.NewServer(td.NewSrvConfig())
	api.Route(srv)
	srv.Run()
}
