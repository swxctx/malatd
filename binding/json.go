package binding

import (
	"encoding/json"

	"github.com/swxctx/malatd"
)

type jsonBinding struct{}

// Name
func (jsonBinding) Name() string {
	return "json"
}

// Bind
func (jsonBinding) Bind(ctx *malatd.Context, obj interface{}) error {
	return decodeJSON(ctx.CallCtx.Request.Body(), obj)
}

// BindBody
func (jsonBinding) BindBody(body []byte, obj interface{}) error {
	return decodeJSON(body, obj)
}

// decodeJSON
func decodeJSON(data []byte, obj interface{}) error {
	return json.Unmarshal(data, obj)
}
