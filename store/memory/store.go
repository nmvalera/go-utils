package memory

import (
	"bytes"
	"context"
	"io"

	store "github.com/kkrt-labs/go-utils/store"
)

type Store struct {
	data map[string][]byte
}

func New() store.Store {
	return &Store{
		data: make(map[string][]byte),
	}
}

// Store stores the data in the memory store
func (s *Store) Store(_ context.Context, key string, reader io.Reader, _ *store.Headers) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	s.data[key] = data
	return nil
}

// Load loads the data from the memory store
func (s *Store) Load(_ context.Context, key string) (io.ReadCloser, *store.Headers, error) {
	data, ok := s.data[key]
	if !ok {
		return nil, nil, store.ErrNotFound
	}

	return io.NopCloser(bytes.NewReader(data)), nil, nil
}

func (s *Store) Copy(_ context.Context, srcKey, dstKey string) error {
	data, ok := s.data[srcKey]
	if !ok {
		return store.ErrNotFound
	}

	s.data[dstKey] = data
	return nil
}

func (s *Store) Delete(_ context.Context, key string) error {
	delete(s.data, key)
	return nil
}
