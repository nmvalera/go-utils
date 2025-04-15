package store

import (
	"context"
	"io"
)

type noOpStore struct{}

// NewNoOpStore returns a new no-op store.
func NewNoOpStore() Store {
	return &noOpStore{}
}

func (s *noOpStore) Store(_ context.Context, _ string, _ io.Reader, _ *Headers) error {
	return nil
}
func (s *noOpStore) Load(_ context.Context, _ string) (io.ReadCloser, *Headers, error) {
	return nil, nil, nil
}
func (s *noOpStore) Delete(_ context.Context, _ string) error  { return nil }
func (s *noOpStore) Copy(_ context.Context, _, _ string) error { return nil }
