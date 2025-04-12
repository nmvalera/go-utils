package s3

import (
	"context"
	"io"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/kkrt-labs/go-utils/common"
	"github.com/kkrt-labs/go-utils/store"
)

type Store struct {
	client *s3.Client

	bucket    string
	keyPrefix string
}

type Options func(*Store) error

func New(s3c *s3.Client, bucket string, opts ...Options) (store.Store, error) {
	s := &Store{
		client: s3c,
		bucket: bucket,
	}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

// Store stores the data in the S3 bucket
func (s *Store) Store(ctx context.Context, key string, reader io.Reader, headers *store.Headers) error {
	input := &s3.PutObjectInput{
		Bucket: &s.bucket,
		Key:    common.Ptr(s.path(key)),
		Body:   reader,
	}

	// Set metadata from headers
	if headers != nil {
		if headers.ContentEncoding != store.ContentEncodingPlain {
			input.ContentEncoding = common.Ptr(headers.ContentEncoding.String())
		}

		if headers.ContentType != store.ContentTypeText {
			input.ContentType = common.Ptr(headers.ContentType.String())
		}

		if headers.KeyValue != nil {
			input.Metadata = headers.KeyValue
		}
	}

	// Store the object
	_, err := s.client.PutObject(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

// Load loads the data from the S3 bucket
// It is the responsibility of the caller to close the returned reader
func (s *Store) Load(ctx context.Context, key string) (io.ReadCloser, *store.Headers, error) {
	output, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    common.Ptr(s.path(key)),
	})
	if err != nil {
		return nil, nil, err
	}

	headers := &store.Headers{}

	if output.ContentType != nil {
		headers.ContentType, _ = store.ParseContentType(*output.ContentType)
	}

	if output.ContentEncoding != nil {
		headers.ContentEncoding, _ = store.ParseContentEncoding(*output.ContentEncoding)
	}

	if output.Metadata != nil {
		headers.KeyValue = output.Metadata
	}

	return output.Body, headers, nil
}

func (s *Store) path(key string) string {
	return filepath.Join(s.keyPrefix, key)
}

// WithKeyPrefix sets the key prefix for the store.
func WithKeyPrefix(prefix string) Options {
	return func(s *Store) error {
		s.keyPrefix = prefix
		return nil
	}
}
