package binding

import (
	"github.com/valyala/fasthttp"
)

const (
	MIMEJSON = "application/json"
)

var (
	JSON_CODEC = jsonCodec{}
)

// Binding
type Binding interface {
	Name() string
	Bind(*fasthttp.RequestCtx, interface{}) error
}

// Default
func Default(contentType string) Binding {
	switch contentType {
	case MIMEJSON:
		return JSON_CODEC
	default:
		return JSON_CODEC
	}
}
