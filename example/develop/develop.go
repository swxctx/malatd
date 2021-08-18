package main

import (
	"github.com/swxctx/malatd"
)

/*
	http://127.0.0.1:8080/malatd
*/
func main() {
	// new server
	srv := malatd.NewServer(malatd.NewSrvConfig())

	// api router
	srv.Get("/malatd", malatdApi)
	srv.Post("/malatd", malatdApi)
	srv.Run()
}

type User struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

// malatd
func malatdApi(ctx *malatd.Context) {

	ctx.String(200, "malatd")
}
