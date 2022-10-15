package AtomicUniversityIssuer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"os"
	"ssi-book/core"
	"ssi-book/protos"
)

type Server struct {
	protos.UnimplementedAtomicIssuerServer

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

func (server *Server) IssueAtomicVC(_ context.Context, msg *protos.MsgRequestAtomicVC) (*protos.MsgResponseAtomicVC, error) {
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

			response := new(protos.MsgResponseAtomicVC)

			if server.Issuer.CredentialSubjectJsonFilePath == "" {
				server.Issuer.CredentialSubjectJsonFilePath = "data/university_vc.json"
			}

			vcList, err := server.Issuer.GenerateAtomicVC()
			if err != nil {
				fmt.Println("Error Generate Atomic VC")
			}
			response.Result = "OK"
			response.Vcs = vcList

			return response, nil
		}
	}

	return nil, errors.New("Error")
}

func (issuer *Issuer) GenerateAtomicVC() ([]*protos.VC, error) {
	var vcList []*protos.VC
	var credentialSubject map[string]interface{}

	if issuer.CredentialSubjectJsonFilePath == "" {
		vcData := make(map[string]interface{})
		vcData["name"] = "HONG KIL DONG"
		credentialSubject = vcData
	} else {
		credentialSubject = loadJson(issuer.CredentialSubjectJsonFilePath) // "custom_vc.json"
	}

	// credentialSubject에서 diploma 클레임부분을 추출한다.
	diploma := credentialSubject["diploma"]

	// 클레임별로 VC를 만들어 vcList에 추가한다.
	switch diploma := diploma.(type) {
	case map[string]interface{}:
		for key, claim := range diploma {
			cs := make(map[string]interface{})

			id, _ := uuid.NewUUID()
			cs["id"] = id
			cs[key] = claim

			// VC 생성.
			vc, _ := core.NewVC(
				"1234567890",
				[]string{"VerifiableCredential", "DiplomaOfUniversity-" + key},
				issuer.did.String(),
				cs,
			)
			// VC에 Issuer의 private key로 서명한다.(JWT 사용)
			token, _ := vc.GenerateJWT(issuer.didDocument.VerificationMethod[0].Id, issuer.kms.PrivateKey)

			vcList = append(vcList, &protos.VC{Name: key, Token: token})
		}
	}

	return vcList, nil
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
