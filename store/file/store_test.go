package file

import (
	"bytes"
	"context"
	"io"
	"path/filepath"
	"testing"

	"github.com/kkrt-labs/go-utils/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImplementsStore(t *testing.T) {
	assert.Implements(t, (*store.Store)(nil), new(Store))
}

func TestFileStore(t *testing.T) {
	dataDir := t.TempDir()
	s := New(dataDir)

	tests := []struct {
		desc    string
		key     string
		keyCopy string
		data    string

		expectedErr      bool
		expectedPath     string
		expectedPathCopy string
	}{
		{
			desc:         "Simple key and no headers",
			key:          "test1",
			data:         "test#1",
			expectedErr:  false,
			expectedPath: "test1",
		},
		{
			desc:         "Key with slash and no headers",
			key:          "test/test2",
			data:         "test#2",
			expectedErr:  false,
			expectedPath: "test/test2",
		},
		{
			desc:         "Key with multiple dots and no headers",
			key:          "test3.txt",
			data:         "test#3",
			expectedErr:  false,
			expectedPath: "test3.txt",
		},
		{
			desc:         "Second store on same key",
			key:          "test1",
			data:         "test#4",
			expectedErr:  false,
			expectedPath: "test1",
		},
		{
			desc:         "Simple key with headers content type and encoding",
			key:          "test5.json.gz",
			data:         "test#5",
			expectedErr:  false,
			expectedPath: "test5.json.gz",
		},
		{
			desc:             "Copy",
			key:              "test1",
			keyCopy:          "test1-copy",
			data:             "test#1",
			expectedErr:      false,
			expectedPath:     "test1",
			expectedPathCopy: "test1-copy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			err := s.Store(context.Background(), tt.key, bytes.NewReader([]byte(tt.data)), nil)
			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.FileExists(t, filepath.Join(dataDir, tt.expectedPath))

			reader, _, err := s.Load(context.Background(), tt.key)
			require.NoError(t, err)

			defer reader.Close()

			content, err := io.ReadAll(reader)
			require.NoError(t, err)
			assert.Equal(t, tt.data, string(content))

			if tt.expectedPathCopy != "" {
				err = s.Copy(context.Background(), tt.key, tt.keyCopy)
				require.NoError(t, err)

				assert.FileExists(t, filepath.Join(dataDir, tt.expectedPathCopy))

				reader, _, err := s.Load(context.Background(), tt.keyCopy)
				require.NoError(t, err)

				defer reader.Close()

				content, err := io.ReadAll(reader)
				require.NoError(t, err)
				assert.Equal(t, tt.data, string(content))

				err = s.Delete(context.Background(), tt.keyCopy)
				require.NoError(t, err)

				assert.NoFileExists(t, filepath.Join(dataDir, tt.expectedPathCopy))
			}

			err = s.Delete(context.Background(), tt.key)
			require.NoError(t, err)

			assert.NoFileExists(t, filepath.Join(dataDir, tt.expectedPath))
		})
	}
}
