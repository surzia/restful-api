package page

import (
	"fmt"
	"restful-api/graph/model"
	"sync"
	"time"
)

// Store is a simple in-memory database of pages; Store methods are
// safe to call concurrently.
type Store struct {
	sync.Mutex

	pages  map[int]*model.Page
	nextId int
}

func New() *Store {
	return &Store{
		pages:  make(map[int]*model.Page),
		nextId: 0,
	}
}

// CreatePage creates a new page in the store.
func (s *Store) CreatePage(text string, tags []string, due time.Time, attachments []*model.Attachment) int {
	s.Lock()
	defer s.Unlock()

	page := &model.Page{
		ID:          s.nextId,
		Text:        text,
		Due:         due,
		Attachments: attachments,
	}

	page.Tags = make([]string, len(tags))
	copy(page.Tags, tags)

	s.pages[s.nextId] = page
	s.nextId++

	return page.ID
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
func (s *Store) UpdatePage(page *model.Page) (*model.Page, error) {
	s.Lock()
	defer s.Unlock()

	_, ok := s.pages[page.ID]
	s.pages[page.ID] = page

	if ok {
		return page, nil
	}

	return nil, fmt.Errorf("page with id=%d not found", page.ID)
}

// GetPage retrieves a page from the store, by id. If no such id exists, an
// error is returned.
func (s *Store) GetPage(id int) (*model.Page, error) {
	s.Lock()
	defer s.Unlock()

	t, ok := s.pages[id]
	if ok {
		return t, nil
	}

	return nil, fmt.Errorf("page with id=%d not found", id)
}

// GetAllPages returns all the pages in the store, in arbitrary order.
func (s *Store) GetAllPages() []*model.Page {
	s.Lock()
	defer s.Unlock()

	ret := make([]*model.Page, 0, len(s.pages))
	for _, page := range s.pages {
		ret = append(ret, page)
	}
	return ret
}

// DeleteAllPages deletes all pages in the store.
func (s *Store) DeleteAllPages() error {
	s.Lock()
	defer s.Unlock()

	s.pages = make(map[int]*model.Page)
	return nil
}

// GetPagesByTag returns all the pages that have the given tag, in arbitrary
// order.
func (s *Store) GetPagesByTag(tag string) []*model.Page {
	s.Lock()
	defer s.Unlock()

	var ret []*model.Page

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
func (s *Store) GetPagesByDueDate(year int, month time.Month, day int) []*model.Page {
	s.Lock()
	defer s.Unlock()

	var ret []*model.Page

	for _, page := range s.pages {
		y, m, d := page.Due.Date()
		if y == year && m == month && d == day {
			ret = append(ret, page)
		}
	}

	return ret
}
