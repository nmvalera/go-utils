package file

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kkrt-labs/go-utils/store"
)

type Store struct {
	dataDir string
}

func New(dataDir string) *Store {
	return &Store{dataDir: dataDir}
}

// Store stores the data in the file
func (f *Store) Store(_ context.Context, key string, reader io.Reader, _ *store.Headers) error {
	filePath := f.filePath(key)
	if err := f.ensureDir(filePath); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (f *Store) ensureDir(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	return nil
}

// Load loads the data from the file
// It is the responsibility of the caller to close the returned reader
func (f *Store) Load(_ context.Context, key string) (io.ReadCloser, *store.Headers, error) {
	filePath := f.filePath(key)
	if !f.fileExists(filePath) {
		return nil, nil, store.ErrNotFound
	}

	o, err := os.Open(filePath)
	return o, nil, err
}

func (f *Store) fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (f *Store) Copy(_ context.Context, srcKey, dstKey string) error {
	srcPath := f.filePath(srcKey)
	if !f.fileExists(srcPath) {
		return store.ErrNotFound
	}

	dstPath := f.filePath(dstKey)
	if err := f.ensureDir(dstPath); err != nil {
		return err
	}

	source, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer source.Close()

	dest, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dest.Close()

	_, err = io.Copy(dest, source)
	return err
}

func (f *Store) Delete(_ context.Context, key string) error {
	filePath := f.filePath(key)
	if !f.fileExists(filePath) {
		return store.ErrNotFound
	}
	return os.Remove(filePath)
}

func (f *Store) filePath(key string) string {
	return filepath.Join(f.dataDir, key)
}
