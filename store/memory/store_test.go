package memory

import (
	"bytes"
	"context"
	"io"
	"testing"

	store "github.com/nmvalera/go-utils/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImplementsStore(t *testing.T) {
	assert.Implements(t, (*store.Store)(nil), new(Store))
}

func TestStoreAndLoad(t *testing.T) {
	s := New()

	err := s.Store(context.Background(), "test", bytes.NewReader([]byte("test")), nil)
	require.NoError(t, err)

	reader, _, err := s.Load(context.Background(), "test")
	require.NoError(t, err)

	data, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Equal(t, []byte("test"), data)
}
