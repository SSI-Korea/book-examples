package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
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

func main() {
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		panic(err)
	}

	conn, err := didregistry.NewDidregistry(common.HexToAddress(CONTRACT_ADDRESS), client)
	if err != nil {
		panic(err)
	}

	txOpts := getAccountAuth(client, ACCOUNT_PRIVATE_KEY)

	result, err := conn.CreateDid(txOpts, "did:ssikr:1234asdf1234", "{new DID document}")
	if err != nil {
		panic(err)
	}

	_ = result

	didDocument, err := conn.ResolveDid(&bind.CallOpts{}, "did:ssikr:1234asdf1234")
	if err != nil {
		panic(err)
	}

	fmt.Println(didDocument)

}
