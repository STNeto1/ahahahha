package main

import (
	"context"
	"log"
	"net"
	searchpb "search/gen/protos"

	"google.golang.org/grpc"
)

type searchService struct {
	searchpb.UnimplementedSearchServiceServer
}

func (*searchService) FetchCategories(_ context.Context, _ *searchpb.Empty) (*searchpb.FetchCategoriesResponse, error) {
	return &searchpb.FetchCategoriesResponse{
		Categories: []*searchpb.Category{
			{
				Id:       "someId",
				Name:     "Category 1",
				Slug:     "category-1",
				ParentId: nil,
			},
		},
	}, nil
}

func (*searchService) mustEmbedUnimplementedSearchServiceServer() {}

func main() {
	lst, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
	defer lst.Close()

	grpcServer := grpc.NewServer()
	searchpb.RegisterSearchServiceServer(grpcServer, &searchService{})

	log.Println("Server is running on port 1234")
	if err := grpcServer.Serve(lst); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
