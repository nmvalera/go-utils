package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtension(t *testing.T) {
	assert.Equal(t, "", Extension(ContentTypeText, ContentEncodingPlain))
	assert.Equal(t, "protobuf", Extension(ContentTypeProtobuf, ContentEncodingPlain))
	assert.Equal(t, "json", Extension(ContentTypeJSON, ContentEncodingPlain))
	assert.Equal(t, "gz", Extension(ContentTypeText, ContentEncodingGzip))
	assert.Equal(t, "json.zlib", Extension(ContentTypeJSON, ContentEncodingZlib))
	assert.Equal(t, "json.flate", Extension(ContentTypeJSON, ContentEncodingFlate))
}
