package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	ml "github.com/IBM/mathlib"
	"github.com/hyperledger/aries-bbs-go/bbs"
)

func main() {
	pkBase64 := "lOpN7uGZWivVIjs0325N/V0dAhoPomrgfXVpg7pZNdRWwFwJDVxoE7TvRyOx/Qr7GMtShNuS2Px/oScD+SMf08t8eAO78QRNErPzwNpfkP4ppcSTShStFDfFbsv9L9yb"
	pkBytes, _ := base64.RawStdEncoding.DecodeString(pkBase64)

	sigBase64 := "hPbLkeMZZ6KKzkjWoTVHeMeuLJfYWjmdAU1Vg5fZ/VZnIXxxeXBB+q0/EL8XQmWkOMMwEGA/D2dCb4MDuntKZpvHEHlvaFR6l1A4bYj0t2Jd6bYwGwCwirNbmSeIoEmJeRzJ1cSvsL+jxvLixdDPnw=="
	sigBytes, _ := base64.StdEncoding.DecodeString(sigBase64)

	messagesBytes := [][]byte{
		[]byte("message1"),
		[]byte("message2"),
	}

	bls := bbs.New(ml.Curves[ml.BLS12_381_BBS])

	err := bls.Verify(messagesBytes, sigBytes, pkBytes)
	if err != nil {
		fmt.Println("Verify Error.")
	} else {
		fmt.Println("Verify OK.")
	}

	// swap messages order
	invalidMessagesBytes := [][]byte{[]byte("message2"), []byte("message1")}

	err = bls.Verify(invalidMessagesBytes, sigBytes, pkBytes)
	if err != nil {
		fmt.Println("Verify Error.", err) // invalid BLS12-381 signature
	}

	err = bls.Verify(messagesBytes, sigBytes, []byte("invalid"))
	if err != nil {
		fmt.Println("Verify Error.", err) // parse public key: invalid size of public key
	}

	pkBytesInvalid := make([]byte, len(pkBytes))
	_, _ = rand.Read(pkBytesInvalid)

	err = bls.Verify(messagesBytes, sigBytes, pkBytesInvalid)
	if err != nil {
		fmt.Println("Verify Error.", err) // parse public key: invalid size of public key
	}

	err = bls.Verify(messagesBytes, []byte("invalid"), pkBytes)
	if err != nil {
		fmt.Println("Verify Error.", err) // parse signature: invalid size of signature
	}

	sigBytesInvalid := make([]byte, len(sigBytes))

	_, _ = rand.Read(sigBytesInvalid)

	err = bls.Verify(messagesBytes, sigBytesInvalid, pkBytes)
	if err != nil {
		fmt.Println("Verify Error.", err) // parse signature: deserialize G1 compressed signature
	}
}
