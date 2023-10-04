package graph

import (
	searchpb "search/gen/protos"

	"github.com/jmoiron/sqlx"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB           *sqlx.DB
	SearchClient searchpb.SearchServiceClient
}
