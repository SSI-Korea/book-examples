package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"google.golang.org/grpc"
	"log"
	"math/big"
	"net"
	"ssi-book/protos"
	didregistry "ssi-book/vdr/ether"
)

const (
	CONTRACT_ADDRESS    = "0x219E540089D05826c1a422ebfbB8a1C3348886A7"
	ACCOUNT_PRIVATE_KEY = "d31fc68bc7e2d3296e462a9510a6549ac32a23627ef2781e25d76c35a130c82a"
)

func getAccountAuth(client *ethclient.Client, accountAddress string) *bind.TransactOpts {

	privateKey, err := crypto.HexToECDSA(accountAddress)
	if err != nil {
		panic(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("invalid key")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	//fetch the last use nonce of account
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		panic(err)
	}
	//fmt.Println("nounce=", nonce)
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		panic(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		panic(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = big.NewInt(1000000)

	return auth
}

type resolverServer struct {
	protos.UnimplementedResolverServer
}

func byte2string(b []byte) string {
	return string(b[:len(b)])
}

func (server *resolverServer) ResolveDid(ctx context.Context, req *protos.ResolverRequest) (*protos.ResolverResponse, error) {
	log.Printf("Resolve DID: %s\n", req.Did)

	client, err := ethclient.Dial("http://127.0.0.1:8545")
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
