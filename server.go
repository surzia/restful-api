package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"

	"restful-api/page"
)

type PageServer struct {
	store *page.Store
}

type PageRequest struct {
	Text string    `json:"text"`
	Tags []string  `json:"tags"`
	Due  time.Time `json:"due"`
}

type PageResponse struct {
	Id int `json:"id"`
}

func NewPageServer() *PageServer {
	return &PageServer{store: page.New()}
}

// renderJSON renders 'v' as JSON and writes it as a response into w.
func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(js)
}

func (p *PageServer) pageHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling page at %s\n", r.URL.Path)
	if r.URL.Path == "/page/" {
		// response for url "/page/"
		if r.Method == http.MethodPost {
			p.createPageHandler(w, r)
		} else if r.Method == http.MethodPut {
			p.updatePageHandler(w, r)
		} else if r.Method == http.MethodGet {
			p.getAllPagesHandler(w, r)
		} else if r.Method == http.MethodDelete {
			p.deleteAllPagesHandler(w, r)
		} else {
			http.Error(w, fmt.Sprintf("expect method GET, DELETE, PUT or POST at /page/, got %v", r.Method), http.StatusMethodNotAllowed)
			return
		}
	} else {
		// response for url "/page/<id>"
		path := strings.Trim(r.URL.Path, "/")
		pathParts := strings.Split(path, "/")

		// if the url is otherwise
		if len(pathParts) < 2 {
			http.Error(w, "expect /page/<id> in page handler", http.StatusBadRequest)
			return
		}

		// parse id from url
		id, err := strconv.Atoi(pathParts[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodDelete {
			p.deletePageHandler(w, r, id)
		} else if r.Method == http.MethodGet {
			p.getPageHandler(w, r, id)
		} else {
			http.Error(w, fmt.Sprintf("expect method GET or DELETE at /page/<id>, got %v", r.Method), http.StatusMethodNotAllowed)
			return
		}
	}
}

func (p *PageServer) tagHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling tags at %s\n", r.URL.Path)
	// this handler only accept http get
	if r.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("expect method GET /tag/<tag>, got %v", r.Method), http.StatusMethodNotAllowed)
		return
	}

	path := strings.Trim(r.URL.Path, "/")
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 2 {
		http.Error(w, "expect /tag/<tag> path", http.StatusBadRequest)
		return
	}

	tag := pathParts[1]
	pages := p.store.GetPagesByTag(tag)
	renderJSON(w, pages)
}

func (p *PageServer) dueHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling due at %s\n", r.URL.Path)
	// this handler only accept http get
	if r.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("expect method GET /due/<date>, got %v", r.Method), http.StatusMethodNotAllowed)
		return
	}

	path := strings.Trim(r.URL.Path, "/")
	pathParts := strings.Split(path, "/")

	// this handler only accept /due/<year>/<month>/<day>
	badRequestError := func() {
		http.Error(w, fmt.Sprintf("expect /due/<year>/<month>/<day>, got %v", r.URL.Path), http.StatusBadRequest)
	}
	if len(pathParts) != 4 {
		badRequestError()
		return
	}

	year, err := strconv.Atoi(pathParts[1])
	if err != nil {
		badRequestError()
		return
	}
	month, err := strconv.Atoi(pathParts[2])
	if err != nil || month < int(time.January) || month > int(time.December) {
		badRequestError()
		return
	}
	day, err := strconv.Atoi(pathParts[3])
	if err != nil {
		badRequestError()
		return
	}

	pages := p.store.GetPagesByDueDate(year, time.Month(month), day)
	renderJSON(w, pages)
}

func (p *PageServer) createPageHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling create page at %s\n", r.URL.Path)
	// Enforce a JSON Content-Type.
	contentType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediaType != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var pr PageRequest
	if err := dec.Decode(&pr); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := p.store.CreatePage(pr.Text, pr.Tags, pr.Due)
	renderJSON(w, PageResponse{Id: id})
}

func (p *PageServer) updatePageHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling update page at %s\n", r.URL.Path)
	// Enforce a JSON Content-Type.
	contentType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediaType != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var ret page.Page
	if err := dec.Decode(&ret); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	upd, err := p.store.UpdatePage(&ret)
	renderJSON(w, upd)
}

func (p *PageServer) deletePageHandler(w http.ResponseWriter, r *http.Request, id int) {
	log.Printf("handling delete page at %s\n", r.URL.Path)
	err := p.store.DeletePage(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte(fmt.Sprintf("page with id=%d has been deleted.", id)))
}

func (p *PageServer) getPageHandler(w http.ResponseWriter, r *http.Request, id int) {
	log.Printf("handling get page at %s\n", r.URL.Path)
	ret, err := p.store.GetPage(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, ret)
}

func (p *PageServer) getAllPagesHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling get all pages at %s\n", r.URL.Path)
	ret := p.store.GetAllPages()
	renderJSON(w, ret)
}

func (p *PageServer) deleteAllPagesHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling delete all pages at %s\n", r.URL.Path)
	_ = p.store.DeleteAllPages()

	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("All pages have been deleted."))
}
