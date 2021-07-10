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
	"strings"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"github.com/valyala/fasthttprouter"
	"golang.org/x/net/webdav"
)

var (
	handler = &webdav.Handler{
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
	var (
		basicAuthPrefix = []byte("Basic ")
	)

	// Get the Basic Authentication credentials
	auth := ctx.Request.Header.Peek("Authorization")
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
	ctx.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized), fasthttp.StatusUnauthorized)

	return "", "", errors.New("invalid or missing authorization")
}

func (config *Configuration) webdav(ctx *fasthttp.RequestCtx, params fasthttprouter.Params) {
	var (
		webdavResponse netHTTPResponseWriter
		webdavRequest  http.Request
		filepath       = params.ByName("filepath")
		filename       = ""
		prefix         = string(ctx.Request.Header.Peek("x-webdav-prefix"))
	)

	fmt.Println(string(ctx.Request.Header.Method()), filepath)

	ctx.Response.Header.Add("WWW-Authenticate", "Basic realm=Restricted")

	if username, password, err := basicAuth(ctx); err == nil {
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

		// Get file/directory info
		info, err := handler.FileSystem.Stat(context.TODO(), string(filepath))

		fmt.Println("-=-=-=-=-=-=-=", filepath, info)

		if info == nil {
			ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)
			return
		} else if err != nil {
			ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		// Manage filename and directory situation
		if info.IsDir() == true {
			if string(ctx.Request.Header.Method()) == fasthttp.MethodGet {
				ctx.Request.Header.SetMethod("PROPFIND")
			}
		} else {
			filename = info.Name()
		}

		// Convert fasthttp request to net/http compatible for webdav server
		if err := fasthttpadaptor.ConvertRequest(ctx, &webdavRequest, true); err != nil {
			ctx.Error(fasthttp.StatusMessage(fasthttp.StatusBadRequest), fasthttp.StatusBadRequest)
		}

		// Set the webdav request to the requested filepath (no prefix)
		webdavRequest.URL.Path = filepath
		webdavRequest.RequestURI = "http://" + webdavRequest.Host + filepath

		// Correct the Destination header
		if len(webdavRequest.Header.Get("Destination")) > 0 {
			webdavRequest.Header.Set("Destination", strings.Replace(webdavRequest.Header.Get("Destination"), strings.Replace(prefix, ":username", username, 1), "", 1))
		}

		fmt.Println(webdavRequest)

		// Run the webdav request
		handler.ServeHTTP(&webdavResponse, &webdavRequest)

		// Transition webdav response back to fasthttp
		if webdavResponse.StatusCode() > 299 {
			ctx.Error(fasthttp.StatusMessage(webdavResponse.StatusCode()), webdavResponse.StatusCode())
		} else {
			// Set the webdav server information back to the actual http response
			ctx.Response.SetStatusCode(webdavResponse.StatusCode())

			if string(ctx.Request.Header.Method()) != fasthttp.MethodHead {
				ctx.Response.SetBody(webdavResponse.Body())
			}

			for k, v := range map[string]string{
				fasthttp.HeaderContentType:           webdavResponse.Header().Get(fasthttp.HeaderContentType),
				fasthttp.HeaderContentLength:         webdavResponse.Header().Get(fasthttp.HeaderContentLength),
				fasthttp.HeaderExpires:               "Thu, 19 Nov 1981 08:52:00 GMT",
				fasthttp.HeaderCacheControl:          "no-store, no-cache, must-revalidate",
				fasthttp.HeaderPragma:                "no-cache",
				fasthttp.HeaderContentSecurityPolicy: "default-src 'none';",
			} {
				ctx.Response.Header.Set(k, v)
			}

			// If it's a file, add the following
			if filename != "" {
				ctx.Response.Header.Set(fasthttp.HeaderETag, webdavResponse.Header().Get(fasthttp.HeaderETag))
				ctx.Response.Header.Set(fasthttp.HeaderLastModified, webdavResponse.Header().Get(fasthttp.HeaderLastModified))
				ctx.Response.Header.Set(fasthttp.HeaderContentDisposition, `attachment; filename="`+filename+`"`)
				ctx.Response.Header.Set("filename", `"`+filename+`"`)
			}
		}
	} else {
		ctx.Response.Header.Add("www-authenticate", `Basic realm=Restricted"`)
		fmt.Println(err)
	}

	return
}
