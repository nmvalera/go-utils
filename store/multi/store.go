package multistore

import (
	"context"
	"io"

	store "github.com/kkrt-labs/go-utils/store"
	"go.uber.org/multierr"
)

type Store struct {
	stores []store.Store
}

func New(stores ...store.Store) store.Store {
	return &Store{stores: stores}
}

// Store stores the data in all stores (if one store returns an error, it will continue to the next store)
func (m *Store) Store(ctx context.Context, key string, reader io.Reader, headers *store.Headers) error {
	for _, s := range m.stores {
		if err := s.Store(ctx, key, reader, headers); err != nil {
			return err
		}
	}
	return nil
}

// Load loads the data from the first store that doesn't return an error
// If all stores return an error, it returns all errors as a multierr error
// It is the responsibility of the caller to close the returned reader
func (m *Store) Load(ctx context.Context, key string, headers *store.Headers) (io.ReadCloser, error) {
	// Try stores in order until we find the data or encounter an error
	errors := make([]error, 0, len(m.stores))
	for _, s := range m.stores {
		reader, err := s.Load(ctx, key, headers)
		if err == nil {
			return reader, nil
		}
		errors = append(errors, err)
	}
	return nil, multierr.Combine(errors...)
}
