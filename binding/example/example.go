package main

import (
	"github.com/swxctx/malatd"
	"github.com/swxctx/malatd/binding"
)

var (
	myApiware = binding.New(nil, nil)
)

func main() {
	// Check whether these structs meet the requirements of apiware, and register them
	err := myApiware.Register(
		new(fasthttpTestApiware),
	)
	if err != nil {
		panic(err)
	}

	// new server
	srv := malatd.NewServer(malatd.NewSrvConfig())

	// api router
	srv.Get("/malatd", testHandler)
	srv.Post("/malatd", testHandler)

	srv.Run()
}
