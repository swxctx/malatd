package malatd

import (
	"sync"

	"github.com/valyala/fasthttp"
)

// Context
type Context struct {
	// fast http ctx
	Ctx *fasthttp.RequestCtx
	// index
	i int
	// malatd http server
	server *Server
	data   sync.Map
	// server plugins / handler
	plugins Plugins
}

// Plugin
type Plugin func(ctx *Context)

// Plugins
type Plugins []Plugin

// Next
func (c *Context) Next() {
	c.i += 1
	if c.i <= len(c.plugins) {
		c.plugins[c.i-1](c)
	} else {
		c.i = 1
		c.plugins[0](c)
	}
}

// String
func (c *Context) String(code int, msg string) (int, error) {
	c.Ctx.SetStatusCode(code)
	return c.Ctx.WriteString(msg)
}
