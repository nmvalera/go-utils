package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
)

func TestS3ClientObject(t *testing.T) {
	assert.Implements(t, (*S3ObjectClient)(nil), new(s3.Client))
}
