package compress

import (
	"bytes"
	"context"
	"testing"

	store "github.com/kkrt-labs/go-utils/store"
	"github.com/kkrt-labs/go-utils/store/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	memStore := memory.New()

	tests := []struct {
		desc      string
		encoding  store.ContentEncoding
		key       string
		data      []byte
		headers   *store.Headers
		expectErr bool
	}{
		{
			desc:      "plain",
			encoding:  store.ContentEncodingPlain,
			key:       "plain",
			data:      []byte("message to compress"),
			expectErr: false,
		},
		{
			desc:      "gzip",
			encoding:  store.ContentEncodingGzip,
			key:       "gzip",
			data:      []byte("message to compress"),
			expectErr: false,
		},
		{
			desc:      "zlib",
			encoding:  store.ContentEncodingZlib,
			key:       "zlib",
			data:      []byte("message to compress"),
			expectErr: false,
		},
		{
			desc:     "flate",
			encoding: store.ContentEncodingFlate,
			key:      "flate",
			data:     []byte("message to compress"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			s, err := New(memStore, WithContentEncoding(tt.encoding))
			require.NoError(t, err)
			ctx := context.TODO()

			// Store the data
			err = s.Store(ctx, tt.key, bytes.NewReader(tt.data), tt.headers)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			// Load the data
			reader, err := s.Load(ctx, tt.key, tt.headers)
			require.NoError(t, err)

			b := make([]byte, 512)
			n, err := reader.Read(b)
			require.NoError(t, err)
			assert.Equal(t, len(tt.data), n)
			assert.Equal(t, string(tt.data), string(b[:n]))
		})
	}
}
