package media

import (
	"io"
	"os"
	"path/filepath"
)

type Storage interface {
	Store(id string, r io.Reader) error
	Open(id string) (io.ReadCloser, error)
	List() ([]string, error)
}

// Media storage implementation.
// Stores files to a local disk folder.

type LocalDiskStorage struct{ dirPath string }

func NewLocalDiskStorage(dirPath string) *LocalDiskStorage {
	return &LocalDiskStorage{dirPath: dirPath}
}

func (ms *LocalDiskStorage) Store(id string, r io.Reader) error {
	f, err := os.Create(ms.getFilePath(id))
	if err != nil {
		return err
	}
	_, err = io.Copy(f, r)
	return err
}

func (fs *LocalDiskStorage) Open(id string) (io.ReadCloser, error) {
	return os.Open(fs.getFilePath(id))
}

func (fs *LocalDiskStorage) List() ([]string, error) {
	out := []string{}
	return out, filepath.WalkDir(fs.dirPath, func(path string, d os.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		out = append(out, "/"+path)
		return nil
	})
}

func (ms *LocalDiskStorage) getFilePath(id string) string { return filepath.Join(ms.dirPath, id) }
