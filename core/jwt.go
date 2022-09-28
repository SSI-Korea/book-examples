package core

import (
	"crypto/ecdsa"
	"crypto/x509"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/multiformats/go-multibase"
	"strings"
)

func VerifyJwt(token string, pbKey *ecdsa.PublicKey) (bool, error) {
	parts := strings.Split(token, ".")
	err := jwt.SigningMethodES256.Verify(strings.Join(parts[0:2], "."), parts[2], pbKey)
	if err != nil {
		return false, nil
	}
	return true, nil
}

// Parse VC JWT Claim and Verify VC JWT.
// 1. 토큰 파싱
// 2. claims에서 발급자의 DID 추출
// 3. 발급자의 DID로 DID Document 리졸브
// 4. DID Document에서 공개키 추출
// 5. 공개키로 JWT 검증
//
// DID도큐먼트의 key ID를 기준으로 public key의 값을 가져와야 하나,
// 여기서는 1개만 존재한다고 가정하고 첫번째를 사용해서 public key를 만들어 사용한다.
func ParseAndVerifyJwtForVC(tokenString string) (bool, *JwtClaims, error) {
	// 1. tokenString을 파싱하고, token을 검증한다.
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// 2. 발급자의 DID 추출
		claims, ok := token.Claims.(*JwtClaims)
		if !ok {
			return nil, fmt.Errorf("claims is wrong.")
		}
		fmt.Println("Issuer DID: ", claims.Issuer)

		// 3. Resolve한다.
		didDocumentStr, err := ResolveDid(claims.Issuer)
		if err != nil {
			return nil, fmt.Errorf("Failed to Resolve DID.\nError: %x\n", err)
		}

		if didDocumentStr == "" {
			return nil, fmt.Errorf("DID Document not found in VDR.")
		}

		fmt.Println("DID Document: ", didDocumentStr)

		// 4. DID Document에서 공개키 추출.
		//Json string을 DID Document 객체로 생성한다.
		didDocument, err := NewDIDDocumentForString(didDocumentStr)
		if err != nil {
			return nil, fmt.Errorf("Failed generate DID Document from string.\nError: %x\n", err)
		}

		// 첫 번째를 사용한다고 가정한다.
		// TODO: 키 ID(위의 kid)에 해당하는 키 값 구하기.
		pbKeyBaseMultibase := didDocument.VerificationMethod[0].PublicKeyMultibase
		_, bytePubKey, err := multibase.Decode(pbKeyBaseMultibase)
		pbKey, err := x509.ParsePKIXPublicKey(bytePubKey)

		return pbKey, nil //token이 검증되면 token.Valid는 true가 된다.
	})

	// 5. JWT 검증
	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		//fmt.Printf("%v %v", claims.Vc, claims.Issuer)
		return true, claims, nil
	}

	return false, nil, err
}

func ParseAndVerifyJwtForVP(tokenString string) (bool, *JwtClaimsForVP, error) {
	// 1. tokenString을 파싱하고, token을 검증한다.
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaimsForVP{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// 2. 보유자의 DID 추출: JWT의 입장에서는 issuer이다.
		claims, ok := token.Claims.(*JwtClaimsForVP)
		if !ok {
			return nil, fmt.Errorf("claims is wrong.")
		}
		fmt.Println("Issuer DID: ", claims.Issuer)

		// 3. Resolve한다.
		didDocumentStr, err := ResolveDid(claims.Issuer)
		if err != nil {
			return nil, fmt.Errorf("Failed to Resolve DID.\nError: %x\n", err)
		}

		if didDocumentStr == "" {
			return nil, fmt.Errorf("DID Document not found in VDR.")
		}

		fmt.Println("DID Document: ", didDocumentStr)

		// 4. DID Document에서 공개키 추출.
		//Json string을 DID Document 객체로 생성한다.
		didDocument, err := NewDIDDocumentForString(didDocumentStr)
		if err != nil {
			return nil, fmt.Errorf("Failed generate DID Document from string.\nError: %x\n", err)
		}

		// 첫 번째를 사용한다고 가정한다.
		// TODO: 키 ID(위의 kid)에 해당하는 키 값 구하기.
		pbKeyBaseMultibase := didDocument.VerificationMethod[0].PublicKeyMultibase
		_, bytePubKey, err := multibase.Decode(pbKeyBaseMultibase)
		pbKey, err := x509.ParsePKIXPublicKey(bytePubKey)

		return pbKey, nil //token이 검증되면 token.Valid는 true가 된다.
	})

	// 5. JWT 검증
	if claims, ok := token.Claims.(*JwtClaimsForVP); ok && token.Valid {
		//fmt.Printf("%v %v", claims.Vc, claims.Issuer)
		return true, claims, nil
	}

	return false, nil, err
}
