package main

import (
	"context"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"google.golang.org/grpc"
	"log"
	"net"
	"ssi-book/protos"
)

type resolverServer struct {
	protos.UnimplementedResolverServer
}

func byte2string(b []byte) string {
	return string(b[:len(b)])
}

func (server *resolverServer) ResolveDid(ctx context.Context, req *protos.ResolverRequest) (*protos.ResolverResponse, error) {
	log.Printf("Resolve DID: %s\n", req.Did)

	db, err := leveldb.OpenFile("../did_db/dids", nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	data, err := db.Get([]byte(req.Did), nil)
	didDocument := byte2string(data)
	return &protos.ResolverResponse{DidDocument: didDocument}, nil
}

func main() {
	fmt.Println("### Start Resolver ###")
	lis, err := net.Listen("tcp", "localhost:9001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := resolverServer{}
	s := grpc.NewServer()
	protos.RegisterResolverServer(s, &server)

	log.Printf("Resolver Server is listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
