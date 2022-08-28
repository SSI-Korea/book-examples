package main

import (
	"fmt"
)

// https://www.w3.org/TR/did-core/
func main() {
	method := "ssikr"
	specificIdentifier := "abcd1234"

	// [did:DID Method:DID Method-Specific Identifier]
	did := fmt.Sprintf("did:%s:%s", method, specificIdentifier)

	fmt.Printf("DID: %s\n", did)
}
