package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.StrictSlash(true)
	server := NewPageServer()
	router.HandleFunc("/page/", server.createPageHandler).Methods("POST")
	router.HandleFunc("/page/", server.getAllPagesHandler).Methods("GET")
	router.HandleFunc("/page/", server.deleteAllPagesHandler).Methods("DELETE")
	router.HandleFunc("/page/", server.updatePageHandler).Methods("PUT")
	router.HandleFunc("/page/{id:[0-9]+}/", server.getPageHandler).Methods("GET")
	router.HandleFunc("/page/{id:[0-9]+}/", server.deletePageHandler).Methods("DELETE")
	router.HandleFunc("/tag/", server.tagHandler).Methods("GET")
	router.HandleFunc("/due/", server.dueHandler).Methods("GET")

	log.Fatal(http.ListenAndServe("localhost:8880", router))
}
