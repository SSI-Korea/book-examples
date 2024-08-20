package RootofTrustIssuer

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

	registerDid(issuerDid.String(), didDocument)
}

func (server *Server) IssueSimpleVC(_ context.Context, msg *protos.MsgRequestVC) (*protos.MsgResponseVC, error) {
	log.Printf("RootOfTrust: %+v \n", msg)

	fmt.Println("ROT VC 발급!!!!")

	response := new(protos.MsgResponseVC)

	if server.Issuer.CredentialSubjectJsonFilePath == "" {
		server.Issuer.CredentialSubjectJsonFilePath = "data/rot_vc.json"
	}

	vcToken, err := server.Issuer.GenerateSampleVC()
	if err != nil {
	}

	response.Result = "OK"
	response.Vc = vcToken

	return response, nil

}

func (issuer *Issuer) GenerateSampleVC() (string, error) {

	var credentialSubject map[string]interface{}

	if issuer.CredentialSubjectJsonFilePath == "" {
		vcData := make(map[string]interface{})
		vcData["name"] = "HONG KIL DONG"
		credentialSubject = vcData
	} else {
		credentialSubject = loadJson(issuer.CredentialSubjectJsonFilePath) // "custom_vc.json"
	}

	// VC 생성.
	vc, err := core.NewVC(
		"1234567890",
		[]string{"VerifiableCredential", "ResidentRegistrationCredential"},
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

func registerDid(did string, document *core.DIDDocument) error {
	err := core.RegisterDid(did, document.String())
	if err != nil {
		return err
	}

	return nil
}

func loadJson(path string) map[string]interface{} {
	jsonData := make(map[string]interface{})

	data, err := os.Open(path)
	if err != nil {
		return nil
	}

	byteValue, _ := ioutil.ReadAll(data)

	json.Unmarshal(byteValue, &jsonData)

	return jsonData
}
