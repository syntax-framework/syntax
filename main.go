package main

import (
	"github.com/syntax-framework/syntax/admin"
	"log"
	"net/http"
)

// @TODO: https://github.com/fsnotify/fsnotify
func main() {
	handler := admin.SiteAdmin()
	httpAddr := "localhost:8080"
	if err := http.ListenAndServe(httpAddr, handler); err != nil {
		log.Fatalf("ListenAndServe %s: %v", httpAddr, err)
	}
}
