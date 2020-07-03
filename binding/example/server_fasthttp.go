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
