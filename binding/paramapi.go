package binding

import (
	"errors"
	"fmt"

	"reflect"
	"strconv"
	"sync"

	"github.com/valyala/fasthttp"
)

type (
	// ParamsAPI defines a parameter model for an web api.
	ParamsAPI struct {
		name   string
		params []*Param
		//used to create a new struct (non-pointer)
		structType reflect.Type
		//the raw struct pointer
		rawStructPointer interface{}
		// create param name from struct field name
		paramNameFunc ParamNameFunc
		// decode params from request body
		bodyDecodeFunc BodyDecodeFunc
		//when request Content-Type is multipart/form-data, the max memory for body.
		maxMemory int64
	}

	// Schema is a collection of ParamsAPI
	Schema struct {
		lib map[string]*ParamsAPI
		sync.RWMutex
	}

	// Create param name from struct param name
	ParamNameFunc func(fieldName string) (paramName string)

	// Decode params from request body
	BodyDecodeFunc func(dest reflect.Value, body []byte) error
)

var (
	defaultSchema = &Schema{
		lib: map[string]*ParamsAPI{},
	}
)

// NewParamsAPI parses and store the struct object, requires a struct pointer,
// if `paramNameFunc` is nil, `paramNameFunc=toSnake`,
// if `bodyDecodeFunc` is nil, `bodyDecodeFunc=bodyJONS`,
func NewParamsAPI(
	structPointer interface{},
	paramNameFunc ParamNameFunc,
	bodyDecodeFunc BodyDecodeFunc,
) (
	*ParamsAPI,
	error,
) {
	name := reflect.TypeOf(structPointer).String()
	v := reflect.ValueOf(structPointer)
	if v.Kind() != reflect.Ptr {
		return nil, NewError(name, "*", "the binding object must be a struct pointer")
	}
	v = reflect.Indirect(v)
	if v.Kind() != reflect.Struct {
		return nil, NewError(name, "*", "the binding object must be a struct pointer")
	}
	var m = &ParamsAPI{
		name:             name,
		params:           []*Param{},
		structType:       v.Type(),
		rawStructPointer: structPointer,
	}
	if paramNameFunc != nil {
		m.paramNameFunc = paramNameFunc
	} else {
		m.paramNameFunc = toSnake
	}
	if bodyDecodeFunc != nil {
		m.bodyDecodeFunc = bodyDecodeFunc
	} else {
		m.bodyDecodeFunc = bodyJONS
	}
	err := m.addFields([]int{}, m.structType, v)
	if err != nil {
		return nil, err
	}
	defaultSchema.set(m)
	return m, nil
}

// Register is similar to a `NewParamsAPI`, but only return error.
// Parse and store the struct object, requires a struct pointer,
// if `paramNameFunc` is nil, `paramNameFunc=toSnake`,
// if `bodyDecodeFunc` is nil, `bodyDecodeFunc=bodyJONS`,
func Register(
	structPointer interface{},
	paramNameFunc ParamNameFunc,
	bodyDecodeFunc BodyDecodeFunc,
) error {
	_, err := NewParamsAPI(structPointer, paramNameFunc, bodyDecodeFunc)
	return err
}

