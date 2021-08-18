package td

import (
	"encoding/json"
	"net/http"
)

type (
	// Context
	Context struct {
		// fast http ctx
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

// ContentType
func (c *Context) ContentType() string {
	return c.Request.Header.Get("Content-Type")
}

// String response string
func (c *Context) String(code int, msg string) (int, error) {
	c.ResponseWriter.WriteHeader(code)
	return c.ResponseWriter.Write([]byte(msg))
}

// RenderJson response json
func (c *Context) Json(obj interface{}) (int, error) {
	resp, _ := encodeJSON(obj)
	return c.ResponseWriter.Write(resp)
}

// encodeJSON
func encodeJSON(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}
