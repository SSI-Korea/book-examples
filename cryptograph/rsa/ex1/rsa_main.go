package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"log"
	"os"
)

func main() {
	// Generate RSA Key
	pvKey, err := rsa.GenerateKey(rand.Reader, 1024) // 2048, 4096
	if err != nil {
		log.Printf("Failed to Generate RSA.")
		os.Exit(0)
	}

	fmt.Printf("Private Key: %+v\n", pvKey) //%+v 필드명과 함께 출력.

	pvKeyByte := x509.MarshalPKCS1PrivateKey(pvKey)
	pvKeyStr := base64.StdEncoding.EncodeToString(pvKeyByte)
	fmt.Printf("Private Key(str): %+v\n", pvKeyStr)

	pbKey := pvKey.PublicKey
	fmt.Printf("Public Key: %+v\n", pbKey)

	pbKeyByte := x509.MarshalPKCS1PublicKey(&pbKey)
	fmt.Printf("Public Key(Byte): %+v\n", pbKeyByte)

	pbKeyStr := base64.StdEncoding.EncodeToString(pbKeyByte)
	fmt.Printf("Public Key(str): %+v\n", pbKeyStr)

}
