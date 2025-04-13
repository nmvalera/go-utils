package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3ObjectClient is a client that can be used to interact with S3 objects.
//
//go:generate mockgen -source s3.go -destination=mock/s3.go -package=mock S3ObjectClient
type S3ObjectClient interface {
	// GetObject gets an object from S3.
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	// PutObject puts an object into S3.
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	// DeleteObject deletes an object from S3.
	DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
	// CopyObject copies an object from one S3 bucket to another.
	CopyObject(ctx context.Context, params *s3.CopyObjectInput, optFns ...func(*s3.Options)) (*s3.CopyObjectOutput, error)
	// ListObjectsV2 lists objects in an S3 bucket.
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}
