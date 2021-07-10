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

func (config *Configuration) setupRouter() *fasthttprouter.Router {
	// Configure our Router
	r := fasthttprouter.New()

	// Configure our GET method routes
	// r.GET("/", index)
	// r.GET("/ping", ping)
	// r.GET("/status", status)
	// r.GET("/version", config.version)

	// go-cloud and nextcloud compatible webdav
	for _, route := range []string{"/webdav/*filepath", "/remote.php/dav/files/:username/*filepath"} {
		r.GET(route, config.webdav)
		r.POST(route, config.webdav)
		r.PUT(route, config.webdav)
		r.HEAD(route, config.webdav)
		r.OPTIONS(route, config.webdav)
		r.DELETE(route, config.webdav)
		r.Handle("PROPFIND", route, config.webdav)
		r.Handle("MKCOL", route, config.webdav)
		r.Handle("MOVE", route, config.webdav)
		r.Handle("LOCK", route, config.webdav)
		r.Handle("UNLOCK", route, config.webdav)
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
		handler := fasthttp.CompressHandler(r.Handler)

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