func (m *ParamsAPI) addFields(parentIndexPath []int, t reflect.Type, v reflect.Value) error {
	var err error
	var maxMemoryMB int64
	var hasFormData, hasBody bool
	var deep = len(parentIndexPath) + 1
	for i := 0; i < t.NumField(); i++ {
		indexPath := make([]int, deep)
		copy(indexPath, parentIndexPath)
		indexPath[deep-1] = i

		var field = t.Field(i)
		tag, ok := field.Tag.Lookup(TAG_PARAM)
		if !ok {
			if field.Anonymous && field.Type.Kind() == reflect.Struct {
				if err = m.addFields(indexPath, field.Type, v.Field(i)); err != nil {
					return err
				}
			}
			continue
		}

		if tag == TAG_IGNORE_PARAM {
			continue
		}

		if field.Type.Kind() == reflect.Ptr {
			return NewError(t.String(), field.Name, "field can not be a pointer")
		}

		var parsedTags = ParseTags(tag)
		var paramPosition = parsedTags["in"]
		var paramTypeString = field.Type.String()

		switch paramTypeString {
		case fileTypeString:
			if paramPosition != "formData" {
				return NewError(t.String(), field.Name, "when field type is `"+paramTypeString+"`, tag `in` value must be `formData`")
			}
		}

		switch paramPosition {
		case "formData":
			if hasBody {
				return NewError(t.String(), field.Name, "tags of `in(formData)` and `in(body)` can not exist at the same time")
			}
			hasFormData = true
		case "body":
			if hasFormData {
				return NewError(t.String(), field.Name, "tags of `in(formData)` and `in(body)` can not exist at the same time")
			}
			if hasBody {
				return NewError(t.String(), field.Name, "there should not be more than one tag `in(body)`")
			}
			hasBody = true
		default:
			if !TagInValues[paramPosition] {
				return NewError(t.String(), field.Name, "invalid tag `in` value, refer to the following: `query`, `formData`, `body`, or `header`")
			}
		}
		if _, ok := parsedTags["len"]; ok && paramTypeString != "string" && paramTypeString != "[]string" {
			return NewError(t.String(), field.Name, "invalid `len` tag for non-string field")
		}
		if _, ok := parsedTags["range"]; ok {
			switch paramTypeString {
			case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float32", "float64":
			case "[]int", "[]int8", "[]int16", "[]int32", "[]int64", "[]uint", "[]uint8", "[]uint16", "[]uint32", "[]uint64", "[]float32", "[]float64":
			default:
				return NewError(t.String(), field.Name, "invalid `range` tag for non-number field")
			}
		}
		if a, ok := field.Tag.Lookup(TAG_REGEXP); ok {
			if paramTypeString != "string" && paramTypeString != "[]string" {
				return NewError(t.String(), field.Name, "invalid `"+TAG_REGEXP+"` tag for non-string field")
			}
			parsedTags[TAG_REGEXP] = a
		}
		if a, ok := parsedTags["maxmb"]; ok {
			i, err := strconv.ParseInt(a, 10, 64)
			if err != nil {
				return NewError(t.String(), field.Name, "invalid `maxmb` tag, it must be positive integer")
			}
			if i > maxMemoryMB {
				maxMemoryMB = i
			}
		}

		fd := &Param{
			apiName:   m.name,
			indexPath: indexPath,
			tags:      parsedTags,
			rawTag:    field.Tag,
			rawValue:  v.Field(i),
		}

		if errStr, ok := field.Tag.Lookup(TAG_ERR); ok {
			fd.tags[TAG_ERR] = errStr
			fd.err = errors.New(errStr)
		}

		if fd.name, ok = parsedTags["name"]; !ok {
			fd.name = m.paramNameFunc(field.Name)
		}

		fd.isFile = paramTypeString == fileTypeString
		_, fd.isRequired = parsedTags["required"]

		m.params = append(m.params, fd)
	}
	if maxMemoryMB > 0 {
		m.maxMemory = maxMemoryMB * MB
	} else {
		m.maxMemory = defaultMaxMemory
	}
	return nil
}

// GetParamsAPI gets the `*ParamsAPI` object according to the type name
func GetParamsAPI(paramsAPIName string) (*ParamsAPI, error) {
	m, ok := defaultSchema.get(paramsAPIName)
	if !ok {
		return nil, errors.New("struct `" + paramsAPIName + "` is not registered")
	}
	return m, nil
}

// SetParamsAPI caches `*ParamsAPI`
func SetParamsAPI(m *ParamsAPI) {
	defaultSchema.set(m)
}

func (schema *Schema) get(paramsAPIName string) (*ParamsAPI, bool) {
	schema.RLock()
	defer schema.RUnlock()
	m, ok := schema.lib[paramsAPIName]
	return m, ok
}

func (schema *Schema) set(m *ParamsAPI) {
	schema.Lock()
	schema.lib[m.name] = m
	defer schema.Unlock()
}

// Name gets the name
func (paramsAPI *ParamsAPI) Name() string {
	return paramsAPI.name
}

// Params gets the parameter information
func (paramsAPI *ParamsAPI) Params() []*Param {
	return paramsAPI.params
}

// Number returns the number of parameters to be bound
func (paramsAPI *ParamsAPI) Number() int {
	return len(paramsAPI.params)
}

// Raw returns the ParamsAPI's original value
func (paramsAPI *ParamsAPI) Raw() interface{} {
	return paramsAPI.rawStructPointer
}

// MaxMemory gets maxMemory
// when request Content-Type is multipart/form-data, the max memory for body.
func (paramsAPI *ParamsAPI) MaxMemory() int64 {
	return paramsAPI.maxMemory
}

// SetMaxMemory sets maxMemory for the request which Content-Type is multipart/form-data.
func (paramsAPI *ParamsAPI) SetMaxMemory(maxMemory int64) {
	paramsAPI.maxMemory = maxMemory
}

// NewReceiver creates a new struct pointer and the field's values  for its receive parameterste it.
func (paramsAPI *ParamsAPI) NewReceiver() (interface{}, []reflect.Value) {
	object := reflect.New(paramsAPI.structType)
	return object.Interface(), paramsAPI.fieldsForBinding(object.Elem())
}

func (paramsAPI *ParamsAPI) fieldsForBinding(structElem reflect.Value) []reflect.Value {
	count := len(paramsAPI.params)
	fields := make([]reflect.Value, count)
	for i := 0; i < count; i++ {
		value := structElem
		param := paramsAPI.params[i]
		for _, index := range param.indexPath {
			value = value.Field(index)
		}
		fields[i] = value
	}
	return fields
}

// BindByName binds the net/http request params to a new struct and validate it.
func BindByName(
	paramsAPIName string,
	req *fasthttp.RequestCtx,
) (
	interface{},
	error,
) {
	paramsAPI, err := GetParamsAPI(paramsAPIName)
	if err != nil {
		return nil, err
	}
	return paramsAPI.BindNew(req)
}

