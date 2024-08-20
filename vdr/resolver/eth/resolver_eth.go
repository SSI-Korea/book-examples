package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"ssi-book/protos"
	didregistry "ssi-book/vdr/ether"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"google.golang.org/grpc"
)

const (
	CONTRACT_ADDRESS = "0x5FbDB2315678afecb367f032d93F642f64180aa3"
	RPC_ENDPOINT     = "http://127.0.0.1:8545"
)

type resolverServer struct {
	protos.UnimplementedResolverServer
}

func byte2string(b []byte) string {
	return string(b[:len(b)])
}

func (server *resolverServer) ResolveDid(ctx context.Context, req *protos.ResolverRequest) (*protos.ResolverResponse, error) {
	log.Printf("Resolve DID: %s\n", req.Did)

	client, err := ethclient.Dial(RPC_ENDPOINT)
	if err != nil {
		panic(err)
	}

	conn, err := didregistry.NewDidregistry(common.HexToAddress(CONTRACT_ADDRESS), client)
	if err != nil {
		panic(err)
	}

	didDocument, err := conn.ResolveDid(&bind.CallOpts{}, req.Did)
	if err != nil {
		panic(err)
	}

	return &protos.ResolverResponse{DidDocument: didDocument}, nil
}

func main() {
	fmt.Println("### Start Resolver(eth) ###")
	lis, err := net.Listen("tcp", "localhost:9901")
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
