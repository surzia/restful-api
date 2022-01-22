package graph

import "restful-api/page"

type Resolver struct {
	Store *page.Store
}

func NewResolver() *Resolver {
	return &Resolver{Store: page.New()}
}
