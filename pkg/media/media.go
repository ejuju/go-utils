package media

import (
	"io"
	"os"
	"path/filepath"
)

type FileStorage interface {
	Store(id string, r io.Reader) error
	Open(id string) (io.ReadCloser, error)
	List() ([]string, error)
}

// Media storage implementation.
// Stores files to a local disk folder.

type LocalFileStorage struct{ dirPath string }

func NewLocalDiskStorage(dirPath string) *LocalFileStorage {
	return &LocalFileStorage{dirPath: dirPath}
}

func (ms *LocalFileStorage) Store(id string, r io.Reader) error {
	f, err := os.Create(ms.getFilePath(id))
	if err != nil {
		return err
	}
	_, err = io.Copy(f, r)
	return err
}

func (fs *LocalFileStorage) Open(id string) (io.ReadCloser, error) {
	return os.Open(fs.getFilePath(id))
}

func (fs *LocalFileStorage) List() ([]string, error) {
	out := []string{}
	return out, filepath.WalkDir(fs.dirPath, func(path string, d os.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		out = append(out, "/"+path)
		return nil
	})
}

func (ms *LocalFileStorage) getFilePath(id string) string { return filepath.Join(ms.dirPath, id) }
