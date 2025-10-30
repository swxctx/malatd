package binding

import (
	"fmt"
	td "github.com/swxctx/malatd"
	"reflect"
)

type queryBinding struct{}

// Name
func (queryBinding) Name() string {
	return "query"
}

// Bind
func (queryBinding) Bind(ctx *td.Context, obj interface{}) error {
	t := reflect.TypeOf(obj)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("queryBinding: obj must be pointer")
	}
	elem := t.Elem()

	if elem.Kind() != reflect.Struct {
		return nil
	}

	if err := mapForm(obj, ctx.Request.URL.Query()); err != nil {
		return err
	}
	return nil
}
