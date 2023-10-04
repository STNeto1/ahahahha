package main

import (
	"log"
	"net"
	searchpb "search/gen/protos"
	"search/pkg/core"

	"google.golang.org/grpc"
)

func main() {
	lst, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
	defer lst.Close()

	grpcServer := grpc.NewServer()
	searchpb.RegisterSearchServiceServer(grpcServer, core.CreateSearchService(core.InitDB()))

	log.Println("Server is running on port 1234")
	log.Fatalln(grpcServer.Serve(lst))
}
