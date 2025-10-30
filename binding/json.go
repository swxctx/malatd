package binding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

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
	ct := ctx.ContentType()
	rv := reflect.ValueOf(obj)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("Bind: obj must be a pointer")
	}

	if rv.Elem().Kind() == reflect.Map {
		elem := rv.Elem()
		if elem.IsNil() {
			elem.Set(reflect.MakeMap(elem.Type()))
		}

		switch {
		case strings.HasPrefix(ct, "application/x-www-form-urlencoded"):
			if err := ctx.Request.ParseForm(); err != nil {
				return err
			}
			for k, v := range ctx.Request.Form {
				if len(v) > 0 {
					elem.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v[0]))
				}
			}
			return nil

		case strings.HasPrefix(ct, "multipart/form-data"):
			if err := ctx.Request.ParseMultipartForm(32 << 20); err != nil {
				return err
			}
			for k, v := range ctx.Request.MultipartForm.Value {
				if len(v) > 0 {
					elem.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v[0]))
				}
			}
			return nil

		case strings.HasPrefix(ct, "application/json"):
			return decodeJSON(ctx.Request.Body, obj)

		default:
			return decodeJSON(ctx.Request.Body, obj)
		}
	}

	switch {
	case strings.HasPrefix(ct, "application/json"):
		return decodeJSON(ctx.Request.Body, obj)

	case strings.HasPrefix(ct, "application/x-www-form-urlencoded"):
		if err := ctx.Request.ParseForm(); err != nil {
			return err
		}
		formMap := make(map[string]string)
		for k, v := range ctx.Request.Form {
			if len(v) > 0 {
				formMap[k] = v[0]
			}
		}
		b, _ := json.Marshal(formMap)
		return json.Unmarshal(b, obj)

	case strings.HasPrefix(ct, "multipart/form-data"):
		if err := ctx.Request.ParseMultipartForm(32 << 20); err != nil {
			return err
		}
		formMap := make(map[string]string)
		for k, v := range ctx.Request.MultipartForm.Value {
			if len(v) > 0 {
				formMap[k] = v[0]
			}
		}
		b, _ := json.Marshal(formMap)
		return json.Unmarshal(b, obj)

	default:
		return decodeJSON(ctx.Request.Body, obj)
	}
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
		if err == io.EOF {
			return nil
		}
		return err
	}
	return nil
}
