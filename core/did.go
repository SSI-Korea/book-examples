// core/did.go

package core

import (
	"errors"
	"fmt"
	"ssi-book/util"
)

type DID struct {
	did    string
	method string
}

func NewDID(method string, pbKey string) (did *DID, err error) {
	if method == "" || pbKey == "" {
		return nil, errors.New("parameter is not valid")
	}

	var newDid = new(DID)
	newDid.method = method
	specificIdentifier := util.MakeHashBase58(pbKey)
	// DID:Method:specific
	newDid.did = fmt.Sprintf("did:%s:%s", method, specificIdentifier)

	return newDid, nil
}

func (d *DID) String() string {
	return d.did
}
