package binding

import td "github.com/swxctx/malatd"

var (
	binderJson  = JSON
	binderQuery = QUERY
)

// Binder request params binder
func Binder(ctx *td.Context, obj interface{}) error {
	// json
	if err := binderJson.Bind(ctx, obj); err != nil {
		return err
	}
	// query
	if err := binderQuery.Bind(ctx, obj); err != nil {
		return err
	}
	return nil
}
