package memory

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreAndLoad(t *testing.T) {
	store := New()

	err := store.Store(context.Background(), "test", bytes.NewReader([]byte("test")), nil)
	require.NoError(t, err)

	reader, _, err := store.Load(context.Background(), "test")
	require.NoError(t, err)

	data, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Equal(t, []byte("test"), data)
}
