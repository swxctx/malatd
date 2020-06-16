package main

import (
	"github.com/swxctx/malatd"
)

func main() {
	// new server
	srv := malatd.NewServer(malatd.NewSrvConfig())

	// api router
	srv.Get("/malatd", func(ctx *malatd.Context) {
		ctx.String(200, "hello malatd")
	})
	srv.ListenAndServe()
}
