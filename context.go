package td

import (
	"encoding/json"
	"net/http"
	"net/url"
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

// Head
func (c *Context) Head(key string) string {
	return c.Request.Header.Get(key)
}

// ContentType
func (c *Context) ContentType() string {
	return c.Request.Header.Get("Content-Type")
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

// RendString response string
func (c *Context) RendString(code int, msg string) (int, error) {
	c.ResponseWriter.WriteHeader(code)
	return c.ResponseWriter.Write([]byte(msg))
}

// RendJson response json
func (c *Context) RendJson(obj interface{}) (int, error) {
	resp, _ := encodeJSON(obj)
	c.ResponseWriter.Header().Set("content-type","applicaton/json")
	return c.ResponseWriter.Write(resp)
}

// RendRerr response rerror
func (c *Context) RendRerr(rerr *Rerror) (int, error) {
	rerrRsp, _ := rerr.MarshalRerror()
	return c.ResponseWriter.Write(rerrRsp)
}

// encodeJSON
func encodeJSON(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}
