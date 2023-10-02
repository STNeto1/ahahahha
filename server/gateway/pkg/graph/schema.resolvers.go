package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.38

import (
	"context"
	"gateway/pkg/core"
	"gateway/pkg/graph/model"

	"github.com/vektah/gqlparser/v2/gqlerror"
)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (bool, error) {
	err := core.CreateUser(r.DB, input.Name, input.Email, input.Password)
	if err != nil {
		return false, gqlerror.Errorf(err.Error())
	}

	return true, nil
}

// AuthenticateUser is the resolver for the authenticateUser field.
func (r *mutationResolver) AuthenticateUser(ctx context.Context, input model.AuthenticatedUserInput) (string, error) {
	usr, err := core.AuthenticateUser(r.DB, input.Email, input.Password)
	if err != nil {
		return "", gqlerror.Errorf(err.Error())
	}

	return usr.Name, nil
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context, term *string) ([]*model.User, error) {
	users, err := core.SearchUsers(r.DB, term)
	if err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	return core.MapUsers(users), nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	usr, err := core.GetUser(r.DB, id)
	if err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	return core.MapUser(usr), nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
