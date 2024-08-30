package main

import (
	"encoding/base64"
	"fmt"

	ml "github.com/IBM/mathlib"
	"github.com/hyperledger/aries-bbs-go/bbs"
)

func main() {
	pkBase64 := "l0Wtf3gy5f140G5vCoCJw2420hwk6Xw65/DX3ycv1W7/eMky8DyExw+o1s2bmq3sEIJatkiN8f5D4k0766x0UvfbupFX+vVkeqnlOvT6o2cag2osQdMFbBQqAybOM4Gm"
	pkBytes, _ := base64.RawStdEncoding.DecodeString(pkBase64)

	proofBase64 := "AAQFpAE2VALtmriOzSMk/oqid4uJhPQRUVUuyenL/L4w4ykdyh0jCX64EFqCdLP+n8VrkOKXhHPKPoCOdHBOMv96aM15NFg867/MToMeNN0IFzZkzhs37qk1vWWFKReMF+cRsCAmkHO6An1goNHdY/4XquSV3LwykezraWt8+8bLvVn6ciaXBVxVcYkbIXRsVjqbAAAAdIl/C/W5G1pDbLMrUrBAYdpvzGHG25gktAuUFZb/SkIyy0uhtWJk2v6A+D3zkoEBsgAAAAJY/jfJR9kpGbSY5pfz+qPkqyNOTJbs6OEpfBwYGsyC7hspvBGUOYyvuKlS8SvKAXW7hVawAhYJbvnRwzeiP6P9kbZKtLQZIkRQB+mxRSbMk/0JgE1jApHOlPtgbqI9yIouhK9xT2wVZl79qTAwifonAAAABDTDo5VtXR2gloy+au7ai0wcnnzjMJ6ztQHRI1ApV5VuOQ19TYL7SW+C90p3QSZFQ5gtl90PHaUuEAHIb+7ZgbJvh5sc1DjKfThwPx0Ao0w8+xTbLhNlxvo6VE1cfbiuME+miCAibLgHjksQ8ctl322qnblYJLXiS4lvx/jtGvA3"
	proofBytes, _ := base64.StdEncoding.DecodeString(proofBase64)

	nonce := []byte("nonce")

	messagesBytes := [][]byte{
		[]byte("message1"),
		[]byte("message2"),
		[]byte("message3"),
		[]byte("message4"),
	}
	revealedMessagesBytes := [][]byte{messagesBytes[0], messagesBytes[2]}

	bls := bbs.New(ml.Curves[ml.BLS12_381_BBS])

	err := bls.VerifyProof(revealedMessagesBytes, proofBytes, nonce, pkBytes)
	if err != nil {
		fmt.Println("VerifyProof Error. ", err)
	} else {
		fmt.Println("VerifyProof OK.")
	}
}
