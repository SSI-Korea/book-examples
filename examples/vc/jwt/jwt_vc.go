package main

import (
	"crypto/ecdsa"
	"crypto/x509"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/multiformats/go-multibase"
	"log"
	"os"
	"ssi-book/core"
	"strings"
)

// Issuer에 의한 VC 발행 예시.
func main() {
	// 키생성(ECDSA) - 향후 KMS로 대체.
	issuerKeyEcdsa := core.NewEcdsa()

	// DID 생성.
	issuerDid, _ := core.NewDID("ssikr", issuerKeyEcdsa.PublicKeyBase58())

	// DID Document 생성.
	verificationId := fmt.Sprintf("%s#keys-1", issuerDid)
	verificationMethod := []core.VerificationMethod{
		{
			Id:                 verificationId,
			Type:               "EcdsaSecp256k1VerificationKey2019",
			Controller:         issuerDid.String(),
			PublicKeyMultibase: issuerKeyEcdsa.PublicKeyMultibase(),
		},
	}
	didDocument := core.NewDIDDocument(issuerDid.String(), verificationMethod)
	core.RegisterDid(issuerDid.String(), didDocument.String())

	fmt.Println("DID Document: ", didDocument)

	// VC 생성.
	vc, err := core.NewVC(
		"1234567890",
		[]string{"VerifiableCredential", "AlumniCredential"},
		issuerDid.String(),
		map[string]interface{}{
			"id": "abcd1234567890",
			"alumniOf": map[string]interface{}{
				"id": "1234567",
				"name": []map[string]string{
					{
						"value": "Example University",
						"lang":  "en",
					}, {
						"value": "Exemple d'Université",
						"lang":  "fr",
					},
				},
			},
		},
	)

	if err != nil {
		fmt.Println("Failed creation VC.")
		os.Exit(0)
	}

	// VC에 Issuer의 private key로 서명한다.(JWT 사용)
	tokenString, err := vc.GenerateJWT(verificationId, issuerKeyEcdsa.PrivateKey)

	fmt.Println("JWT Token: ", tokenString)

	// 생성된 VC를 검증한다.(public key를 직접 사용해서 검증)
	parts := strings.Split(tokenString, ".")
	err = jwt.SigningMethodES256.Verify(strings.Join(parts[0:2], "."), parts[2], issuerKeyEcdsa.PublicKey)

	if err != nil {
		panic("VC is Not verified.")
	}

	token, err := jwt.ParseWithClaims(tokenString, &core.JwtClaims{}, nil)
	if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
		fmt.Println("unexpected signing method: %v", token.Header["alg"])
	}

	fmt.Println("Token Header: ", token.Header)

	claims, ok := token.Claims.(*core.JwtClaims)

	if !ok {
		panic("claim has error.")
	}

	fmt.Println("Issuer DID: ", claims.Issuer)

	didDocumentStr, err := core.ResolveDid(claims.Issuer)
	if err != nil {
		log.Printf("Failed to Resolve DID.\nError: %x\n", err)
	}

	if didDocumentStr == "" {
		log.Printf("DID Document not found in VDR.")
	}

	fmt.Println("DID Document: ", didDocumentStr)

	//Json string을 DID Document 객체로 생성한다.
	didDocument, err = core.NewDIDDocumentForString(didDocumentStr)

	if err != nil {
		log.Printf("Failed generate DID Document from string.\nError: %x\n", err)
	}

	// 첫 번째를 사용한다고 가정한다.
	// TODO: 키 ID(위의 kid)에 해당하는 키 값 구하기.
	pbKeyBaseMultibase := didDocument.VerificationMethod[0].PublicKeyMultibase
	_, bytePubKey, err := multibase.Decode(pbKeyBaseMultibase)
	pbKey, err := x509.ParsePKIXPublicKey(bytePubKey)

	isVerify, err := core.VerifyJwt(tokenString, pbKey.(*ecdsa.PublicKey))
	if isVerify {
		fmt.Println("Verified.")
	} else {
		fmt.Println("Not Verified.", err)
	}

	// ParseAndVerifyJwtForVC() 함수를 이용한 verify
	isVerify, claims, err = core.ParseAndVerifyJwtForVC(tokenString)
	if isVerify {
		fmt.Println("Verified2.")
		fmt.Println("claims: ", claims)
	} else {
		fmt.Println("Not Verified.", err)
	}

}
