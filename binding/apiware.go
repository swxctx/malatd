package binding

import (
	"errors"

	"github.com/valyala/fasthttp"
)

type (
	Apiware struct {
		ParamNameFunc
		BodyDecodeFunc
	}
)

// Create a new apiware engine.
// Parse and store the struct object, requires a struct pointer,
// if `paramNameFunc` is nil, `paramNameFunc=toSnake`,
// if `bodyDecodeFunc` is nil, `bodyDecodeFunc=bodyJONS`,
func New(bodyDecodeFunc BodyDecodeFunc, paramNameFunc ParamNameFunc) *Apiware {
	return &Apiware{
		ParamNameFunc:  paramNameFunc,
		BodyDecodeFunc: bodyDecodeFunc,
	}
}

// Check whether structs meet the requirements of apiware, and register them.
// note: requires a structure pointer.
func (a *Apiware) Register(structPointers ...interface{}) error {
	var errStr string
	for _, obj := range structPointers {
		err := Register(obj, a.ParamNameFunc, a.BodyDecodeFunc)
		if err != nil {
			errStr += err.Error() + "\n"
		}
	}
	if len(errStr) > 0 {
		return errors.New(errStr)
	}
	return nil
}

// Bind the fasthttp request params to the structure and validate.
// note: structPointer must be structure pointer.
func (a *Apiware) Bind(
	structPointer interface{},
	reqCtx *fasthttp.RequestCtx,
) error {
	return Bind(structPointer, reqCtx)
}
