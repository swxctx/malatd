# Binding
Base By [henrylee2cn Apiware](https://github.com/henrylee2cn/apiware)

将`net/http`及`fasthttp`请求的指定参数绑定到结构体，并验证参数值的合法性。
建议您可以使用结构体作为web框架的Handler，并用该中间件快速绑定请求参数，节省了大量参数类型转换与有效性验证的工作。

# Demo 示例

```
package main

import (
	"github.com/swxctx/malatd"
	"github.com/swxctx/malatd/binding"
)

var (
	myApiware = binding.New(nil, nil)
)

func main() {
	// Check whether these structs meet the requirements of apiware, and register them
	err := myApiware.Register(
		new(fasthttpTestApiware),
	)
	if err != nil {
		panic(err)
	}

	// new server
	srv := malatd.NewServer(malatd.NewSrvConfig())

	// api router
	srv.Get("/malatd", testHandler)
	srv.Post("/malatd", testHandler)

	srv.Run()
}

package main

import (
	"encoding/json"
	"net/http"

	"github.com/swxctx/malatd"
)

// fasthttpTestApiware
type fasthttpTestApiware struct {
	Num   float32 `param:"in(query),name(n),range(0.1:10.19)" json:"num"`
	Title string  `param:"in(query),nonzero" json:"title"`
	Desc  Desc    `param:"in(body),name(desc)" json:"desc"`
}

// Desc
type Desc struct {
	A string `json:"a"`
}

// testHandler
func testHandler(ctx *malatd.Context) {
	// bind params
	params := new(fasthttpTestApiware)
	err := myApiware.Bind(params, ctx.CallCtx)
	b, _ := json.MarshalIndent(params, "", " ")

	if err != nil {
		ctx.CallCtx.SetStatusCode(http.StatusBadRequest)
		ctx.CallCtx.Write(append([]byte(err.Error()+"\n"), b...))
	} else {
		ctx.CallCtx.SetStatusCode(http.StatusOK)
		ctx.CallCtx.Write(b)
	}
}
```

# Struct&Tag 结构体及其标签

tag   |   key    | required |     value     |   desc
------|----------|----------|---------------|----------------------------------
param |    in    | only one |     query     | (position of param) e.g. url: `http://www.abc.com/a?b={query}`
param |    in    | only one |     formData  | (position of param) e.g. `request body: a=123&b={formData}`
param |    in    | only one |     body      | (position of param) request body can be any content
param |    in    | only one |     header    | (position of param) request header info
param |   name   |    no    |   (e.g. `id`)  | specify request param`s name
param | required |    no    |    required   | request param is required
param |   desc   |    no    |   (e.g. `id`)  | request param description
param |   len    |    no    | (e.g. `3:6``3`) | length range of param's value
param |   range  |    no    | (e.g. `0:10`)  | numerical range of param's value
param |  nonzero |    no    |    nonzero    | param`s value can not be zero
param |   maxmb  |    no    |  (e.g. `32`)   | when request Content-Type is multipart/form-data, the max memory for body.(multi-param, whichever is greater)
regexp|          |    no    |(e.g. `^\w+$`)| param value can not be null
err   |          |    no    |(e.g. `incorrect password format`)| the custom error for binding or validating
**NOTES**:
* the binding object must be a struct pointer
* the binding struct's field can not be a pointer
* `regexp` or `param` tag is only usable when `param:"type(xxx)"` is exist
* if the `param` tag is not exist, anonymous field will be parsed
* when the param's position(`in`) is `formData` and the field's type is `multipart.FileHeader`, the param receives file uploaded
* param tags `in(formData)` and `in(body)` can not exist at the same time
* there should not be more than one `in(body)` param tag

# Field Types 结构体字段类型

base    |   slice    | special
--------|------------|-------------------------------------------------------
string  |  []string  | [][]byte
byte    |  []byte    | [][]uint8
uint8   |  []uint8   | multipart.FileHeader (only for `formData` param)
bool    |  []bool    | http.Cookie (only for `net/http`'s `cookie` param)
int     |  []int     | fasthttp.Cookie (only for `fasthttp`'s `cookie` param)
int8    |  []int8    | struct (struct type only for `body` param or as an anonymous field to extend params)
int16   |  []int16   |
int32   |  []int32   |
int64   |  []int64   |
uint8   |  []uint8   |
uint16  |  []uint16  |
uint32  |  []uint32  |
uint64  |  []uint64  |
float32 |  []float32 |
float64 |  []float64 |
