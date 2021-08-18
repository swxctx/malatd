package binding

import td "github.com/swxctx/malatd"

type queryBinding struct{}

// Name
func (queryBinding) Name() string {
	return "query"
}

// Bind
func (queryBinding) Bind(ctx *td.Context, obj interface{}) error {
	if err := mapForm(obj, ctx.Request.URL.Query()); err != nil {
		return err
	}
	return nil
}
