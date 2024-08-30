package main

import (
	"encoding/base64"
	"fmt"

	ml "github.com/IBM/mathlib"
	"github.com/hyperledger/aries-bbs-go/bbs"
)

func main() {

	pkBase64 := "sVEbbh9jDPGSBK/oT/EeXQwFvNuC+47rgq9cxXKrwo6G7k4JOY/vEcfgZw9Vf/TpArbIdIAJCFMDyTd7l2atS5zExAKX0B/9Z3E/mgIZeQJ81iZ/1HUnUCT2Om239KFx"
	pkBytes, _ := base64.RawStdEncoding.DecodeString(pkBase64)

	proofBase64 := "AAIBiN4EL9psRsIUlwQah7a5VROD369PPt09Z+jfzamP+/114a5RfWVMju3NCUl2Yv6ahyIdHGdEfxhC985ShlGQrRPLa+crFRiu2pfnAk+L6QMNooVMQhzJc2yYgktHen4QhsKV3IGoRRUs42zqPTP3BdqIPQeLgjDVi1d1LXEnP+WFQGEQmTKWTja4u1MsERdmAAAAdIb6HuFznhE3OByXN0Xp3E4hWQlocCdpExyNlSLh3LxK5duCI/WMM7ETTNS0Ozxe3gAAAAIuALkiwplgKW6YmvrEcllWSkG3H+uHEZzZGL6wq6Ac0SuktQ4n84tZPtMtR9vC1Rsu8f7Kwtbq1Kv4v02ct9cvj7LGcitzg3u/ZO516qLz+iitKeGeJhtFB8ggALcJOEsebPFl12cYwkieBbIHCBt4AAAAAxgEHt3iqKIyIQbTYJvtrMjGjT4zuimiZbtE3VXnqFmGaxVTeR7dh89PbPtsBI8LLMrCvFFpks9D/oTzxnw13RBmMgMlc1bcfQOmE9DZBGB7NCdwOnT7q4TVKhswOITKTQ=="
	proofBytes, _ := base64.StdEncoding.DecodeString(proofBase64)

	nonce := []byte("nonce")

	messagesBytes := [][]byte{[]byte("message1"), []byte("message2")}
	revealedMessagesBytes := messagesBytes[:1]

	bls := bbs.New(ml.Curves[ml.BLS12_381_BBS])

	err := bls.VerifyProof(revealedMessagesBytes, proofBytes, nonce, pkBytes)
	if err != nil {
		fmt.Println("VerifyProof Error. ", err)
	} else {
		fmt.Println("VerifyProof OK.")
	}

}
