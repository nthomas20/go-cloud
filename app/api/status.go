/*
 * Filename: status.go
 * Author: Nathaniel Thomas
 */

package api

import (
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
)

func index(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	ctx.WriteString("GO Cloud")
}

func ping(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	ctx.WriteString("OK")
}

func status(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	ctx.WriteString("OK")
}

func (config *Configuration) version(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	ctx.WriteString("Version:    " + config.Version + "\n")
	ctx.WriteString("Build Date: " + config.BuildDate)
}
