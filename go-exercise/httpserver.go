package main

import (
	"fmt"
	"io"
	"net/http"
)

func goHello(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "go hello------------")
}

type hello struct {
}

func (h *hello) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Method, req.URL.Path)
	io.WriteString(w, "hello world\n")
}

func main() {
	http.HandleFunc("/hello1111", goHello)
	s := &http.Server{
		Addr:           "localhost:6002",
		Handler:        new(hello),
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()
}
