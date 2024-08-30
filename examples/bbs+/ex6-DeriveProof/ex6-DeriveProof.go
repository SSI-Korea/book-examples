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

	privKeyBytes, _ := privKey.Marshal()

	messagesBytes := [][]byte{
		[]byte("message1"),
		[]byte("message2"),
		[]byte("message3"),
		[]byte("message4"),
	}
	bls := bbs.New(curve)

	signatureBytes, err := bls.Sign(messagesBytes, privKeyBytes)
	if err != nil {
		fmt.Println("Sign Error. ", err)
	} else {
		fmt.Println("Sign OK.")
	}

	pubKeyBytes, _ := pubKey.Marshal()

	err = bls.Verify(messagesBytes, signatureBytes, pubKeyBytes)
	if err != nil {
		fmt.Println("Verify Error. ", err)
	} else {
		fmt.Println("Verify OK.")
	}

	nonce := []byte("nonce")
	revealedIndexes := []int{0, 2}
	proofBytes, err := bls.DeriveProof(messagesBytes, signatureBytes, nonce, pubKeyBytes, revealedIndexes)
	if err != nil {
		fmt.Println("DeriveProof Error. ", err)
	} else {
		fmt.Println("DeriveProof OK.")
	}

	/** ##########	*/
	revealedMessages := make([][]byte, len(revealedIndexes))
	for i, ind := range revealedIndexes {
		revealedMessages[i] = messagesBytes[ind]
	}

	err = bls.VerifyProof(revealedMessages, proofBytes, nonce, pubKeyBytes)
	if err != nil {
		fmt.Println("VerifyProof Error. ", err)
	} else {
		fmt.Println("VerifyProof OK.")
	}

}
