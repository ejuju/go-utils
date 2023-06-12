package uid

import (
	"crypto/rand"
	"encoding/hex"
)

type ID []byte

func MustNewID(length int) ID {
	buf := make([]byte, length)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return buf
}

func (id ID) Hex() string   { return hex.EncodeToString(id) }
func (id ID) Bytes() []byte { return id }
