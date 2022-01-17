package page

import (
	"fmt"
	"sync"
	"time"
)

type Page struct {
	Id   int       `json:"id"`
	Text string    `json:"text"`
	Tags []string  `json:"tags"`
	Due  time.Time `json:"due"`
}

// Store is a simple in-memory database of pages; Store methods are
// safe to call concurrently.
type Store struct {
	sync.Mutex

	pages  map[int]*Page
	nextId int
}

func New() *Store {
	return &Store{
		pages:  make(map[int]*Page),
		nextId: 0,
	}
}

// CreatePage creates a new page in the store.
func (s *Store) CreatePage(text string, tags []string, due time.Time) int {
	s.Lock()
	defer s.Unlock()

	page := &Page{
		Id:   s.nextId,
		Text: text,
		Due:  due,
	}

	page.Tags = make([]string, len(tags))
	copy(page.Tags, tags)

	s.pages[s.nextId] = page
	s.nextId++

	return page.Id
}

// DeletePage deletes the page with the given id. If no such id exists, an error
// is returned.
func (s *Store) DeletePage(id int) error {
	s.Lock()
	defer s.Unlock()

	_, ok := s.pages[id]
	if ok {
		delete(s.pages, id)
		return nil
	}

	return fmt.Errorf("page with id=%d not found", id)
}

// UpdatePage updates the page with the given id and new page. If no such id
// exists, an error is returned, else page with given id should be updated.
func (s *Store) UpdatePage(page *Page) (*Page, error) {
	s.Lock()
	defer s.Unlock()

	_, ok := s.pages[page.Id]
	s.pages[page.Id] = page

	if ok {
		return page, nil
	}

	return nil, fmt.Errorf("page with id=%d not found", page.Id)
}

// GetPage retrieves a page from the store, by id. If no such id exists, an
// error is returned.
func (s *Store) GetPage(id int) (*Page, error) {
	s.Lock()
	defer s.Unlock()

	t, ok := s.pages[id]
	if ok {
		return t, nil
	}

	return nil, fmt.Errorf("page with id=%d not found", id)
}

// GetAllPages returns all the pages in the store, in arbitrary order.
func (s *Store) GetAllPages() []*Page {
	s.Lock()
	defer s.Unlock()

	ret := make([]*Page, 0, len(s.pages))
	for _, page := range s.pages {
		ret = append(ret, page)
	}
	return ret
}

// DeleteAllPages deletes all pages in the store.
func (s *Store) DeleteAllPages() error {
	s.Lock()
	defer s.Unlock()

	s.pages = make(map[int]*Page)
	return nil
}

// GetPagesByTag returns all the pages that have the given tag, in arbitrary
// order.
func (s *Store) GetPagesByTag(tag string) []*Page {
	s.Lock()
	defer s.Unlock()

	var ret []*Page

search:
	for _, page := range s.pages {
		for _, pageTag := range page.Tags {
			if pageTag == tag {
				ret = append(ret, page)
				continue search
			}
		}
	}

	return ret
}

// GetPagesByDueDate returns all the pages that have the given due date, in
// arbitrary order.
func (s *Store) GetPagesByDueDate(year int, month time.Month, day int) []*Page {
	s.Lock()
	defer s.Unlock()

	var ret []*Page

	for _, page := range s.pages {
		y, m, d := page.Due.Date()
		if y == year && m == month && d == day {
			ret = append(ret, page)
		}
	}

	return ret
}
