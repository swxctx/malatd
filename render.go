package td

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

/**
    @date: 2022/3/14
**/

// Render response json
func (c *Context) Render(obj interface{}) (int, error) {
	if c.server != nil && c.server.customRender != nil {
		return c.server.customRender(c, obj)
	}
	resp, _ := EncodeJSON(obj)
	c.ResponseWriter.Header().Set("Content-type", "application/json")
	return c.ResponseWriter.Write(resp)
}

// EncodeJSON marshal json
func EncodeJSON(obj interface{}) ([]byte, error) {
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

// Stream 方法用于发送流式响应
func (c *Context) Stream(step func(w io.Writer) bool) error {
	flusher, ok := c.ResponseWriter.(http.Flusher)
	if !ok {
		return fmt.Errorf("stream: ResponseWriter does not implement http.Flusher")
	}

	for {
		keepOpen := step(c.ResponseWriter)
		if !keepOpen {
			break
		}
		flusher.Flush()
	}
	return nil
}
