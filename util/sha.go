package util

import (
	"crypto/sha256"
	"encoding/hex"
)

func SHA256(data string) string {
	s := sha256.New()
	s.Write([]byte(data))
	bs := s.Sum(nil)

	rbs := make([]byte, s.BlockSize())
	hex.Encode(rbs, bs)

	return string(rbs)
}
