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
	errors := make([]error, 0, len(m.stores))
	for _, s := range m.stores {
		if err := s.Store(ctx, key, reader, headers); err != nil {
			errors = append(errors, err)
		}
	}
	return multierr.Combine(errors...)
}

// Load loads the data from the first store that doesn't return an error
// If all stores return an error, it returns all errors as a multierr error
// It is the responsibility of the caller to close the returned reader
func (m *Store) Load(ctx context.Context, key string) (io.ReadCloser, *store.Headers, error) {
	// Try stores in order until we find the data or encounter an error
	errors := make([]error, 0, len(m.stores))
	for _, s := range m.stores {
		reader, headers, err := s.Load(ctx, key)
		if err == nil && reader != nil {
			return reader, headers, nil
		}

		errors = append(errors, err)
	}

	err := multierr.Combine(errors...)
	if err == nil {
		return nil, nil, store.ErrNotFound
	}

	return nil, nil, err
}

func (m *Store) Copy(ctx context.Context, srcKey, dstKey string) error {
	errors := make([]error, 0, len(m.stores))
	for _, s := range m.stores {
		if err := s.Copy(ctx, srcKey, dstKey); err != nil {
			errors = append(errors, err)
		}
	}
	return multierr.Combine(errors...)
}

func (m *Store) Delete(ctx context.Context, key string) error {
	errors := make([]error, 0, len(m.stores))
	for _, s := range m.stores {
		if err := s.Delete(ctx, key); err != nil {
			errors = append(errors, err)
		}
	}
	return multierr.Combine(errors...)
}
