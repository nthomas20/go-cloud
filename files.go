// go get golang.org/x/net/webdav
// go run files.go -openbrowser -http=127.0.0.1:9090
package main

import (
	"context"
	"flag"
	"net/http"
	"log"
	"fmt"
	"net"
	"time"
	"runtime"
	"os"
	"os/exec"
	"golang.org/x/net/webdav"
)

var (
	httpListen  = flag.String("http", "127.0.0.1:8080", "host:port to listen on")
	openBrowser = flag.Bool("openbrowser", false, "open browser automatically")
	g_username  = flag.String("username", "admin", "the default username")
	g_password  = flag.String("password", "password", "the default password")
	root_folder  = flag.String("root", "./data", "the default password")
)

var (
	handler = &webdav.Handler{
		FileSystem: webdav.Dir(*root_folder),
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

func main() {
	flag.Parse()

	if _, err := os.Stat(*root_folder); os.IsNotExist(err) {
		os.Mkdir(*root_folder, 0755)
	}

	server := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		username, password, authOK := r.BasicAuth()

		if authOK == false {
			http.Error(w, "Not authorized", 401)
			return
		}

		if username != *g_username || password != *g_password {
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

	listener, err := net.Listen("tcp", *httpListen)
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		fmt.Println("Started")
		url := "http://" + *httpListen
		if waitServer(url) && *openBrowser && startBrowser(url) {
			log.Printf("A browser window should open. If not, please visit %s", url)
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

// startBrowser tries to open the URL in a browser, and returns
// whether it succeed.
func startBrowser(url string) bool {
	// try to start the browser
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}
