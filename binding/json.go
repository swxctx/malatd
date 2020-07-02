package binding

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

type (
	jsonCodec struct{}
)

func (jsonCodec) Name() string {
	return "json"
}

// Bind
func (jsonCodec) Bind(ctx *fasthttp.RequestCtx, obj interface{}) error {
	return decodeJSON(ctx.PostBody(), obj)
}

// BindBody
func (jsonCodec) BindBody(body []byte, obj interface{}) error {
	return decodeJSON(body, obj)
}

// decodeJSON
func decodeJSON(r []byte, obj interface{}) error {
	if err := json.Unmarshal(r, obj); err != nil {
		return err
	}
	return nil
}
