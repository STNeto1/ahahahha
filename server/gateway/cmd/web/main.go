package main

import (
	"gateway/pkg/core"
	"gateway/pkg/graph"
	"log"
	searchpb "search/gen/protos"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial(":1234", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("failed to connect search grpc server", err)
	}
	defer conn.Close()

	e := echo.New()

	resolverHandler := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		DB:           core.InitDB(),
		SearchClient: searchpb.NewSearchServiceClient(conn),
	}}))
	playgroundHandler := playground.Handler("GraphQL", "/query")

	// Middleware
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(core.UserMiddleware())
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
