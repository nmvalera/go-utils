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
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
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

// Load loads the data from the file
// It is the responsibility of the caller to close the returned reader
func (f *Store) Load(_ context.Context, key string) (io.ReadCloser, *store.Headers, error) {
	o, err := os.Open(f.filePath(key))
	return o, nil, err
}

func (f *Store) filePath(key string) string {
	return filepath.Join(f.dataDir, key)
}
