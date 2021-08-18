package binding

import (
	"github.com/swxctx/malatd"
	"net/http"
)

type queryBinding struct{}

// Name
func (queryBinding) Name() string {
	return "query"
}

// Bind
func (queryBinding) Bind(ctx *malatd.Context, obj interface{}) error {
	values := ctx.CallCtx.QueryArgs().
	if err := mapForm(obj, values); err != nil {
		return err
	}
	return nil
}