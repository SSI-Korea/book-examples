package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"github.com/multiformats/go-multibase"
	"math/big"
	"ssi-book/core"
)

// Actor #1
// 1. 키쌍 생성
// 2. DID 생성
// 3. DID Document 생성
// 4. DID Document 등록
// 5. 접속 및 DID 전송(생략)
//
// Actor #2
// 6. DID Document 리졸브
// 7. VerificationMethod(또는 Authentication)의 키 선택
// 8. 키 ID와 랜덤 메시지 전달(challenge)
//
// Actor #1
// 9. 키 ID에 해당하는 개인키로 메시지에 서명 후 전달(response)
//
// Actor #2
// 10. 검증
func main() {
	// 1. 키쌍 생성
	keyEcdsa := core.NewEcdsa()

	// 2. DID 생성.
	did, _ := core.NewDID("ssikr", keyEcdsa.PublicKeyBase58())

	// 3. DID Document 생성.
	verificationId := fmt.Sprintf("%s#keys-1", did)
	verificationMethod := []core.VerificationMethod{
		{
			Id:                 verificationId,
			Type:               "EcdsaSecp256k1VerificationKey2019",
			Controller:         did.String(),
			PublicKeyMultibase: keyEcdsa.PublicKeyMultibase(),
		},
	}
	didDocument := core.NewDIDDocument(did.String(), verificationMethod)
	// 4. DID Document 등록
	core.RegisterDid(did.String(), didDocument.String())

	// 5. 전송했다고 가정

	// 6. DID Document 리졸브
	didDocumentStr, err := core.ResolveDid(did.String())

	if err != nil {
		fmt.Errorf("Failed to Resolve DID.\nError: %x\n", err)
	}

	if didDocumentStr == "" {
		fmt.Errorf("DID Document not found in VDR.")
	}

	didDoc, err := core.NewDIDDocumentForString(didDocumentStr)
	if err != nil {
		fmt.Errorf("Failed generate DID Document from string.\nError: %x\n", err)
	}

	// 7. VerificationMethod(또는 Authentication)의 키 선택 - 첫 번째를 사용한다고 가정한다.
	challengeKeyId := didDoc.VerificationMethod[0].Id
	n, err := rand.Int(rand.Reader, big.NewInt(10000))
	challengeMsg := fmt.Sprintf("ssikr-challenge-msg-%d", n)

	// 8. 키 ID와 랜덤 메시지 전달(challenge) - 가정

	// 9. 키 ID에 해당하는 개인키로 메시지에 서명 후 전달(response)
	_ = challengeKeyId
	digest := sha256.Sum256([]byte(challengeMsg))
	signature, err := ecdsa.SignASN1(rand.Reader, keyEcdsa.PrivateKey, digest[:])

	// 10. 검증
	pbKeyBaseMultibase := didDoc.VerificationMethod[0].PublicKeyMultibase
	_, bytePubKey, err := multibase.Decode(pbKeyBaseMultibase)
	pbKey, err := x509.ParsePKIXPublicKey(bytePubKey)

	isVerify := ecdsa.VerifyASN1(pbKey.(*ecdsa.PublicKey), digest[:], signature)
	if isVerify {
		fmt.Println("### Challenge Success ###")
	} else {
		fmt.Println("### Challenge Fail ###")
	}

}
