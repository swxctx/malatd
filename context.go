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
	// handlers
	handlers Handlers
}

// HandlerFunc
type HandlerFunc func(ctx *Context)

// Handlers
type Handlers []HandlerFunc

// Next
func (c *Context) Next() {
	c.i += 1
	if c.i <= len(c.handlers) {
		c.handlers[c.i-1](c)
	} else {
		c.i = 1
		c.handlers[0](c)
	}
}

// String
func (c *Context) String(code int, msg string) (int, error) {
	c.Ctx.SetStatusCode(code)
	return c.Ctx.WriteString(msg)
}
