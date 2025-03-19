package compress

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"context"
	"fmt"
	"io"

	store "github.com/kkrt-labs/go-utils/store"
)

type Store struct {
	store           store.Store
	contentEncoding store.ContentEncoding
}

type Options func(*Store) error

func New(s store.Store, opts ...Options) (*Store, error) {
	cs := &Store{
		store:           s,
		contentEncoding: store.ContentEncodingPlain,
	}

	for _, opt := range opts {
		if err := opt(cs); err != nil {
			return nil, err
		}
	}

	return cs, nil
}

// Store stores the data in the store
func (c *Store) Store(ctx context.Context, key string, reader io.Reader, headers *store.Headers) error {
	var compressedReader io.Reader

	switch c.contentEncoding {
	case store.ContentEncodingPlain:
		compressedReader = reader
	case store.ContentEncodingGzip:
		buf := new(bytes.Buffer)
		gw := gzip.NewWriter(buf)
		defer gw.Close()
		if _, err := io.Copy(gw, reader); err != nil {
			return fmt.Errorf("failed to compress with gzip: %w", err)
		}

		if err := gw.Flush(); err != nil {
			return fmt.Errorf("failed to compress with gzip: %w", err)
		}

		compressedReader = buf

	case store.ContentEncodingZlib:
		buf := new(bytes.Buffer)
		zw := zlib.NewWriter(buf)
		defer zw.Close()
		if _, err := io.Copy(zw, reader); err != nil {
			return fmt.Errorf("failed to compress with zlib: %w", err)
		}

		if err := zw.Flush(); err != nil {
			return fmt.Errorf("failed to compress with zlib: %w", err)
		}

		compressedReader = buf

	case store.ContentEncodingFlate:
		buf := new(bytes.Buffer)
		fw, err := flate.NewWriter(buf, flate.BestCompression)
		if err != nil {
			return fmt.Errorf("failed to create flate writer: %w", err)
		}
		defer fw.Close()
		if _, err := io.Copy(fw, reader); err != nil {
			return fmt.Errorf("failed to compress with flate: %w", err)
		}

		if err := fw.Flush(); err != nil {
			return fmt.Errorf("failed to compress with flate: %w", err)
		}

		compressedReader = buf
	default:
		return fmt.Errorf("unsupported content encoding: %s", c.contentEncoding)
	}

	if headers == nil {
		headers = &store.Headers{}
	}
	headers.ContentEncoding = c.contentEncoding

	return c.store.Store(ctx, key, compressedReader, headers)
}

// Load loads the data from the store
// It is the responsibility of the caller to close the returned reader
func (c *Store) Load(ctx context.Context, key string, headers *store.Headers) (io.ReadCloser, error) {
	if headers == nil {
		headers = &store.Headers{}
	}
	headers.ContentEncoding = c.contentEncoding

	reader, err := c.store.Load(ctx, key, headers)
	if err != nil {
		return nil, err
	}

	switch c.contentEncoding {
	case store.ContentEncodingPlain:
		return reader, nil
	case store.ContentEncodingGzip:
		r, err := gzip.NewReader(reader)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress with gzip: %w", err)
		}
		return r, nil
	case store.ContentEncodingZlib:
		r, err := zlib.NewReader(reader)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress with zlib: %w", err)
		}
		return r, nil
	case store.ContentEncodingFlate:
		return flate.NewReader(reader), nil
	default:
		return nil, fmt.Errorf("unsupported content encoding: %s", c.contentEncoding)
	}
}

func WithContentEncoding(encoding store.ContentEncoding) Options {
	return func(s *Store) error {
		s.contentEncoding = encoding
		return nil
	}
}
