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

// Tries 5 times before panicking
func MustNewID(length int) ID {
	var id ID
	var err error
	for i := 0; i < 5; i++ {
		id, err = NewID(length)
		if err == nil {
			return id
		}
	}
	panic(err)
}

func (id ID) Hex() string    { return hex.EncodeToString(id) }
func (id ID) String() string { return id.Hex() }
func (id ID) Bytes() []byte  { return id }
