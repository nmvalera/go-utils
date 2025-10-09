package store

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Store is an interface for storing and loading objects.
//
//go:generate mockgen -destination=./mock/store.go -package=mock github.com/nmvalera/go-utils/store Store
type Store interface {
	// Store stores an object in the store.
	//
	// The key is the identifier for the object.
	// The reader is the object to store.
	// The headers are optional metadata about the object.
	Store(ctx context.Context, key string, reader io.Reader, headers *Headers) error

	// Load loads an object from the store.
	//
	// The key is the identifier for the object.
	// The headers are optional metadata about the object.
	// It is the responsibility of the caller to close the returned reader
	Load(ctx context.Context, key string) (io.ReadCloser, *Headers, error)

	// Delete deletes an object from the store.
	Delete(ctx context.Context, key string) error

	// Copy copies an object from one store to another.
	Copy(ctx context.Context, srcKey, dstKey string) error
}

var ErrNotFound = errors.New("not found")

// Headers are optional metadata about an object to store/load
type Headers struct {
	// ContentType is the type of the object
	ContentType ContentType

	// ContentEncoding is the compression algorithm used to store the object.
	ContentEncoding ContentEncoding

	// KeyValue is a map of key-value pairs to store/load with the object.
	KeyValue map[string]string
}

func (h *Headers) GetContentType() (string, error) {
	return strings.TrimPrefix(h.ContentType.String(), "application/"), nil
}
func (h *Headers) GetContentEncoding() (ContentEncoding, error) {
	switch h.ContentEncoding {
	case ContentEncodingGzip:
		return ContentEncodingGzip, nil
	case ContentEncodingZlib:
		return ContentEncodingZlib, nil
	case ContentEncodingFlate:
		return ContentEncodingFlate, nil
	case ContentEncodingPlain:
		return ContentEncodingPlain, nil
	}
	return -1, fmt.Errorf("invalid compression: %s", h.ContentEncoding)
}
