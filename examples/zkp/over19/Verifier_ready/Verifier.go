package main

import (
	"bufio"
	"os"
	"ssi-book/examples/zkp/over19"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

func main() {
	// 1. 회로 컴파일
	var circuit over19.AgeCircuit
	ccs, _ := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)

	// 2. 설정 단계
	provingKey, verificationKey, _ := groth16.Setup(ccs)

	// 3. 공개 파라미터(검증 키와 증명 키) 저장 또는 공유
	verificationKeyFile, _ := os.Create("../datas/verification_key")
	wb1 := bufio.NewWriter(verificationKeyFile)
	verificationKey.WriteRawTo(wb1)
	wb1.Flush()
	verificationKeyFile.Close()

	provingKeyFile, _ := os.Create("../datas/proving_key")
	wb2 := bufio.NewWriter(provingKeyFile)
	provingKey.WriteRawTo(wb2) // alternatively, provingKey.WriteTo(&buf)
	wb2.Flush()
	provingKeyFile.Close()

	// 4. 메타데이터 저장 (증명자를 위한 정보)
	metaData := []byte(`{
        "circuit": "AgeVerification",
        "inputs": {
            "Age": {
                "type": "secret",
                "description": "Your actual age (must be an integer)"
            },
            "MinAge": {
                "type": "public",
                "description": "Minimum age for verification (set to 19)",
                "value": 19
            }
        },
        "description": "Prove that your age is at least 19 without revealing your actual age"
    }`)

	// 메타데이터를 파일로 저장
	os.WriteFile("../datas/circuit_metadata.json", metaData, 0644)
}
