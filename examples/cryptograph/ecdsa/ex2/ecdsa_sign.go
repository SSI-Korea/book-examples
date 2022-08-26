package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"
)

type Signature struct {
	R *big.Int
	S *big.Int
}

func (s Signature) String() string {
	return s.R.String() + s.S.String()
}

func sign(digest []byte, pvKey *ecdsa.PrivateKey) (*Signature, error) {
	r := big.NewInt(0)
	s := big.NewInt(0)

	r, s, err := ecdsa.Sign(rand.Reader, pvKey, digest)
	if err != nil {
		return nil, err //errors.New("failed to sign to msg.")
	}

	// prepare a signature structure to marshal into json
	signature := &Signature{
		R: r,
		S: s,
	}
	/*
		signature := r.Bytes()
		signature = append(signature, s.Bytes()...)
	*/
	return signature, nil
}

func SignASN1(digest []byte, pvKey *ecdsa.PrivateKey) ([]byte, error) {
	signature, err := ecdsa.SignASN1(rand.Reader, pvKey, digest)
	if err != nil {
		return nil, err //errors.New("failed to sign to msg.")
	}

	return signature, nil
}

func SignToString(digest []byte, pvKey *ecdsa.PrivateKey) (string, error) {
	signature, err := sign(digest, pvKey)
	if err != nil {
		return "", err
	}

	return signature.String(), nil
}

func main() {
	pvKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader) // elliptic.p224, elliptic.P384(), elliptic.P521()

	if err != nil {
		log.Println("ECDSA Keypair generation was Fail!")
	}

	//pbKey := &pvKey.PublicKey

	msg := "Hello SSI-KOREA."
	digest := sha256.Sum256([]byte(msg))
	signature, err := sign(digest[:], pvKey)
	if err != nil {
		log.Printf("Failed to sign msg.")
	}

	signatureASN1, err := SignASN1(digest[:], pvKey)
	if err != nil {
		log.Printf("Failed to sign msg.")
	}

	fmt.Printf("########## Sign ##########\n")
	fmt.Printf("===== Message =====\n")
	fmt.Printf("Msg: %s\n", msg)
	fmt.Printf("Digest: %x\n", digest)
	fmt.Printf("R: %s, S: %s\n", signature.R, signature.S)
	fmt.Printf("Signature: %+v\n", signature.String())
	fmt.Printf("SignatureASN1: %+v\n", signatureASN1)

}
