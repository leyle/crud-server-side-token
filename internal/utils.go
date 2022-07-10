package internal

import (
	"crypto/md5"
	"encoding/hex"
)

const keyLength = 32

func makeKeyLength32(key []byte) []byte {
	if len(key) == keyLength {
		return key
	}
	m := md5.New()
	m.Write(key)
	tmp := hex.EncodeToString(m.Sum(nil))
	return []byte(tmp)
}
