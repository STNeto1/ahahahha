package loaders

import (
	"context"
	"models"
	"strings"
	"time"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type ctxKey struct {
	name string
}

var loadersKey = &ctxKey{"dataloaders"}

// handleError creates array of result with the same error repeated for as many items requested
func handleError[T any](itemsLength int, err error) []*dataloader.Result[T] {
	result := make([]*dataloader.Result[T], itemsLength)
	for i := 0; i < itemsLength; i++ {
		result[i] = &dataloader.Result[T]{Error: err}
	}
	return result
}

// Loaders wrap your data loaders to inject via middleware
type Loaders struct {
	UserLoader *dataloader.Loader[string, *models.User]
}

// NewLoaders instantiates data loaders for the middleware
func NewLoaders(conn *sqlx.DB) *Loaders {
	// define the data loader
	ur := &userReader{db: conn}
	return &Loaders{
		UserLoader: dataloader.NewBatchedLoader(ur.getUsers, dataloader.WithWait[string, *models.User](time.Millisecond)),
	}
}

// userReader reads Users from a database
type userReader struct {
	db *sqlx.DB
}

// getUsers implements a batch function that can retrieve many users by ID,
// for use in a dataloader
func (u *userReader) getUsers(ctx context.Context, userIds []string) []*dataloader.Result[*models.User] {
	stmt, err := u.db.PrepareContext(ctx, `SELECT id, name FROM users WHERE id IN (?`+strings.Repeat(",?", len(userIds)-1)+`)`)
	if err != nil {
		return handleError[*models.User](len(userIds), err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, userIds)
	if err != nil {
		return handleError[*models.User](len(userIds), err)
	}
	defer rows.Close()

	result := make([]*dataloader.Result[*models.User], 0, len(userIds))
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			result = append(result, &dataloader.Result[*models.User]{Error: err})
			continue
		}
		result = append(result, &dataloader.Result[*models.User]{Data: &user})
	}
	return result
}

func LoaderMiddleware(db *sqlx.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			loader := NewLoaders(db)

			ctx := context.WithValue(c.Request().Context(), loadersKey, loader)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

// For returns the dataloader for a given context
func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}

// GetUser returns single user by id efficiently
func GetUser(ctx context.Context, userID string) (*models.User, error) {
	loaders := For(ctx)
	return loaders.UserLoader.Load(ctx, userID)()
}

// GetUsers returns many users by ids efficiently
func GetUsers(ctx context.Context, userIDs []string) ([]*models.User, []error) {
	loaders := For(ctx)
	return loaders.UserLoader.LoadMany(ctx, userIDs)()
}
