package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
)

func sign(msg string, pvKey *rsa.PrivateKey) ([]byte, error) {
	digest := sha256.Sum256([]byte(msg))

	signature, err := rsa.SignPKCS1v15(
		rand.Reader,
		pvKey,
		crypto.SHA256,
		digest[:],
	)

	return signature, err
}

func main() {
	pvKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic("Failed to Generate RSA.")
	}

	msg := "Hello World."
	sig, err := sign(msg, pvKey)
	if err != nil {
		panic("Failed to Sign.")
	}

	fmt.Printf("Signature: %v\n", sig)

	sigBase58 := base58.Encode(sig)
	fmt.Printf("Signature(base58): %v\n", sigBase58)
}
