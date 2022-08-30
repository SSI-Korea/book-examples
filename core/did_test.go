package core

import (
	"fmt"
	"strings"
	"testing"
)

func TestGenerateDID(t *testing.T) {
	method := "ssikr"
	// 키생성(ECDSA) - 향후 KMS로 대체.
	kms := NewEcdsa()

	// DID 생성.
	did, err := NewDID(method, kms.PublicKeyBase58())

	if err != nil {
		t.Error("Failed to generate DID.")
	}

	if did == nil || !strings.HasPrefix(did.String(), fmt.Sprintf("did:%s:", method)) {
		t.Error("Failed to generate DID.")
	}
}
