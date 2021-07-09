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

	r.GET("/webdav", config.webdav)
	r.POST("/webdav", config.webdav)
	r.PUT("/webdav", config.webdav)
	r.HEAD("/webdav", config.webdav)
	r.OPTIONS("/webdav", config.webdav)

	r.GET("/webdav/*filepath", config.webdav)
	r.POST("/webdav/*filepath", config.webdav)
	r.PUT("/webdav/*filepath", config.webdav)
	r.HEAD("/webdav/*filepath", config.webdav)
	r.OPTIONS("/webdav/*filepath", config.webdav)
	r.DELETE("/webdav/*filepath", config.webdav)

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
