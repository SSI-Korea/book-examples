package BankIssuer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/exp/slices"
	"io/ioutil"
	"log"
	"math"
	"os"
	"ssi-book/core"
	"ssi-book/protos"
	"time"
)

type Server struct {
	protos.UnimplementedMultipleIssuerServer

	Issuer *Issuer
}

type Issuer struct {
	kms         *core.ECDSAManager
	did         *core.DID
	didDocument *core.DIDDocument

	CredentialSubjectJsonFilePath string
}

const (
	VC_MODE_TEST   = 1
	VC_MODE_CUSTOM = 2
)

var (
	vcMode = VC_MODE_TEST
)

type VC_CUSTOM_CLAIM struct {
	data map[string]interface{}
}

func (server *Server) IssueMultipleVC(_ context.Context, msg *protos.MsgRequestMultipleVC) (*protos.MsgResponseMultipleVC, error) {
	log.Printf("IssueMultipleVC MSG: %+v \n", msg)
	isVerify, claims, err := core.ParseAndVerifyJwtForVP(msg.Vp)
	if !isVerify || err != nil {
		fmt.Println("VP is NOT verified.")
		return nil, errors.New(fmt.Sprintf("VP is invalid: %s", err))
	}

	fmt.Println("VP is verified.")

	checkSelfCertification, checkDiploma, checkEmployee := false, false, false

	var vcEmployee map[string]interface{}

	for i, vc := range claims.Vp.VerifiableCredential {
		fmt.Println("VC: ", vc)
		isVerify, claims, err := core.ParseAndVerifyJwtForVC(vc)
		if !isVerify || err != nil {
			fmt.Println("VC #", i, " is NOT verified.")
			return nil, errors.New(fmt.Sprintf("VC is invalid: %s", err))
		}

		fmt.Println("VC is verified.")

		vcType := claims.Vc.Type
		vcClaims := claims.Vc.CredentialSubject
		if slices.Contains(vcType, "SelfCertification") && vcClaims["name"] == "HONG KIL DONG" && vcClaims["birthDate"] == "2000-01-01" {
			checkSelfCertification = true
		}
		if slices.Contains(vcType, "DiplomaOfUniversity") {
			vcDiploma := vcClaims["diploma"].(map[string]interface{})
			if vcDiploma["join"] == "2015-03-01" && vcDiploma["graduation"] == "2020-03-01" {
				checkDiploma = true
			}
		}
		if slices.Contains(vcType, "CertificateOfEmployment") {
			vcEmployee = vcClaims["employee"].(map[string]interface{})
			if vcEmployee["name"] == "HONG KIL DONG" {
				checkEmployee = true
			}

		}
	}

	if checkSelfCertification && checkDiploma && checkEmployee {
		fmt.Println("VC 발급!!!!")

		response := new(protos.MsgResponseMultipleVC)

		server.Issuer.CredentialSubjectJsonFilePath = "data/bank_account_vc.json"

		vcAccountToken, err := server.Issuer.GenerateSampleVC("AccountCredential")
		if err != nil {
			return nil, errors.New("VC Generate Error")
		}

		response.Result = "OK"
		response.Vc = append(response.Vc, vcAccountToken)

		nowDate := time.Now()
		joinDate, _ := time.Parse("2006-01-02", vcEmployee["join"].(string))
		diff := nowDate.Sub(joinDate)
		diffDays := math.Floor(diff.Minutes() / 60 / 24)

		if diffDays >= 180 {
			server.Issuer.CredentialSubjectJsonFilePath = "data/bank_loan_vc.json"

			vcLoanToken, err := server.Issuer.GenerateSampleVC("LoanCredential")
			if err != nil {
				return nil, errors.New("VC Generate Error")
			}
			response.Vc = append(response.Vc, vcLoanToken)
		}

		return response, nil
	}

	return nil, errors.New("VC condition Error")
}

func (issuer *Issuer) GenerateDID() {
	// 키생성(ECDSA)
	issuer.kms = core.NewEcdsa()

	// DID 생성.
	issuerDid, _ := core.NewDID("ssikr", issuer.kms.PublicKeyBase58())

	issuer.did = issuerDid

	// DID Document 생성.
	verificationId := fmt.Sprintf("%s#keys-1", issuerDid)
	verificationMethod := []core.VerificationMethod{
		{
			Id:                 verificationId,
			Type:               core.VERIFICATION_KEY_TYPE_SECP256K1,
			Controller:         issuerDid.String(),
			PublicKeyMultibase: issuer.kms.PublicKeyMultibase(),
		},
	}
	didDocument := core.NewDIDDocument(issuerDid.String(), verificationMethod)
	issuer.didDocument = didDocument

	RegisterDid(issuerDid.String(), didDocument)
}

func (issuer *Issuer) GenerateSampleVC(typ string) (string, error) {

	var credentialSubject map[string]interface{}

	if issuer.CredentialSubjectJsonFilePath == "" {

	} else {
		credentialSubject = LoadJson(issuer.CredentialSubjectJsonFilePath) // "custom_vc.json"
	}

	// VC 생성.
	vc, err := core.NewVC(
		"1234567890abccde",
		[]string{"VerifiableCredential", typ},
		issuer.did.String(),
		credentialSubject,
	)

	if err != nil {
		return "", errors.New("Failed creation VC.")
	}

	// VC에 Issuer의 private key로 서명한다.(JWT 사용)
	token, err := vc.GenerateJWT(issuer.didDocument.VerificationMethod[0].Id, issuer.kms.PrivateKey)

	return token, nil
}

func RegisterDid(did string, document *core.DIDDocument) error {
	err := core.RegisterDid(did, document.String())
	if err != nil {
		return err
	}
	return nil

}

func LoadJson(path string) map[string]interface{} {
	jsonData := make(map[string]interface{})

	data, err := os.Open(path)
	if err != nil {
		return nil
	}

	vcMode = VC_MODE_CUSTOM
	byteValue, _ := ioutil.ReadAll(data)

	json.Unmarshal(byteValue, &jsonData)

	return jsonData
}
