package kv

import (
	"errors"
	"fmt"
	"os"
)

type WriteOp string

const (
	WriteOpUnknown     WriteOp = "unknown row kind"
	WriteOpPutKey      WriteOp = "put key"
	WriteOpPutKeyValue WriteOp = "put key-value"
	WriteOpDelete      WriteOp = "delete row by key"
)

var (
	ErrKeyNotFound      = errors.New("key not found")
	ErrKeyEmpty         = errors.New("key is empty")
	ErrKeyAlreadyExists = errors.New("key already exists")
	ErrUnknownWriteOp   = errors.New("unknown write operation")
)

// Opens a read-only and a write-only file handler.
func openFileROWO(fpath string) (*os.File, *os.File, error) {
	roFile, err := os.OpenFile(fpath, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, nil, fmt.Errorf("open read-only file: %w", err)
	}
	woFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return nil, nil, fmt.Errorf("open write-only file: %w", err)
	}

	return roFile, woFile, nil
}
