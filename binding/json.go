package binding

import (
	"bytes"
	"encoding/json"
	"io"

	td "github.com/swxctx/malatd"
)

var (
	EnableDecoderUseNumber = false
)

type jsonBinding struct{}

// Name
func (jsonBinding) Name() string {
	return "json"
}

// Bind
func (jsonBinding) Bind(ctx *td.Context, obj interface{}) error {
	return decodeJSON(ctx.Request.Body, obj)
}

func (jsonBinding) BindBody(body []byte, obj interface{}) error {
	return decodeJSON(bytes.NewReader(body), obj)
}

func decodeJSON(r io.Reader, obj interface{}) error {
	decoder := json.NewDecoder(r)
	if EnableDecoderUseNumber {
		decoder.UseNumber()
	}
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return nil
}
