// util/utils.go

package util

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"os"
)

func MakeHash(plain string) []byte {
	digest := sha256.Sum256([]byte(plain))
	return digest[:]
}

func MakeHashBase58(plain string) string {
	return base58.Encode(MakeHash(plain))
}

func MakeHashHex(plain string) string {
	return hex.EncodeToString(MakeHash(plain))
}

func PressKey(msg string) {
	kbReader := bufio.NewReader(os.Stdin)

	fmt.Println(msg)
	kbReader.ReadString('\n')
}
