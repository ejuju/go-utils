package kv

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type DB struct {
	mu         sync.RWMutex       // Expose mutex: client handles concurrency concerns
	format     *Format            // For encoding and decoding data
	fileRO     *os.File           // Read-only
	fileWO     *os.File           // Write-only
	fileOffset int                // Current file write offset (to record new row positions)
	fileRefs   map[string]FileRef // Maps database rows by key to their offset and size on file
}

// Reference to a specific range of bytes in a file.
// Used to map a row to its location on file.
type FileRef struct {
	Offset int
	Size   int
}

// Instanciates a new DB.
// Opens underlying file handles for the given path (for reading and writing data to a file on disk)
// and extract initial data.
func NewDB(fpath string, chars *Format) (*DB, error) {
	var err error
	db := &DB{format: chars, fileRefs: make(map[string]FileRef)}

	// Open file (we need one read-only and one write-only handle)
	db.fileRO, db.fileWO, err = openFileROWO(fpath)
	if err != nil {
		return nil, err
	}

	// Extract file refs from file and store DB offset
	db.fileOffset, err = extractFileRefs(db.fileRO, db.format, db.fileRefs)
	if err != nil {
		return nil, fmt.Errorf("extract file refs: %w", err)
	}

	return db, nil
}

// Gracefully closes the database.
// Close underyling file handles.
func (db *DB) Close() error {
	if err := db.fileRO.Close(); err != nil {
		return fmt.Errorf("close read-only file: %w", err)
	}
	if err := db.fileWO.Close(); err != nil {
		return fmt.Errorf("close write-only file: %w", err)
	}
	return nil
}

// Get the corresponding offset and size of a row on file.
// If the key file ref is not found, the key does not exist.
func (db *DB) KeyFileRef(k []byte) (FileRef, bool) { ref, ok := db.fileRefs[string(k)]; return ref, ok }

// Reports whether a key is known (has been set before).
func (db *DB) KeyExists(k []byte) bool { _, ok := db.fileRefs[string(k)]; return ok }

// Reports the number of unique keys in the database.
func (db *DB) Count() int { return len(db.fileRefs) }

// Reports the underlying datafile path.
func (db *DB) FilePath() string { return db.fileRO.Name() }

// Calls Sync on the underlying write-only file handle.
func (db *DB) Sync() error { return db.fileWO.Sync() }

// Calls lock on the underlying mutex.
func (db *DB) Lock()   { db.mu.Lock() }
func (db *DB) Unlock() { db.mu.Unlock() }

// Put set a key or key-value pair in the database.
func (db *DB) Put(k, v []byte) error {
	// Encode row bytes
	var b []byte
	var err error
	if v == nil {
		b, err = db.format.Encode(WriteOpPutKey, k, nil)
	} else {
		b, err = db.format.Encode(WriteOpPutKeyValue, k, v)
	}
	if err != nil {
		return fmt.Errorf("encoding: %w", err)
	}
	// Append row to file
	_, err = db.fileWO.Write(b)
	if err != nil {
		return fmt.Errorf("append row to file: %w", err)
	}
	db.fileRefs[string(k)] = FileRef{Offset: db.fileOffset, Size: len(b)}
	db.fileOffset += len(b)
	return err
}

// Removes a key from the database.
// Further attempts to access this key will result in a ErrKeyNotFound.
func (db *DB) Delete(k []byte) error {
	// Encode row bytes
	b, err := db.format.Encode(WriteOpDelete, k, nil)
	if err != nil {
		return fmt.Errorf("encoding: %w", err)
	}
	// Append row to file
	_, err = db.fileWO.Write(b)
	if err != nil {
		return fmt.Errorf("file write: %w", err)
	}
	// Remove file ref
	delete(db.fileRefs, string(k))
	return err
}

// Returns the value associated with the given key.
func (db *DB) Get(k []byte) ([]byte, error) {
	// Get file ref and fail with ErrKeyNotFound if key does not exist
	ref, ok := db.KeyFileRef(k)
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrKeyNotFound, k)
	}

	// Read row bytes from file
	row := make([]byte, ref.Size)
	_, err := db.fileRO.ReadAt(row, int64(ref.Offset))
	if err != nil {
		return nil, fmt.Errorf("read key %q at offset %d with size %d: %w", k, ref.Offset, ref.Size, err)
	}

	// Parse bytes and return extracted value
	_, _, v, err := db.format.ParseRowFromBytes(row)
	if err != nil {
		return nil, fmt.Errorf("parse row: %w", err)
	}

	return v, nil
}

// Iterates over all the keys in the database.
// The callback may return true to exit the loop.
// Keys are not ordered.
func (db *DB) ForEachKey(callback func(k []byte) (stop bool)) {
	for k := range db.fileRefs {
		stop := callback([]byte(k))
		if stop {
			break
		}
	}
}

// Implements the io.WriterTo interface.
// Writes database data to the given writer.
func (db *DB) WriteTo(w io.Writer) (int64, error) {
	n, _, err := db.CompactTo(w)
	return int64(n), err
}

// Writes the database data to the writer, skipping deleted and stale data.
func (db *DB) CompactTo(w io.Writer) (int, map[string]FileRef, error) {
	offset := 0
	newRefs := make(map[string]FileRef, len(db.fileRefs))
	for k, ref := range db.fileRefs {
		row := make([]byte, ref.Size)
		_, err := db.fileRO.ReadAt(row, int64(ref.Offset))
		if err != nil {
			return offset, nil, fmt.Errorf("read row %q: %w", k, err)
		}
		n, err := w.Write(row)
		if err != nil {
			return offset, nil, fmt.Errorf("write row %q: %w", k, err)
		}
		newRefs[k] = FileRef{Offset: offset, Size: n}
		offset += n
	}
	return offset, newRefs, nil
}
