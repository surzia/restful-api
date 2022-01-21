package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	server := NewPageServer()
	router.POST("/page/", server.createPageHandler)
	router.GET("/page/", server.getAllPagesHandler)
	router.DELETE("/page/", server.deleteAllPagesHandler)
	router.PUT("/page/", server.updatePageHandler)
	router.GET("/page/{id:[0-9]+}/", server.getPageHandler)
	router.DELETE("/page/{id:[0-9]+}/", server.deletePageHandler)
	router.GET("/tag/", server.tagHandler)
	router.GET("/due/", server.dueHandler)

	log.Fatal(http.ListenAndServe("localhost:8880", router))
}
