package main

import (
	"fmt"
	"os"
	"ssi-book/core"
)

func main() {
	issuerKeyEcdsa := core.NewEcdsa()

	// DID 생성.
	issuerDid, _ := core.NewDID("ssikr", issuerKeyEcdsa.PublicKeyBase58())

	// VC 생성.
	vc, err := core.NewVC(
		"1234567890",
		[]string{"VerifiableCredential", "AlumniCredential"},
		issuerDid.String(),
		map[string]interface{}{
			"id": "1234567890",
			"alumniOf": map[string]interface{}{
				"id": "1234567",
				"name": []map[string]string{
					{
						"value": "Example University",
						"lang":  "en",
					}, {
						"value": "Exemple d'Université",
						"lang":  "fr",
					},
				},
			},
		},
	)

	if err != nil {
		fmt.Println("Failed creation VC.")
		os.Exit(0)
	}

	fmt.Println("### New VC test ###")
	fmt.Println(vc)
}
