package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"net"
	"ssi-book/protos"
	didregistry "ssi-book/vdr/ether"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"google.golang.org/grpc"
)

const (
	CONTRACT_ADDRESS    = "0xa80525f1811e1809546413b29b8731d8f71e72bf"
	ACCOUNT_PRIVATE_KEY = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	RPC_ENDPOINT        = "http://127.0.0.1:8545"
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
	fmt.Println("nounce=", nonce)

	// chainID, err := client.ChainID(context.Background())
	// if err != nil {
	// 	panic(err)
	// }

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth := bind.NewKeyedTransactor(privateKey)
	if err != nil {
		panic(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = gasPrice        //big.NewInt(16231732509)

	return auth
}

type registrarServer struct {
	protos.UnimplementedRegistrarServer
}

func (server *registrarServer) RegisterDid(ctx context.Context, req *protos.RegistrarRequest) (*protos.RegistrarResponse, error) {
	log.Printf("Register DID: %s\n", req.Did)

	client, err := ethclient.Dial(RPC_ENDPOINT)
	// client, err := ethclient.Dial("")
	if err != nil {
		panic(err)
	}

	instance, err := didregistry.NewDidregistry(common.HexToAddress(CONTRACT_ADDRESS), client)
	if err != nil {
		panic(err)
	}

	auth := getAccountAuth(client, ACCOUNT_PRIVATE_KEY)

	tx, err := instance.RegisterDid(auth, req.Did, req.DidDocument)
	if err != nil {
		panic(err)
	}

	fmt.Printf("tx sent: %s\n", tx.Hash().Hex())

	_ = tx

	return &protos.RegistrarResponse{Result: tx.Hash().Hex()}, nil
}

func main() {
	fmt.Println("### Start Registrar(eth) ###")
	lis, err := net.Listen("tcp", "localhost:9900")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := registrarServer{}
	s := grpc.NewServer()
	protos.RegisterRegistrarServer(s, &server)

	log.Printf("Registrar Server is listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
