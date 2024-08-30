package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	ml "github.com/IBM/mathlib"
	"github.com/hyperledger/aries-bbs-go/bbs"
)

func generateKeyPairRandom(curve *ml.Curve) (*bbs.PublicKey, *bbs.PrivateKey, error) {
	seed := make([]byte, 32)

	_, err := rand.Read(seed)
	if err != nil {
		panic(err)
	}

	bbs := bbs.NewBBSLib(curve)

	return bbs.GenerateKeyPair(sha256.New, seed)
}

var curve = ml.Curves[ml.BLS12_381_BBS]

func main() {
	pubKey, privKey, _ := generateKeyPairRandom(curve)

	bls := bbs.New(curve)

	messagesBytes := [][]byte{[]byte("message1"), []byte("message2")}

	signatureBytes, _ := bls.SignWithKey(messagesBytes, privKey)
	fmt.Println("signatureBytes: ", signatureBytes)

	pubKeyBytes, _ := pubKey.Marshal()
	fmt.Println("pubKeyBytes: ", pubKeyBytes)
}
