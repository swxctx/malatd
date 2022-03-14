package td

import (
	"net/http"
	"net/url"
)

type (
	// Context
	Context struct {
		// Request
		Request *http.Request
		// ResponseWriter
		ResponseWriter http.ResponseWriter
		// index
		index int
		// http server
		server *Server
		// server plugins || handler
		plugins Plugins
	}
)

// Plugin
type Plugin func(ctx *Context)

// Plugins
type Plugins []Plugin

// Next
func (c *Context) Next() {
	c.index += 1
	if c.index <= len(c.plugins) {
		c.plugins[c.index-1](c)
	} else {
		c.index = 1
		c.plugins[0](c)
	}
}

// Head
func (c *Context) Head(key string) string {
	return c.Request.Header.Get(key)
}

// HeadSet
func (c *Context) HeadSet(key, value string) {
	c.ResponseWriter.Header().Set(key, value)
}

// ContentType
func (c *Context) ContentType() string {
	return c.Request.Header.Get("Content-Type")
}

// ContentType
func (c *Context) ContentTypeSet(value string) {
	c.ResponseWriter.Header().Set("Content-Type", value)
}

// GetRemoteIP
func (c *Context) GetRemoteIP() string {
	return c.Request.RemoteAddr
}

// QueryValues
func (c *Context) QueryValues() url.Values {
	return c.Request.URL.Query()
}

// Query
func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}
