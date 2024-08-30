package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	ml "github.com/IBM/mathlib"
	"github.com/hyperledger/aries-bbs-go/bbs"
)

func generateKeyPairRandom(curve *ml.Curve) (*bbs.PublicKey, *bbs.PrivateKey, error) {
	seed := make([]byte, 32)

	_, err := rand.Read(seed)
	if err != nil {
		panic(err)
	}

	bbs := bbs.NewBBSLib(curve)

	return bbs.GenerateKeyPair(sha256.New, seed)
}

var curve = ml.Curves[ml.BLS12_381_BBS]

func main() {

	pubKey, privKey, _ := generateKeyPairRandom(curve)

	pubKeyBytes, _ := pubKey.Marshal()

	blindMsgCount := 2

	messagesBytes := [][]byte{
		[]byte("message1"),
		[]byte("message2"),
		[]byte("message3"),
		[]byte("message4"),
	}

	pubKeyWithGenerators, err := pubKey.ToPublicKeyWithGenerators(len(messagesBytes))
	if err != nil {
		fmt.Println("ToPublicKeyWithGenerators Error. ", err)
	}

	blindedMessagesBytes := [][]byte{
		[]byte("message1"),
		nil,
		nil,
		[]byte("message4"),
	}

	clearMessagesBytes := [][]byte{
		nil,
		[]byte("message2"),
		[]byte("message3"),
		nil,
	}

	// requester generates commitment to blind messages
	cb := bbs.NewCommitmentBuilder(blindMsgCount + 1)
	for i, msg := range blindedMessagesBytes {
		if msg == nil {
			continue
		}

		cb.Add(pubKeyWithGenerators.H[i], bbs.FrFromOKM(msg, curve))
	}
	blinding := curve.NewRandomZr(rand.Reader)
	cb.Add(pubKeyWithGenerators.H0, blinding)
	b_req := cb.Build()

	// signer adds its component
	cb = bbs.NewCommitmentBuilder(len(messagesBytes) - blindMsgCount + 2)
	for i, msg := range clearMessagesBytes {
		if msg == nil {
			continue
		}

		cb.Add(pubKeyWithGenerators.H[i], bbs.FrFromOKM(msg, curve))
	}
	cb.Add(b_req, curve.NewZrFromInt(1))
	cb.Add(curve.GenG1, curve.NewZrFromInt(1))
	comm := cb.Build()

	// signer signs
	scheme := bbs.New(curve)
	sig, err := scheme.SignWithKeyB(comm, len(messagesBytes), privKey)
	if err != nil {
		fmt.Println("SignWithKeyB Error. ", err)
	}

	// requester unblinds
	signature, err := bbs.NewBBSLib(curve).ParseSignature(sig)
	if err != nil {
		fmt.Println("NewBBSLib Error. ", err)
	}

	signature.S = curve.ModAdd(signature.S, blinding, curve.GroupOrder)
	sig, _ = signature.ToBytes()

	// requester verifies
	err = scheme.Verify(messagesBytes, sig, pubKeyBytes)
	if err != nil {
		fmt.Println("Verify Error. ", err)
	}

}
