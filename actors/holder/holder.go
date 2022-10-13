package holder

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"os"
	"ssi-book/core"
	"ssi-book/protos"
	"time"
)

type Holder struct {
	Kms         *core.ECDSAManager
	Did         *core.DID
	DidDocument *core.DIDDocument
	VCList      []string
}

func (holder *Holder) GenerateDID() {
	// 키생성(ECDSA) - 향후 KMS로 대체.
	holder.Kms = core.NewEcdsa()

	// DID 생성.
	did, _ := core.NewDID("ssikr", holder.Kms.PublicKeyBase58())

	holder.Did = did

	// DID Document 생성.
	verificationId := fmt.Sprintf("%s#keys-1", did)
	verificationMethod := []core.VerificationMethod{
		{
			Id:                 verificationId,
			Type:               core.VERIFICATION_KEY_TYPE_SECP256K1,
			Controller:         did.String(),
			PublicKeyMultibase: holder.Kms.PublicKeyMultibase(),
		},
	}
	didDocument := core.NewDIDDocument(did.String(), verificationMethod)
	holder.DidDocument = didDocument
}

func (holder *Holder) GenerateFirstVC() {
	// VC 생성.
	vc, _ := core.NewVC(
		"1234567890",
		[]string{"VerifiableCredential", "SelfCertification"},
		holder.Did.String(), // Issuer did
		map[string]interface{}{
			"id":        "1234567890",
			"name":      "HONG KIL DONG",
			"mobile":    "010-1234-1234",
			"birthDate": "2000-01-01",
			"gender":    "M",
		},
	)

	vcJwt, _ := vc.GenerateJWT(holder.DidDocument.VerificationMethod[0].Id, holder.Kms.PrivateKey)
	holder.VCList = append(holder.VCList, vcJwt)
}

func (holder *Holder) GenerateVP() (string, error) {

	vcList := holder.VCList

	vp, err := core.NewVP(
		"12345678901111",
		[]string{"VerifiablePresentaion"},
		holder.Did.String(),
		vcList,
	)
	if err != nil {
		fmt.Println("ERROR")
		os.Exit(0)
	}

	vpToken := vp.GenerateJWT(holder.DidDocument.VerificationMethod[0].Id, holder.Kms.PrivateKey)

	return vpToken, nil
}

func (holder *Holder) RequestVCToUniversityIssuer(vpToken string) error {
	conn, err := grpc.Dial("localhost:1121", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("UniversityIssuer not connect: %v", err)
		return err
	}
	defer conn.Close()
	c := protos.NewSimpleIssuerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.IssueSimpleVC(ctx, &protos.MsgRequestVC{
		Did: holder.Did.String(),
		Vp:  vpToken,
	})
	if err != nil {
		log.Printf("could not request: %v", err)
		return err
	}

	fmt.Printf("UniversityIssuer's response: %s\n", res.Result)
	fmt.Printf("UniversityIssuer's response VC: %s\n", res.Vc)
	if res.Result == "OK" {
		holder.VCList = append(holder.VCList, res.Vc)
	}

	return nil
}

func (holder *Holder) RequestVCToCompanyIssuer(vpToken string) error {
	conn, err := grpc.Dial("localhost:1122", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("CompanyIssuer not connect: %v", err)
		return err
	}
	defer conn.Close()
	c := protos.NewSimpleIssuerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.IssueSimpleVC(ctx, &protos.MsgRequestVC{
		Did: holder.Did.String(),
		Vp:  vpToken,
	})
	if err != nil {
		log.Printf("could not request: %v", err)
		return err
	}

	fmt.Printf("CompanyIssuer's response: %s\n", res.Result)
	fmt.Printf("CompanyIssuer's response VC: %s\n", res.Vc)
	if res.Result == "OK" {
		holder.VCList = append(holder.VCList, res.Vc)
	}

	return nil
}

func (holder *Holder) RequestVCToBankIssuer(vpToken string) error {
	conn, err := grpc.Dial("localhost:1123", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("BankIssuer not connect: %v", err)
		return err
	}
	defer conn.Close()
	c := protos.NewMultipleIssuerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.IssueMultipleVC(ctx, &protos.MsgRequestMultipleVC{
		Did: holder.Did.String(),
		Vp:  vpToken,
	})
	if err != nil {
		log.Printf("could not request: %v", err)
		return err
	}

	fmt.Printf("BankIssuer's response: %s\n", res.Result)

	if res.Result == "OK" {
		for _, vc := range res.Vc {
			fmt.Printf("BankIssuer's response VC: %s\n", vc)
			holder.VCList = append(holder.VCList, vc)
		}
	}

	return nil
}
