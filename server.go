package main

import (
	"net/http"
	"restful-api/page"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type PageServer struct {
	store *page.Store
}

type PageRequest struct {
	Text string    `json:"text"`
	Tags []string  `json:"tags"`
	Due  time.Time `json:"due"`
}

func NewPageServer() *PageServer {
	return &PageServer{store: page.New()}
}

func (p *PageServer) tagHandler(c *gin.Context) {
	tag := c.Params.ByName("tag")
	tasks := p.store.GetPagesByTag(tag)
	c.JSON(http.StatusOK, tasks)
}

func (p *PageServer) dueHandler(c *gin.Context) {
	badRequestError := func() {
		c.String(http.StatusBadRequest, "expect /due/<year>/<month>/<day>, got %v", c.FullPath())
	}

	year, err := strconv.Atoi(c.Params.ByName("year"))
	if err != nil {
		badRequestError()
		return
	}

	month, err := strconv.Atoi(c.Params.ByName("month"))
	if err != nil || month < int(time.January) || month > int(time.December) {
		badRequestError()
		return
	}

	day, err := strconv.Atoi(c.Params.ByName("day"))
	if err != nil {
		badRequestError()
		return
	}

	tasks := p.store.GetPagesByDueDate(year, time.Month(month), day)
	c.JSON(http.StatusOK, tasks)
}

func (p *PageServer) createPageHandler(c *gin.Context) {
	var ret PageRequest
	if err := c.ShouldBindJSON(&ret); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	id := p.store.CreatePage(ret.Text, ret.Tags, ret.Due)
	c.JSON(http.StatusOK, gin.H{"Id": id})
}

func (p *PageServer) updatePageHandler(c *gin.Context) {
	var ret page.Page
	if err := c.ShouldBindJSON(&ret); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	_, err := p.store.UpdatePage(&ret)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}
	c.JSON(http.StatusOK, gin.H{"msg": "update page succeed!"})
}

func (p *PageServer) deletePageHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if err = p.store.DeletePage(id); err != nil {
		c.String(http.StatusNotFound, err.Error())
	}
}

func (p *PageServer) getPageHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	task, err := p.store.GetPage(id)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, task)
}

func (p *PageServer) getAllPagesHandler(c *gin.Context) {
	allTasks := p.store.GetAllPages()
	c.JSON(http.StatusOK, allTasks)
}

func (p *PageServer) deleteAllPagesHandler(*gin.Context) {
	_ = p.store.DeleteAllPages()
}
