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

	privKeyBytes, _ := privKey.Marshal()

	signatureBytes, err := bls.Sign(messagesBytes, privKeyBytes)
	if err != nil {
		fmt.Println("Sign error. ", err)
	} else {
		fmt.Println("Sign OK.")
	}

	// require.NoError(t, err)
	// require.NotEmpty(t, signatureBytes)
	// require.Len(t, signatureBytes, curve.CompressedG1ByteSize+2*32)

	pubKeyBytes, _ := pubKey.Marshal()

	err = bls.Verify(messagesBytes, signatureBytes, pubKeyBytes)
	if err != nil {
		fmt.Println("Verify error. ", err)
	} else {
		fmt.Println("Verify OK.")
	}

	// // invalid private key bytes
	// signatureBytes, err = bls.Sign(messagesBytes, []byte("invalid"))
	// require.Error(t, err)
	// require.EqualError(t, err, "unmarshal private key: invalid size of private key")
	// require.Nil(t, signatureBytes)

	// // at least one message must be passed
	// signatureBytes, err = bls.Sign([][]byte{}, privKeyBytes)
	// require.Error(t, err)
	// require.EqualError(t, err, "messages are not defined")
	// require.Nil(t, signatureBytes)
}
