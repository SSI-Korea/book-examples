package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
)

func verifyProof() {
	vkFile, _ := os.Open("../datas/verification_key")
	rb1 := bufio.NewReader(vkFile)
	vk := groth16.NewVerifyingKey(ecc.BN254)
	vk.ReadFrom(rb1)
	vkFile.Close()

	pFile, _ := os.Open("../datas/proof")
	rb2 := bufio.NewReader(pFile)
	proof := groth16.NewProof(ecc.BN254)
	proof.ReadFrom(rb2)

	// 파일 읽기
	data, err := os.ReadFile("../datas/public_witness.bin")
	if err != nil {
		// 에러 처리
		panic("public witness file is not exist or wrong.")
	}

	// recreate a witness
	wit, _ := witness.New(ecc.BN254.ScalarField())
	// Binary unmarshalling
	err2 := wit.UnmarshalBinary(data)
	if err2 != nil {
		// 에러 처리
		panic("public witness UnmarshalBinary fail.")
	}

	pw, _ := wit.Public()

	err = groth16.Verify(proof, vk, pw)
	if err != nil {
		fmt.Println("증명 검증 실패: ", err)
	} else {
		fmt.Println("증명 검증 성공: 증명자는 19세 이상입니다.")
	}
}

func main() {
	verifyProof()
}
