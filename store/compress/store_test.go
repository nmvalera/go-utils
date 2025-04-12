package compress

import (
	"bytes"
	"context"
	"testing"

	store "github.com/kkrt-labs/go-utils/store"
	"github.com/kkrt-labs/go-utils/store/memory"
	"github.com/kkrt-labs/go-utils/store/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestStore(t *testing.T) {
	tests := []struct {
		desc        string
		encoding    store.ContentEncoding
		key         string
		expectedKey string
		data        []byte
		headers     *store.Headers
		expectErr   bool
	}{
		{
			desc:        "plain",
			encoding:    store.ContentEncodingPlain,
			key:         "test1",
			expectedKey: "test1",
			data:        []byte("message to compress"),
			expectErr:   false,
		},
		{
			desc:        "gzip",
			encoding:    store.ContentEncodingGzip,
			key:         "test2",
			expectedKey: "test2.gz",
			data:        []byte("message to compress"),
			expectErr:   false,
		},
		{
			desc:        "zlib",
			encoding:    store.ContentEncodingZlib,
			key:         "test3",
			expectedKey: "test3.zlib",
			data:        []byte("message to compress"),
			expectErr:   false,
		},
		{
			desc:        "flate",
			encoding:    store.ContentEncodingFlate,
			key:         "test4",
			expectedKey: "test4.flate",
			data:        []byte("message to compress"),
			expectErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			memStore := memory.New()
			mockStore := mock.NewMockStore(ctrl)

			s, err := New(mockStore, WithContentEncoding(tt.encoding))
			require.NoError(t, err)
			ctx := context.TODO()

			mockStore.EXPECT().Store(ctx, tt.expectedKey, gomock.Any(), gomock.Any()).DoAndReturn(memStore.Store)

			// Store the data
			err = s.Store(ctx, tt.key, bytes.NewReader(tt.data), tt.headers)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			// Load the data
			mockStore.EXPECT().Load(ctx, tt.expectedKey).DoAndReturn(memStore.Load)
			reader, _, err := s.Load(ctx, tt.key)
			require.NoError(t, err)

			b := make([]byte, 512)
			n, err := reader.Read(b)
			require.NoError(t, err)
			assert.Equal(t, len(tt.data), n)
			assert.Equal(t, string(tt.data), string(b[:n]))
		})
	}
}
