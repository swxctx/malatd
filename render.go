package td

import "encoding/json"

/**
    @date: 2022/3/14
**/

// Render response json
func (c *Context) Render(obj interface{}) (int, error) {
	resp, _ := encodeJSON(obj)
	c.ResponseWriter.Header().Set("content-type", "applicaton/json")
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
	c.ResponseWriter.Header().Set("content-type", "applicaton/json")
	c.index = abortIndex
	return c.ResponseWriter.Write(rerrRsp)
}