// Bind binds the net/http request params to the `structPointer` param and validate it.
// note: structPointer must be struct pointer.
func Bind(
	structPointer interface{},
	req *fasthttp.RequestCtx,
) error {
	paramsAPI, err := GetParamsAPI(reflect.TypeOf(structPointer).String())
	if err != nil {
		return err
	}
	return paramsAPI.BindAt(structPointer, req)
}

// BindAt binds the net/http request params to a struct pointer and validate it.
// note: structPointer must be struct pointer.
func (paramsAPI *ParamsAPI) BindAt(
	structPointer interface{},
	req *fasthttp.RequestCtx,
) error {
	name := reflect.TypeOf(structPointer).String()
	if name != paramsAPI.name {
		return errors.New("the structPointer's type `" + name + "` does not match type `" + paramsAPI.name + "`")
	}
	return paramsAPI.BindFields(
		paramsAPI.fieldsForBinding(reflect.ValueOf(structPointer).Elem()),
		req,
	)
}

// BindNew binds the net/http request params to a struct pointer and validate it.
func (paramsAPI *ParamsAPI) BindNew(
	req *fasthttp.RequestCtx,
) (
	interface{},
	error,
) {
	structPrinter, fields := paramsAPI.NewReceiver()
	err := paramsAPI.BindFields(fields, req)
	return structPrinter, err
}

// RawBind binds the net/http request params to the original struct pointer and validate it.
func (paramsAPI *ParamsAPI) RawBind(
	req *fasthttp.RequestCtx,
) (
	interface{},
	error,
) {
	var fields []reflect.Value
	for _, param := range paramsAPI.params {
		fields = append(fields, param.rawValue)
	}
	err := paramsAPI.BindFields(fields, req)
	return paramsAPI.rawStructPointer, err
}

// BindFields binds the net/http request params to a struct and validate it.
// Must ensure that the param `fields` matches `paramsAPI.params`.
func (paramsAPI *ParamsAPI) BindFields(
	fields []reflect.Value,
	req *fasthttp.RequestCtx,
) (
	err error,
) {
	defer func() {
		if p := recover(); p != nil {
			err = NewError(paramsAPI.name, "?", fmt.Sprint(p))
		}
	}()

	var formValues = formValues(req)
	for i, param := range paramsAPI.params {
		value := fields[i]
		switch param.In() {
		case "query":
			paramValuesBytes := req.QueryArgs().PeekMulti(param.name)
			if len(paramValuesBytes) > 0 {
				var paramValues = make([]string, len(paramValuesBytes))
				for i, b := range paramValuesBytes {
					paramValues[i] = string(b)
				}
				if err = convertAssign(value, paramValues); err != nil {
					return param.myError(err.Error())
				}
			} else if len(paramValuesBytes) == 0 && param.IsRequired() {
				return param.myError("missing query param")
			}

		case "formData":
			// Can not exist with `body` param at the same time
			if param.IsFile() {
				if fh, err := req.FormFile(param.name); err == nil {
					value.Set(reflect.ValueOf(fh).Elem())
				} else if param.IsRequired() {
					return param.myError("missing formData param")
				}
				continue
			}

			paramValues, ok := formValues[param.name]
			if ok {
				if err = convertAssign(value, paramValues); err != nil {
					return param.myError(err.Error())
				}
			} else if param.IsRequired() {
				return param.myError("missing formData param")
			}

		case "body":
			// Theoretically there should be at most one `body` param, and can not exist with `formData` at the same time
			body := req.PostBody()
			if body != nil {
				if err = paramsAPI.bodyDecodeFunc(value, body); err != nil {
					return param.myError(err.Error())
				}
			} else if param.IsRequired() {
				return param.myError("missing body param")
			}

		case "header":
			paramValueBytes := req.Request.Header.Peek(param.name)
			if paramValueBytes != nil {
				if err = convertAssign(value, []string{string(paramValueBytes)}); err != nil {
					return param.myError(err.Error())
				}
			} else if param.IsRequired() {
				return param.myError("missing header param")
			}
		}
		if err = param.validate(value); err != nil {
			return err
		}
	}
	return
}

// formValues returns all post data values with their keys
// multipart, formValues data, post arguments
func formValues(req *fasthttp.RequestCtx) map[string][]string {
	// first check if we have multipart formValues
	multipartForm, err := req.MultipartForm()
	if err == nil {
		//we have multipart formValues
		return multipartForm.Value
	}
	valuesAll := make(map[string][]string)
	// if no multipart and post arguments ( means normal formValues   )
	if req.PostArgs().Len() == 0 {
		return valuesAll // no found
	}
	req.PostArgs().VisitAll(func(k []byte, v []byte) {
		key := string(k)
		value := string(v)
		// for slices
		if valuesAll[key] != nil {
			valuesAll[key] = append(valuesAll[key], value)
		} else {
			valuesAll[key] = []string{value}
		}
	})
	return valuesAll
}
