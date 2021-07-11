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
	return
}

func ping(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	ctx.WriteString("OK")
	return
}

func status(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	ctx.WriteString("OK")
	return
}

func (config *Configuration) version(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	ctx.WriteString("Version:    " + config.Version + "\n")
	ctx.WriteString("Build Date: " + config.BuildDate)
	return
}
