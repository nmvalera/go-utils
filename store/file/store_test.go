package file

import (
	"bytes"
	"context"
	"io"
	"path/filepath"
	"testing"

	store "github.com/kkrt-labs/go-utils/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileStore(t *testing.T) {
	dataDir := t.TempDir()
	s, err := New(dataDir)
	require.NoError(t, err)

	tests := []struct {
		desc    string
		key     string
		data    string
		headers *store.Headers

		expectedErr  bool
		expectedPath string
	}{
		{
			desc:         "Simple key and no headers",
			key:          "test1",
			data:         "test#1",
			headers:      nil,
			expectedErr:  false,
			expectedPath: "test1",
		},
		{
			desc:         "Key with slash and no headers",
			key:          "test/test2",
			data:         "test#2",
			headers:      nil,
			expectedErr:  false,
			expectedPath: "test/test2",
		},
		{
			desc:         "Key with multiple dots and no headers",
			key:          "test3.txt",
			data:         "test#3",
			headers:      nil,
			expectedErr:  false,
			expectedPath: "test3.txt",
		},
		{
			desc:         "Second store on same key",
			key:          "test1",
			data:         "test#4",
			headers:      nil,
			expectedErr:  false,
			expectedPath: "test1",
		},
		{
			desc: "Simple key with headers content type and encoding",
			key:  "test5",
			data: "test#5",
			headers: &store.Headers{
				ContentType:     store.ContentTypeJSON,
				ContentEncoding: store.ContentEncodingGzip,
			},
			expectedErr:  false,
			expectedPath: "test5.json.gz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			err := s.Store(context.Background(), tt.key, bytes.NewReader([]byte(tt.data)), tt.headers)
			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.FileExists(t, filepath.Join(dataDir, tt.expectedPath))

			reader, err := s.Load(context.Background(), tt.key, tt.headers)
			require.NoError(t, err)

			defer reader.Close()

			content, err := io.ReadAll(reader)
			require.NoError(t, err)
			assert.Equal(t, tt.data, string(content))
		})
	}
}

func TestExtension(t *testing.T) {
	assert.Equal(t, "txt", Extension(store.ContentTypeText, store.ContentEncodingPlain))
	assert.Equal(t, "protobuf", Extension(store.ContentTypeProtobuf, store.ContentEncodingPlain))
	assert.Equal(t, "json", Extension(store.ContentTypeJSON, store.ContentEncodingPlain))
	assert.Equal(t, "txt.gz", Extension(store.ContentTypeText, store.ContentEncodingGzip))
	assert.Equal(t, "txt.zlib", Extension(store.ContentTypeText, store.ContentEncodingZlib))
	assert.Equal(t, "txt.flate", Extension(store.ContentTypeText, store.ContentEncodingFlate))
}
