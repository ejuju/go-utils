package uid

import (
	"crypto/rand"
	"encoding/hex"
)

type ID []byte

func NewID(length int) (ID, error) {
	buf := make([]byte, length)
	_, err := rand.Read(buf)
	return buf, err
}

func MustNewID(length int) ID {
	id, err := NewID(length)
	if err != nil {
		panic(err)
	}
	return id
}

func (id ID) Hex() string    { return hex.EncodeToString(id) }
func (id ID) String() string { return id.Hex() }
func (id ID) Bytes() []byte  { return id }
