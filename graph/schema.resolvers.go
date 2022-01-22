package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"restful-api/graph/generated"
	"restful-api/graph/model"
	"time"
)

func (r *mutationResolver) CreatePage(_ context.Context, input model.NewPage) (*model.Page, error) {
	attachments := make([]*model.Attachment, 0, len(input.Attachments))
	for _, a := range input.Attachments {
		attachments = append(attachments, (*model.Attachment)(a))
	}
	id := r.Store.CreatePage(input.Text, input.Tags, input.Due, attachments)
	page, err := r.Store.GetPage(id)
	return page, err
}

func (r *mutationResolver) DeletePage(_ context.Context, id int) (*bool, error) {
	return nil, r.Store.DeletePage(id)
}

func (r *mutationResolver) DeleteAllPages(context.Context) (*bool, error) {
	return nil, r.Store.DeleteAllPages()
}

func (r *queryResolver) GetAllPages(context.Context) ([]*model.Page, error) {
	return r.Store.GetAllPages(), nil
}

func (r *queryResolver) GetPage(_ context.Context, id int) (*model.Page, error) {
	return r.Store.GetPage(id)
}

func (r *queryResolver) GetPagesByTag(_ context.Context, tag string) ([]*model.Page, error) {
	return r.Store.GetPagesByTag(tag), nil
}

func (r *queryResolver) GetPagesByDue(_ context.Context, due time.Time) ([]*model.Page, error) {
	y, m, d := due.Date()
	return r.Store.GetPagesByDueDate(y, m, d), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
