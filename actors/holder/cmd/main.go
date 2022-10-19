package main

import (
	"fmt"
	"ssi-book/actors/holder"
	"ssi-book/core"
	"ssi-book/util"
)

func main() {
	fmt.Println("### Start HOLDER's Wallet ###")
	// New Holder
	hldr := new(holder.Holder)
	hldr.AtomicVCList = make(map[string]string)

	util.PressKey("1. DID를 생성합니다. [아무키나 입력하세요.]")
	hldr.GenerateDID()
	fmt.Printf("DID: %s\n", hldr.Did.String())
	fmt.Printf("DID Document: %+v\n", hldr.DidDocument)

	util.PressKey("2. DID를 VDR에 등록합니다. [아무키나 입력하세요.]")
	core.RegisterDid(hldr.Did.String(), hldr.DidDocument.String())

	// 최초 VC를 발급한다.
	util.PressKey("3. 최초 VC를 발급합니다. [아무키나 입력하세요.]")
	hldr.GenerateFirstVC()
	fmt.Println("First VC: ", hldr.VCList[0])

	// UniversityIssuer에게 졸업증명 VC를 요청한다.
	util.PressKey("4. UniversityIssuer에게 졸업증명 VC를 요청한다. [아무키나 입력하세요.]")
	vpToken, _ := hldr.GenerateVP()

	fmt.Printf("VP Token: %s\n", vpToken)

	hldr.RequestVCToUniversityIssuer(vpToken)

	// CompanyIssuer에게 재직증명 VC를 요청한다.
	util.PressKey("5. CompanyIssuer에게 재직증명 VC를 요청한다. [아무키나 입력하세요.]")
	vpToken, _ = hldr.GenerateVP()

	fmt.Printf("VP Token: %s\n", vpToken)

	hldr.RequestVCToCompanyIssuer(vpToken)

	// BankIssuer에게 재직증명 VC를 요청한다.
	util.PressKey("6. BankIssuer에게 계좌 VC와 대출 VC를 요청한다. [아무키나 입력하세요.]")
	vpToken, _ = hldr.GenerateVP()

	fmt.Printf("VP Token: %s\n", vpToken)

	hldr.RequestVCToBankIssuer(vpToken)

	// AtomicUniversityIssuer에게 졸업증명 Atomic VC를 요청한다.
	util.PressKey("7. AtomicUniversityIssuer에게 졸업증명 Atomic VC를 요청한다. [아무키나 입력하세요.]")
	vpToken, _ = hldr.GenerateVP()

	hldr.RequestVCToAtomicUniversityIssuer(vpToken)

	// AtomicUniversityIssuer에게 졸업증명 Atomic VC를 요청한다.
	util.PressKey("8. Atomic VC 목록을 출력한다. [아무키나 입력하세요.]")
	hldr.PrintAtomicVC()

}
