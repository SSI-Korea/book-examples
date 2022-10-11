package CompanyIssuer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"ssi-book/core"
	"ssi-book/protos"
)

type Server struct {
	protos.UnimplementedSimpleIssuerServer

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

func (server *Server) IssueSimpleVC(_ context.Context, msg *protos.MsgRequestVC) (*protos.MsgResponseVC, error) {
	log.Printf("IssueSimpleVC MSG: %+v \n", msg)
	isVerify, claims, err := core.ParseAndVerifyJwtForVP(msg.Vp)
	if !isVerify || err != nil {
		fmt.Println("VP is NOT verified.")
		return nil, errors.New(fmt.Sprintf("VP is invalid: %s", err))
	}

	fmt.Println("VP is verified.")

	for i, vc := range claims.Vp.VerifiableCredential {
		fmt.Println("VC: ", vc)
		isVerify, claims, err := core.ParseAndVerifyJwtForVC(vc)
		if !isVerify || err != nil {
			fmt.Println("VC #", i, " is NOT verified.")
			return nil, errors.New(fmt.Sprintf("VC is invalid: %s", err))
		}
		fmt.Println("VC is verified.")
		vcClaims := claims.Vc.CredentialSubject
		if vcClaims["name"] == "HONG KIL DONG" && vcClaims["birthDate"] == "2000-01-01" {

			fmt.Println("VC 발급!!!!")

			response := new(protos.MsgResponseVC)

			server.Issuer.CredentialSubjectJsonFilePath = "university_vc.json"

			vcToken, err := server.Issuer.GenerateSampleVC()
			if err != nil {

			}
			response.Result = "OK"
			response.Vc = vcToken

			return response, nil
		}
	}

	return nil, errors.New("Error")

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

func (issuer *Issuer) GenerateSampleVC() (string, error) {

	var credentialSubject map[string]interface{}

	if issuer.CredentialSubjectJsonFilePath == "" {

	} else {
		credentialSubject = LoadJson(issuer.CredentialSubjectJsonFilePath) // "custom_vc.json"
	}

	// VC 생성.
	vc, err := core.NewVC(
		"1234567890",
		[]string{"VerifiableCredential", "DiplomaOfUniversity"},
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
