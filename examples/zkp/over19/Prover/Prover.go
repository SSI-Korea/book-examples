package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"ssi-book/examples/zkp/over19"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

func main() {
	// 1. 메타데이터 로드
	// JSON 파일 경로
	filePath := "../datas/circuit_metadata.json"

	// 파일 열기
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// JSON 데이터 디코딩
	decoder := json.NewDecoder(file)
	var config over19.CircuitConfig
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// 데이터 사용
	var minAge int
	for inputName, inputInfo := range config.Inputs {
		if inputName == "MinAge" {
			minAge = inputInfo.Value
		}
	}

	// 2. 공개 파라미터(증명 키) 로드
	pkFile, _ := os.Open("../datas/proving_key")
	w := bufio.NewReader(pkFile)
	pk := groth16.NewProvingKey(ecc.BN254)
	pk.ReadFrom(w)
	pkFile.Close()

	// 3. 입력값 준비
	assignment := over19.AgeCircuit{
		Age:    20, // 증명자의 실제 나이 (비밀 입력)
		MinAge: minAge,
	}

	// 4. 위트니스를 JSON으로 변환
	// witnessJSON, _ := json.Marshal(witness)

	// 5. JSON 위트니스로부터 gnark 위트니스 생성
	witness, errNewWitness := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	if errNewWitness != nil {
		fmt.Println("Witness Error: ", errNewWitness)
		panic("witness fail.")
	}

	// fmt.Println("witness: ", witness)

	publicWitness, _ := witness.Public()

	// Binary marshalling
	data, _ := publicWitness.MarshalBinary()

	// 파일 저장
	err0 := os.WriteFile("../datas/public_witness.bin", data, 0644)
	if err0 != nil {
		// 에러 처리
		panic("public_witness.bin fail.")
	}

	var circuit over19.AgeCircuit
	ccs, _ := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)

	// 6. 증명 생성
	proof, errProve := groth16.Prove(ccs, pk, witness)
	if errProve != nil {
		fmt.Println("Proof Error: ", errProve)
		panic("Proof generate fail.")
	}

	// 7. 증명을 파일로 저장 (검증자에게 전송하기 위해)
	// proof.WriteFile("proof.json")
	proofFile, _ := os.Create("../datas/proof")
	proof.WriteTo(proofFile)
	proofFile.Close()

	fmt.Println("증명이 생성되어 proof 파일로 저장되었습니다.")
}
