package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	server := NewPageServer()
	mux.HandleFunc("/page/", server.pageHandler)
	mux.HandleFunc("/tag/", server.tagHandler)
	mux.HandleFunc("/due/", server.dueHandler)

	log.Fatal(http.ListenAndServe("localhost:8880", mux))
}
