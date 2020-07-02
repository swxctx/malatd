package malatd

import (
	"github.com/swxctx/malatd/binding"
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

// ContentType
func (c *Context) ContentType() string {
	return string(c.Ctx.Request.Header.ContentType())
}

// String response string
func (c *Context) String(code int, msg string) (int, error) {
	c.Ctx.SetStatusCode(code)
	return c.Ctx.WriteString(msg)
}

// Bind
func (c *Context) Bind(obj interface{}) error {
	return binding.Default(c.ContentType()).Bind(c.Ctx, obj)
}
