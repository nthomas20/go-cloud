package webdav

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/webdav"
)

var (
	httpListen = "127.0.0.1:8080"
	username   = "admin"
	password   = "password"
	rootFolder = "./data"
)

var (
	handler = &webdav.Handler{
		FileSystem: webdav.Dir(rootFolder),
		LockSystem: webdav.NewMemLS(),
	}
)

// responseWriterNoBody is a wrapper used to suprress the body of the response
// to a request. Mainly used for HEAD requests.
type responseWriterNoBody struct {
	http.ResponseWriter
}

// newResponseWriterNoBody creates a new responseWriterNoBody.
func newResponseWriterNoBody(w http.ResponseWriter) *responseWriterNoBody {
	return &responseWriterNoBody{w}
}

// Header executes the Header method from the http.ResponseWriter.
func (w responseWriterNoBody) Header() http.Header {
	return w.ResponseWriter.Header()
}

// Write suprresses the body.
func (w responseWriterNoBody) Write(data []byte) (int, error) {
	return 0, nil
}

// WriteHeader writes the header to the http.ResponseWriter.
func (w responseWriterNoBody) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
}

// Run : Run the webdav server
func Run() {
	flag.Parse()

	if _, err := os.Stat(rootFolder); os.IsNotExist(err) {
		os.Mkdir(rootFolder, 0755)
	}

	server := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		authUsername, authPassword, authOK := r.BasicAuth()

		if authOK == false {
			http.Error(w, "Not authorized", 401)
			return
		}

		if authUsername != username || authPassword != password {
			http.Error(w, "Not authorized", 401)
			return
		}

		if r.Method == http.MethodHead {
			w = newResponseWriterNoBody(w)
		}

		if r.Method == http.MethodGet {
			info, err := handler.FileSystem.Stat(context.TODO(), r.URL.Path)
			if err == nil && info.IsDir() {
				r.Method = "PROPFIND"
			}
		}

		handler.ServeHTTP(w, r)
	})

	listener, err := net.Listen("tcp", httpListen)
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		fmt.Println("Started")
		url := "http://" + httpListen
		if waitServer(url) {
			log.Printf("Please visit %s", url)
		} else {
			log.Printf("Please open your web browser and visit %s", url)
		}
	}()
	log.Fatal(http.Serve(listener, server), nil)
}

// waitServer waits some time for the http Server to start
// serving url. The return value reports whether it starts.
func waitServer(url string) bool {
	tries := 20
	for tries > 0 {
		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			return true
		}
		time.Sleep(100 * time.Millisecond)
		tries--
	}
	return false
}
