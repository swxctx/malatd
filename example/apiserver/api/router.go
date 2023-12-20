package api

import (
	td "github.com/swxctx/malatd"
	"github.com/swxctx/malatd/example/apiserver/plugin"
)

func Route(srv *td.Server) {
	srv.Get("/", indexHandle)

	login := srv.Group("/user")
	login.Post("/login", plugin.AuthPlugin, loginHandle)
	login.Get("/login", plugin.AuthPlugin, loginHandle)
}
