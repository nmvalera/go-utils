package s3

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/nmvalera/go-utils/aws/mock"
	"github.com/nmvalera/go-utils/common"
	"github.com/nmvalera/go-utils/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestImplementsStore(t *testing.T) {
	assert.Implements(t, (*store.Store)(nil), new(Store))
}

func TestStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockS3Client := mock.NewMockS3ObjectClient(ctrl)

	testBucketName := "test-bucket"
	s3Store, err := New(mockS3Client, testBucketName, WithKeyPrefix("test-prefix"))
	require.NoError(t, err)
	require.NotNil(t, s3Store)

	ctx := context.TODO()

	t.Run("Store", func(t *testing.T) {
		mockS3Client.EXPECT().PutObject(
			ctx,
			gomock.Cond(func(obj *s3.PutObjectInput) bool {
				match := obj.Bucket != nil && *obj.Bucket == testBucketName
				match = match && obj.Key != nil && *obj.Key == "test-prefix/test-key-store"
				if obj.Body == nil {
					return false
				}
				body, err := io.ReadAll(obj.Body)
				if err != nil {
					return false
				}
				match = match && string(body) == "test-data-store"

				match = match && obj.ContentType != nil && *obj.ContentType == store.ContentTypeJSON.String()
				match = match && obj.ContentEncoding != nil && *obj.ContentEncoding == store.ContentEncodingGzip.String()
				match = match && obj.Metadata != nil && obj.Metadata["test-key-store"] == "test-value-store"

				return match
			}),
		).Return(nil, nil)

		err = s3Store.Store(
			ctx,
			"test-key-store",
			strings.NewReader("test-data-store"),
			&store.Headers{
				ContentType:     store.ContentTypeJSON,
				ContentEncoding: store.ContentEncodingGzip,
				KeyValue: map[string]string{
					"test-key-store": "test-value-store",
				},
			},
		)
		assert.NoError(t, err)
	})

	t.Run("Load", func(t *testing.T) {
		mockS3Client.EXPECT().GetObject(
			ctx,
			gomock.Cond(func(obj *s3.GetObjectInput) bool {
				match := obj.Bucket != nil && *obj.Bucket == testBucketName
				match = match && obj.Key != nil && *obj.Key == "test-prefix/test-key-load"
				return match
			}),
		).Return(
			&s3.GetObjectOutput{
				Body:            io.NopCloser(strings.NewReader("test-data-load")),
				ContentType:     common.Ptr(store.ContentTypeJSON.String()),
				ContentEncoding: common.Ptr(store.ContentEncodingGzip.String()),
				Metadata: map[string]string{
					"test-key-load": "test-value-load",
				},
			},
			nil,
		)
		body, headers, err := s3Store.Load(ctx, "test-key-load")
		require.NoError(t, err)
		require.NotNil(t, body)
		require.NotNil(t, headers)
		assert.Equal(t, store.ContentTypeJSON, headers.ContentType)
		assert.Equal(t, store.ContentEncodingGzip, headers.ContentEncoding)
		assert.Equal(t, "test-value-load", headers.KeyValue["test-key-load"])

		b, err := io.ReadAll(body)
		require.NoError(t, err)
		assert.Equal(t, "test-data-load", string(b))
	})

	t.Run("Copy", func(t *testing.T) {
		mockS3Client.EXPECT().CopyObject(
			ctx,
			gomock.Cond(func(obj *s3.CopyObjectInput) bool {
				match := obj.Bucket != nil && *obj.Bucket == testBucketName
				match = match && obj.Key != nil && *obj.Key == "test-prefix/test-key-copy"
				match = match && obj.CopySource != nil && *obj.CopySource == fmt.Sprintf("%s/test-prefix/test-key-copy", testBucketName)
				return match
			}),
		).Return(nil, nil)

		err = s3Store.Copy(ctx, "test-key-copy", "test-key-copy")
		assert.NoError(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		mockS3Client.EXPECT().DeleteObject(
			ctx,
			gomock.Cond(func(obj *s3.DeleteObjectInput) bool {
				match := obj.Bucket != nil && *obj.Bucket == testBucketName
				match = match && obj.Key != nil && *obj.Key == "test-prefix/test-key-delete"
				return match
			}),
		).Return(nil, nil)

		err = s3Store.Delete(ctx, "test-key-delete")
		assert.NoError(t, err)
	})
}
