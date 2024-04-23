package crypto

import (
	"golang.org/x/crypto/sha3"
)

const hashSize = 32

func Hash(ms ...[]byte) []byte {
	h := sha3.NewShake128()
	for _, m := range ms {
		h.Write(m)
	}
	ret := make([]byte, hashSize)
	h.Read(ret)

	return ret
}
