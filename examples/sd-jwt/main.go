package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"github.com/golang-jwt/jwt"
)

type SDJWTClaims struct {
	jwt.StandardClaims
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"_sd:address"`
}

func generateSalt() string {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		log.Fatal(err)
	}
	return base64.RawURLEncoding.EncodeToString(salt)
}

func issueSDJWT(privateKey *rsa.PrivateKey) (string, string, error) {
	claims := SDJWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "test-issuer",
		},
		Name:    "John Doe",
		Age:     30,
		Address: generateSalt(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", "", err
	}

	disclosures := map[string]string{
		"address": "123 Main St",
	}

	disclosureJSON, err := json.Marshal(disclosures)
	if err != nil {
		return "", "", err
	}

	return signedToken, base64.RawURLEncoding.EncodeToString(disclosureJSON), nil
}

func verifySDJWT(tokenString string, publicKey *rsa.PublicKey, disclosures string) error {
	token, err := jwt.ParseWithClaims(tokenString, &SDJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(*SDJWTClaims); ok && token.Valid {
		fmt.Printf("Name: %v\n", claims.Name)
		fmt.Printf("Age: %v\n", claims.Age)

		var disclosureMap map[string]string
		disclosureBytes, _ := base64.RawURLEncoding.DecodeString(disclosures)
		json.Unmarshal(disclosureBytes, &disclosureMap)

		if address, ok := disclosureMap["address"]; ok {
			fmt.Printf("Address: %v\n", address)
		}
	} else {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func main() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}

	sdJWT, disclosures, err := issueSDJWT(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("SD-JWT: %v\n", sdJWT)
	fmt.Printf("Disclosures: %v\n", disclosures)

	err = verifySDJWT(sdJWT, &privateKey.PublicKey, disclosures)
	if err != nil {
		log.Fatal(err)
	}
}
