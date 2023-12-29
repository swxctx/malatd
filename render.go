package td

import (
	"encoding/json"
	"net/http"
)

/**
    @date: 2022/3/14
**/

// Render response json
func (c *Context) Render(obj interface{}) (int, error) {
	resp, _ := encodeJSON(obj)
	c.ResponseWriter.Header().Set("Content-type", "application/json")
	return c.ResponseWriter.Write(resp)
}

// encodeJSON
func encodeJSON(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

// Render response json
func (c *Context) RenderString(resp string) (int, error) {
	return c.ResponseWriter.Write([]byte(resp))
}

// RenderRerr response rerror
func (c *Context) RenderRerr(rerr *Rerror) (int, error) {
	rerrRsp, _ := rerr.MarshalRerror()
	c.ResponseWriter.Header().Set("Content-type", "application/json")
	c.index = abortIndex
	return c.ResponseWriter.Write(rerrRsp)
}

// Redirect request redirect
func (c *Context) Redirect(targetUrl string, code int) {
	http.Redirect(c.ResponseWriter, c.Request, targetUrl, code)
}

// renderNotFound
func renderNotFound(response http.ResponseWriter, request *http.Request) {
	notFoundResp, _ := RerrNotFound.MarshalRerror()
	response.Header().Set("Content-type", "application/json")
	response.Write(notFoundResp)
}
