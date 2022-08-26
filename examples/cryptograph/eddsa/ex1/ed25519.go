package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

func main() {
	pubKey, pvKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Println("Fail to generate key.")
	}

	fmt.Println("Private Key: ", pvKey)
	fmt.Println("Public Key: ", pubKey)

	msg := "Hello SSI-KOREA."
	digest := sha256.Sum256([]byte(msg))

	signature := ed25519.Sign(pvKey, digest[:])

	isVerify := ed25519.Verify(pubKey, digest[:], signature)
	
	if isVerify {
		fmt.Println("Verified.")
	} else {
		fmt.Println("Not Verified.")
	}
}
