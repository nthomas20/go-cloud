/*
 *
 */

package api

import (
	"log"
	"os"

	"gitea.nthomas20.net/nathaniel/go-cloud/app/models"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
	"github.com/valyala/fasthttprouter"
)

// API : Export of API Interface
type API interface {
	// Primary API Runners
	setupRouter() *fasthttprouter.Router
	Run() bool
}

// Configuration : API runtime Configuration
type Configuration struct {
	Configuration *models.Configuration
	Version       string
	BuildDate     string
}

func trackPrefix(prefix string, next fasthttprouter.Handle) fasthttprouter.Handle {
	return fasthttprouter.Handle(func(ctx *fasthttp.RequestCtx, params fasthttprouter.Params) {
		// Set prefix header
		ctx.Request.Header.Add("x-webdav-prefix", prefix)

		// Call the next step
		next(ctx, params)
	})
}

func cors(next fasthttprouter.Handle) fasthttprouter.Handle {
	var (
		corsAllowHeaders     = "authorization"
		corsAllowMethods     = "HEAD,GET,POST,PUT,DELETE,OPTIONS,PROPFIND,MKCOL,MOVE,LOCK,UNLOCK"
		corsAllowOrigin      = "*"
		corsAllowCredentials = "true"
	)

	return fasthttprouter.Handle(func(ctx *fasthttp.RequestCtx, params fasthttprouter.Params) {
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", corsAllowCredentials)
		ctx.Response.Header.Set("Access-Control-Allow-Headers", corsAllowHeaders)
		ctx.Response.Header.Set("Access-Control-Allow-Methods", corsAllowMethods)
		ctx.Response.Header.Set("Access-Control-Allow-Origin", corsAllowOrigin)

		// Call the next step
		next(ctx, params)
	})
}

func (config *Configuration) setupRouter() *fasthttprouter.Router {
	// Configure our Router
	r := fasthttprouter.New()
	r.RedirectFixedPath = false
	r.RedirectTrailingSlash = false

	// Configure our GET method routes
	// r.GET("/", index)
	// r.GET("/ping", ping)
	// r.GET("/status", status)
	// r.GET("/version", config.version)

	// go-cloud and nextcloud compatible webdav
	// https://cs.opensource.google/go/x/net/+/04defd46:webdav/webdav.go;l=42
	for _, prefix := range []string{"/webdav", "/remote.php/dav/files/:username"} {
		route := prefix + "/*filepath"

		r.GET(route, trackPrefix(prefix, config.webdav))
		r.POST(route, trackPrefix(prefix, config.webdav))
		r.PUT(route, trackPrefix(prefix, config.webdav))
		r.HEAD(route, trackPrefix(prefix, config.webdav))
		r.OPTIONS(route, trackPrefix(prefix, config.webdav))
		r.DELETE(route, trackPrefix(prefix, config.webdav))
		r.Handle("PROPFIND", route, trackPrefix(prefix, config.webdav))
		r.Handle("PROPPATCH", route, trackPrefix(prefix, config.webdav))
		r.Handle("MKCOL", route, trackPrefix(prefix, config.webdav))
		r.Handle("COPY", route, trackPrefix(prefix, config.webdav))
		r.Handle("MOVE", route, trackPrefix(prefix, config.webdav))
		r.Handle("LOCK", route, trackPrefix(prefix, config.webdav))
		r.Handle("UNLOCK", route, trackPrefix(prefix, config.webdav))
	}

	return r
}

// Run : Run the server
func (config *Configuration) Run() bool {
	log.Println("Configuring Routes")

	r := config.setupRouter()

	log.Println("Configuring TCP Listener")

	// Launch our listener!
	go func() {
		// Trigger compression handler
		// handler := fasthttp.CompressHandler(r.Handler)
		handler := r.Handler

		ln, err := reuseport.Listen("tcp4", ":"+config.Configuration.Port)
		if err != nil {
			log.Println("Error in REUSEPORT listener:", err)
			if err := fasthttp.ListenAndServe(":"+config.Configuration.Port, handler); err != nil {
				log.Println("Error in fasthttp Server:", err)
				os.Exit(1)
			}
		} else {
			if err = fasthttp.Serve(ln, handler); err != nil {
				log.Println("Error in fasthttp Server:", err)
				os.Exit(1)
			}
		}
	}()

	return true
}
