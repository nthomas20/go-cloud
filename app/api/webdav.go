/*
 * Filename: webdav.go
 * Author: Nathaniel Thomas
 * Many thanks to basic authorization pattern from: https://github.com/valyala/fasthttprouter/blob/master/examples/auth.go
 * netHTTPResponseWriter utilized from: https://github.com/valyala/fasthttp/blob/master/fasthttpadaptor/adaptor.go
 */

package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"github.com/valyala/fasthttprouter"
	"golang.org/x/net/webdav"
)

var (
	basicAuthPrefix = []byte("Basic ")
	handler         = &webdav.Handler{
		// FileSystem: webdav.Dir(rootFolder),
		LockSystem: webdav.NewMemLS(),
	}
)

type netHTTPResponseWriter struct {
	statusCode int
	h          http.Header
	body       []byte
}

func (w *netHTTPResponseWriter) StatusCode() int {
	if w.statusCode == 0 {
		return http.StatusOK
	}
	return w.statusCode
}

func (w *netHTTPResponseWriter) Header() http.Header {
	if w.h == nil {
		w.h = make(http.Header)
	}
	return w.h
}

func (w *netHTTPResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func (w *netHTTPResponseWriter) Write(p []byte) (int, error) {
	w.body = append(w.body, p...)
	return len(p), nil
}

func (w *netHTTPResponseWriter) Body() []byte {
	return w.body
}

func basicAuth(ctx *fasthttp.RequestCtx) (string, string, error) {
	// Get the Basic Authentication credentials
	auth := ctx.Request.Header.Peek("Authorization")
	fmt.Println(string(auth))
	if bytes.HasPrefix(auth, basicAuthPrefix) {
		// Check credentials
		payload, err := base64.StdEncoding.DecodeString(string(auth[len(basicAuthPrefix):]))

		if err == nil {
			pair := bytes.SplitN(payload, []byte(":"), 2)
			if len(pair) == 2 {
				return string(pair[0]), string(pair[1]), nil
			}
		}
	}

	// Request Basic Authentication otherwise
	ctx.Response.Header.Add("WWW-Authenticate", `Basic realm="Restricted"`)
	ctx.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized), fasthttp.StatusUnauthorized)

	return "", "", errors.New("invalid or missing authorization")
}

func (config *Configuration) webdav(ctx *fasthttp.RequestCtx, params fasthttprouter.Params) {
	var (
		webdavResponse netHTTPResponseWriter
		webdavRequest  http.Request
		filepath       = params.ByName("filepath")
	)

	fmt.Println(string(ctx.Request.RequestURI()))
	fmt.Println(filepath)

	if username, password, err := basicAuth(ctx); err == nil {
		fmt.Println(username, password)

		// Check username and password against available configuration
		if _, found := config.Configuration.Accounts[username]; found == false {
			ctx.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized), fasthttp.StatusUnauthorized)
			return
		}

		if _, found := config.Configuration.Accounts[username].Passwords[password]; found == false {
			ctx.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized), fasthttp.StatusUnauthorized)
			return
		}

		// TODO: Check active status of account

		// Configure directory for user
		handler.FileSystem = webdav.Dir(config.Configuration.Accounts[username].RootDirectory)

		// Manage directory properties
		if string(ctx.Request.Header.Method()) == "GET" {
			info, err := handler.FileSystem.Stat(context.TODO(), string(filepath))

			if err == nil && info.IsDir() {
				fmt.Println("DIRECTORY")
				ctx.Request.Header.SetMethod("PROPFIND")
			}
		}

		fmt.Println(string(ctx.Request.Header.Method()))

		// Convert fasthttp request to net/http compatible for webdav server
		if err := fasthttpadaptor.ConvertRequest(ctx, &webdavRequest, true); err != nil {
			ctx.Error(fasthttp.StatusMessage(fasthttp.StatusBadRequest), fasthttp.StatusBadRequest)
		}

		// Set the webdav request to the requested filepath (no prefix)
		webdavRequest.URL.Path = filepath

		// Run the webdav request
		handler.ServeHTTP(&webdavResponse, &webdavRequest)

		// fmt.Println("----------------")
		// fmt.Println("content-type", webdavResponse.Header().Get("content-type"))
		// fmt.Println("content-length", webdavResponse.Header().Get("content-length"))

		// Transition webdav response back to fasthttp
		if webdavResponse.StatusCode() > 299 {
			ctx.Error(fasthttp.StatusMessage(webdavResponse.StatusCode()), webdavResponse.StatusCode())
		} else {
			// Set the webdav server information back to the actual http response
			ctx.Response.SetStatusCode(webdavResponse.StatusCode())
			ctx.Response.SetBody(webdavResponse.Body())
			ctx.Response.Header.Set("content-type", webdavResponse.Header().Get("content-type"))
			ctx.Response.Header.Set("content-length", webdavResponse.Header().Get("content-length"))

			// TODO: Deal with meta data or modifications
		}
	} else {
		ctx.Response.Header.Add("www-authenticate", `Basic realm=Restricted"`)
		fmt.Println(err)
	}

	return
}
