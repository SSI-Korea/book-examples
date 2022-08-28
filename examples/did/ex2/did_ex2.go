package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"ssi-book/core"
)

func NewDID(method string, pbKey string) (string, error) {
	if method == "" || pbKey == "" {
		return "", errors.New("parameter is not valid")
	}

	digest := sha256.Sum256([]byte(pbKey))
	specificIdentifier := base58.Encode(digest[:])

	// DID:Method:specific
	did := fmt.Sprintf("did:%s:%s", method, specificIdentifier)

	return did, nil
}

func main() {
	var method = "ssikr"

	kms := new(core.ECDSAManager)
	kms.Generate()

	did, err := NewDID(method, kms.PublicKeyMultibase())

	if err != nil {
		panic(fmt.Sprintf("Failed to generate DID, error: %v\n", err))
	}

	fmt.Println("### New DID ###")
	fmt.Printf("did => %s\n", did)
}
