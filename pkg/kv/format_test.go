package kv

import (
	"bytes"
	"errors"
	"testing"
)

func TestFormat(t *testing.T) {

	t.Run("encode row", func(t *testing.T) {
		tests := []struct {
			description  string
			inputWriteOp WriteOp
			inputK       []byte
			inputV       []byte
			wantRow      string
			wantErr      func(error) bool
		}{
			{
				description:  "valid put key",
				inputWriteOp: WriteOpPutKey,
				inputK:       []byte("MyKey"),
				inputV:       nil,
				wantRow:      "- MyKey\n",
			},
			{
				description:  "valid put key-value",
				inputWriteOp: WriteOpPutKeyValue,
				inputK:       []byte("MyKey"),
				inputV:       []byte("MyValue"),
				wantRow:      "= MyKey MyValue\n",
			},
			{
				description:  "valid delete",
				inputWriteOp: WriteOpDelete,
				inputK:       []byte("MyKey"),
				inputV:       nil,
				wantRow:      "! MyKey\n",
			},
			{
				description:  "valid put key-value with zero-length value",
				inputWriteOp: WriteOpPutKeyValue,
				inputK:       []byte("MyKey"),
				inputV:       []byte(""),
				wantRow:      "= MyKey \n", // note the trailing whitespace after the key
			},
			{
				description:  "valid put key-value with nil value",
				inputWriteOp: WriteOpPutKeyValue,
				inputK:       []byte("MyKey"),
				inputV:       nil,
				wantRow:      "= MyKey \n", // note the trailing whitespace after the key
			},
			{
				description:  "should fail with ErrKeyEmpty on put with zero-length key",
				inputWriteOp: WriteOpPutKeyValue,
				inputK:       []byte(""),
				inputV:       nil,
				wantRow:      "",
				wantErr:      func(err error) bool { return errors.Is(err, ErrKeyEmpty) },
			},
			{
				description:  "should fail with ErrKeyEmpty on put with nil key",
				inputWriteOp: WriteOpPutKeyValue,
				inputK:       nil,
				inputV:       nil,
				wantErr:      func(err error) bool { return errors.Is(err, ErrKeyEmpty) },
			},
			{
				description:  "should fail with ErrKeyEmpty on delete with zero-length key",
				inputWriteOp: WriteOpDelete,
				inputK:       nil,
				inputV:       []byte(""),
				wantErr:      func(err error) bool { return errors.Is(err, ErrKeyEmpty) },
			},
			{
				description:  "should fail with ErrKeyEmpty on delete with nil key",
				inputWriteOp: WriteOpDelete,
				inputK:       nil,
				wantErr:      func(err error) bool { return errors.Is(err, ErrKeyEmpty) },
			},
			{
				description:  "should fail on put key if key contains row end",
				inputWriteOp: WriteOpPutKey,
				inputK:       []byte("My\nKey"), // note the line-break
				wantErr:      func(err error) bool { return err != nil },
			},
			{
				description:  "should fail on delete key if key contains row end",
				inputWriteOp: WriteOpDelete,
				inputK:       []byte("My\nKey"), // note the line-break
				wantErr:      func(err error) bool { return err != nil },
			},
			{
				description:  "should fail on put key-value if key contains value prefix",
				inputWriteOp: WriteOpPutKeyValue,
				inputK:       []byte("My Key"), // note the whitespace
				inputV:       []byte(""),
				wantErr:      func(err error) bool { return err != nil },
			},
		}

		for _, test := range tests {
			t.Run(test.description, func(t *testing.T) {
				row, err := DefaultFormat.Encode(test.inputWriteOp, test.inputK, test.inputV)
				wantErr := test.wantErr
				if wantErr == nil {
					wantErr = func(err error) bool { return err == nil }
				}
				if !wantErr(err) {
					t.Fatal("unexpected error:", err)
				}
				if string(row) != test.wantRow {
					t.Fatalf("want %q but got %q", test.wantRow, row)
				}
			})
		}
	})

	t.Run("decode row", func(t *testing.T) {
		tests := []struct {
			description string
			inputRow    []byte
			wantWriteOp WriteOp
			wantK       []byte
			wantV       []byte
		}{
			{
				description: "valid put key",
				inputRow:    []byte("- MyKey\n"),
				wantWriteOp: WriteOpPutKey,
				wantK:       []byte("MyKey"),
				wantV:       nil,
			},
			{
				description: "valid put key-value",
				inputRow:    []byte("= MyKey MyValue\n"),
				wantWriteOp: WriteOpPutKeyValue,
				wantK:       []byte("MyKey"),
				wantV:       []byte("MyValue"),
			},
			{
				description: "valid delete",
				inputRow:    []byte("! MyKey\n"),
				wantWriteOp: WriteOpDelete,
				wantK:       []byte("MyKey"),
				wantV:       nil,
			},
			{
				description: "valid put key-value with empty value",
				inputRow:    []byte("= MyKey \n"), // not the trailing whitespace
				wantWriteOp: WriteOpPutKeyValue,
				wantK:       []byte("MyKey"),
				wantV:       []byte(""),
			},
		}

		for _, test := range tests {
			t.Run(test.description, func(t *testing.T) {
				writeOp, k, v, err := DefaultFormat.ParseRowFromBytes(test.inputRow)
				if err != nil {
					t.Fatal(err)
				}
				if writeOp != test.wantWriteOp {
					t.Fatalf("want write-op %q not %q", test.wantWriteOp, writeOp)
				}
				if !bytes.Equal(test.wantK, k) {
					t.Fatalf("want key %q not %q", test.wantK, k)
				}
				if test.wantV == nil && v != nil {
					t.Fatalf("want nil value not %q", v)
				}
				if test.wantV != nil && !bytes.Equal(test.wantV, v) {
					t.Fatalf("want value %q but got %q", test.wantK, k)
				}
			})
		}
	})
}
