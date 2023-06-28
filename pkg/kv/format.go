package kv

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

// Holds characters used for encoding and decoding row data.
type Format struct {
	PutKey        byte // Default: '-'
	PutKeyValue   byte // Default: '='
	Delete        byte // Default: '!'
	KeyPrefix     byte // Default: ' ' (whitespace)
	ValuePrefix   byte // Default: ' ' (whitespace)
	RowEnd        byte // Default: '\n' (line-break)
	SizeStart     byte // Default: '('
	SizeSeparator byte // Default: '|'
	SizeEnd       byte // Default: ')'
	SizeBase      int  // Default: 10 (decimal)
}

// TextFileFormat with default values.
var DefaultFormat = &Format{
	PutKey:      '-',
	PutKeyValue: '=',
	Delete:      '!',
	KeyPrefix:   ' ',
	ValuePrefix: ' ',
	RowEnd:      '\n',
}

func (chars *Format) Encode(kind WriteOp, k, v []byte) ([]byte, error) {
	if len(k) == 0 {
		return nil, ErrKeyEmpty
	}
	if (kind == WriteOpDelete || kind == WriteOpPutKey) && len(v) > 0 {
		return nil, fmt.Errorf("row %q must receive nil value", kind)
	}
	if (kind == WriteOpDelete || kind == WriteOpPutKey) && bytes.IndexByte(k, chars.RowEnd) != -1 {
		return nil, fmt.Errorf("invalid key contains row end %q: %q", chars.RowEnd, k)
	}
	if kind == WriteOpPutKeyValue && bytes.IndexByte(k, chars.ValuePrefix) != -1 {
		return nil, fmt.Errorf("invalid key contains value prefix %q: %q", chars.ValuePrefix, k)
	}
	if kind == WriteOpPutKeyValue && bytes.IndexByte(v, chars.RowEnd) != -1 {
		return nil, fmt.Errorf("invalid value contains row end %q: %q", chars.RowEnd, k)
	}

	out := []byte{}

	// Add first two characters: command + key-prefix
	switch kind {
	default:
		return nil, fmt.Errorf("%w: %q", ErrUnknownWriteOp, kind)
	case WriteOpPutKey:
		out = append(out, chars.PutKey)
	case WriteOpPutKeyValue:
		out = append(out, chars.PutKeyValue)
	case WriteOpDelete:
		out = append(out, chars.Delete)
	}
	out = append(out, chars.KeyPrefix)

	// Add key
	out = append(out, k...)

	// Add value if row is of kind "put key-value"
	if kind == WriteOpPutKeyValue {
		out = append(out, chars.ValuePrefix)
		out = append(out, v...)
	}

	out = append(out, chars.RowEnd)
	return out, nil
}

// Parse a row assuming a byte slice containing the whole row (without trailing line-break)
func (chars *Format) ParseRowFromBytes(b []byte) (WriteOp, []byte, []byte, error) {
	// Fail if row is less than 3 characters long.
	// For ex: the shortest possible row is `- 1` (cmd + key-prefix + key + trailing char)
	if len(b) < 4 {
		return WriteOpUnknown, nil, nil, fmt.Errorf("too short to be a valid row: %q", b)
	}

	// Find row kind based on first character
	var kind WriteOp
	firstChar := b[0]
	switch firstChar {
	case chars.PutKey:
		kind = WriteOpPutKey
	case chars.PutKeyValue:
		kind = WriteOpPutKeyValue
	case chars.Delete:
		kind = WriteOpDelete
	default:
		return WriteOpUnknown, nil, nil, fmt.Errorf("%w: %q", ErrUnknownWriteOp, firstChar)
	}

	// Fail if second character is not key prefix
	secondChar := b[1]
	if secondChar != chars.KeyPrefix {
		return kind, nil, nil, fmt.Errorf("second char should be %q not %q", chars.KeyPrefix, secondChar)
	}

	// Read key and optional value
	var k []byte
	var v []byte // nil if key-only row (delete or put-key)
	if kind == WriteOpPutKey || kind == WriteOpDelete {
		k = b[2 : len(b)-1] // Read from third char to before last
	} else if kind == WriteOpPutKeyValue {
		valuePrefixOffset := bytes.IndexByte(b[2:], chars.ValuePrefix)
		if valuePrefixOffset == -1 {
			return kind, nil, nil, fmt.Errorf("value prefix %q not found in row: %q", chars.ValuePrefix, b)
		}
		k = b[2 : 2+valuePrefixOffset]          // Read key from third char to value prefix (excl.)
		v = b[2+valuePrefixOffset+1 : len(b)-1] // Read value from prefix (excl.) to before last
	}

	// Check trailing line-break
	lastChar := b[len(b)-1]
	if lastChar != chars.RowEnd {
		return kind, k, v, fmt.Errorf("last char should be %q not %q", chars.RowEnd, lastChar)
	}

	return kind, k, v, nil
}

func extractFileRefs(r io.Reader, format *Format, refs map[string]FileRef) (int, error) {
	offset := 0
	s := bufio.NewScanner(r)
	for s.Scan() {
		row := append(s.Bytes(), '\n')
		writeOp, k, _, err := format.ParseRowFromBytes(row)
		if err != nil {
			return offset, err
		}
		if writeOp == WriteOpDelete {
			delete(refs, string(k))
			continue
		}
		refs[string(k)] = FileRef{Offset: offset, Size: len(row)}
		offset += len(row)
	}
	return offset, s.Err()
}
