package s3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"github.com/nmvalera/go-utils/aws"
	"github.com/nmvalera/go-utils/common"
	"github.com/nmvalera/go-utils/store"
)

// Store is a store that uses S3 as the underlying storage.
type Store struct {
	client aws.S3ObjectClient

	bucket    string
	keyPrefix string
}

type Options func(*Store) error

func New(s3c aws.S3ObjectClient, bucket string, opts ...Options) (store.Store, error) {
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
		Bucket: common.Ptr(s.bucket),
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
		Bucket: common.Ptr(s.bucket),
		Key:    common.Ptr(s.path(key)),
	})
	if err != nil {
		var aerr smithy.APIError
		if errors.As(err, &aerr) && aerr.ErrorCode() == "NoSuchKey" {
			return nil, nil, store.ErrNotFound
		}
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

// Copy copies an object from one key to another
func (s *Store) Copy(ctx context.Context, srcKey, dstKey string) error {
	_, err := s.client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     common.Ptr(s.bucket),
		Key:        common.Ptr(s.path(dstKey)),
		CopySource: common.Ptr(fmt.Sprintf("%s/%s", s.bucket, s.path(srcKey))),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: common.Ptr(s.bucket),
		Key:    common.Ptr(s.path(key)),
	})
	return err
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
