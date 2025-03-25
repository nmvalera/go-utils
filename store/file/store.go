package file

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/kkrt-labs/go-utils/store"
)

type Store struct {
	dataDir string
}

func New(dataDir string) *Store {
	return &Store{dataDir: dataDir}
}

// Store stores the data in the file
func (f *Store) Store(_ context.Context, key string, reader io.Reader, headers *store.Headers) error {
	filePath := f.filePath(key, headers)
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
func (f *Store) Load(_ context.Context, key string, headers *store.Headers) (io.ReadCloser, error) {
	return os.Open(f.filePath(key, headers))
}

func (f *Store) filePath(key string, headers *store.Headers) string {
	filePath := filepath.Join(f.dataDir, key)
	var (
		ct store.ContentType
		ce store.ContentEncoding
	)

	if headers != nil {
		ct = headers.ContentType
		ce = headers.ContentEncoding
	}
	ext := Extension(ct, ce)
	if ext != "" {
		filePath = fmt.Sprintf("%s.%s", filePath, ext)
	}

	return filePath
}

func Extension(ct store.ContentType, ce store.ContentEncoding) string {
	var parts []string

	switch ct {
	case store.ContentTypeText:
		parts = append(parts, "txt")
	case store.ContentTypeProtobuf:
		parts = append(parts, "protobuf")
	case store.ContentTypeJSON:
		parts = append(parts, "json")
	}

	switch ce {
	case store.ContentEncodingGzip:
		parts = append(parts, "gz")
	case store.ContentEncodingZlib:
		parts = append(parts, "zlib")
	case store.ContentEncodingFlate:
		parts = append(parts, "flate")
	}

	return strings.Join(parts, ".")
}
