package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	ml "github.com/IBM/mathlib"
	"github.com/hyperledger/aries-bbs-go/bbs"
)

func generateKeyPairRandom(curve *ml.Curve) (*bbs.PublicKey, *bbs.PrivateKey, error) {
	seed := make([]byte, 32)

	_, err := rand.Read(seed)
	if err != nil {
		panic(err)
	}

	bbs := bbs.NewBBSLib(curve)

	return bbs.GenerateKeyPair(sha256.New, seed)
}

func main() {
	// 1. 키 생성
	var curve = ml.Curves[ml.BLS12_381_BBS]

	publicKey, privateKey, err := generateKeyPairRandom(curve)
	if err != nil {
		panic(err)
	}

	privKeyBytes, _ := privateKey.Marshal()
	pubKeyBytes, _ := publicKey.Marshal()

	// 2. 메시지 준비
	messagesBytes := [][]byte{
		[]byte("Alice"),
		[]byte("1990-01-01"), // 생년월일
		[]byte("Seoul"),
		[]byte("male"),
	}

	bls := bbs.New(curve)

	// 3. 서명 생성
	signatureBytes, err := bls.Sign(messagesBytes, privKeyBytes)
	// fmt.Println(signatureBytes)
	if err != nil {
		panic(err)
	}

	err = bls.Verify(messagesBytes, signatureBytes, pubKeyBytes)
	if err != nil {
		panic(err)
	}

	// 4. 증명 생성 (나이가 30세 이상이라는 증명)
	nonce := []byte("nonce")
	// revealedIndexes := []int{0, 2}
	revealedIndexes := []int{1} // 생년월일 인덱스

	proofBytes, err := bls.DeriveProof(messagesBytes, signatureBytes, nonce, pubKeyBytes, revealedIndexes)
	// fmt.Println(proofBytes)
	if err != nil {
		panic(err)
	}

	revealedMessages := make([][]byte, len(revealedIndexes))
	for i, ind := range revealedIndexes {
		revealedMessages[i] = messagesBytes[ind]
	}

	fmt.Println(string(revealedMessages[0]))

	// 5. 증명 검증
	err = bls.VerifyProof(revealedMessages, proofBytes, nonce, pubKeyBytes)

	if err != nil {
		fmt.Println("Verification failed.", err)
	} else {
		fmt.Println("Verification success")
	}
}
