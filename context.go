package malatd

import (
	"github.com/valyala/fasthttp"
)

type (
	// Context
	Context struct {
		// fast http ctx
		CallCtx *fasthttp.RequestCtx
		// index
		i int
		// malatd http server
		server *Server
		// server plugins / handler
		plugins Plugins
	}
)

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

// ContentType
func (c *Context) ContentType() string {
	return string(c.CallCtx.Request.Header.ContentType())
}

// String response string
func (c *Context) String(code int, msg string) (int, error) {
	c.CallCtx.SetStatusCode(code)
	return c.CallCtx.WriteString(msg)
}
