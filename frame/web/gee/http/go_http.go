package http

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

var (
	httpInst *HelloServer
	httpOnce sync.Once
)

type HelloServer struct {}

func GetHelloServer() *HelloServer {
	if httpInst == nil {
		httpOnce.Do(func() {
			httpInst = &HelloServer{}
		})
	}
	return httpInst
}

// init
func (s *HelloServer) init() {
	//http.HandleFunc("/", indexHandler)
	//http.HandleFunc("/header", headerHandler)
	log.Fatal(http.ListenAndServe(":6789", s))
}

func (s *HelloServer) ServeHTTP(w http.ResponseWriter, req *http.Request)  {
	switch req.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	case "/header":
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	default:
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}

// indexHandler returns r.URL.Path
func indexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
}

// handler echoes r.URL.Header
func headerHandler(w http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
}