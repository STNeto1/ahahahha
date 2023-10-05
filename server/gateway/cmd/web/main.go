package main

import (
	"context"
	"errors"
	"gateway/pkg/core"
	"gateway/pkg/graph"
	"log"
	searchpb "search/gen/protos"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	db := core.InitDB()

	conn, err := grpc.Dial(":1234", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("failed to connect search grpc server", err)
	}
	defer conn.Close()

	e := echo.New()

	resolverHandler := handler.NewDefaultServer(createGraphqlSchema(db, conn))
	playgroundHandler := playground.Handler("GraphQL", "/query")

	// Middleware
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(core.UserMiddleware(db))
	e.Use(core.CookieMiddleWare())

	e.POST("/query", func(c echo.Context) error {
		req := c.Request()
		res := c.Response()
		resolverHandler.ServeHTTP(res, req)
		return nil
	})

	e.GET("/", func(c echo.Context) error {
		req := c.Request()
		res := c.Response()
		playgroundHandler.ServeHTTP(res, req)
		return nil
	})

	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}

func createGraphqlSchema(db *sqlx.DB, searchConn *grpc.ClientConn) graphql.ExecutableSchema {
	return graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		DB:           db,
		SearchClient: searchpb.NewSearchServiceClient(searchConn),
	},
		Directives: graph.DirectiveRoot{
			Auth: func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
				if usr := core.GetUserFromContext(ctx); usr == nil {
					return nil, errors.New("unauthorized")
				}

				return next(ctx)
			},
			Staff: func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
				if usr := core.GetUserFromContext(ctx); usr == nil || usr.Role != "staff" {
					return nil, errors.New("unauthorized")
				}

				return next(ctx)
			},
		},
	})
}
