package store

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoOpStore(t *testing.T) {
	store := NewNoOpStore()
	require.NotNil(t, store)

	assert.NoError(t, store.Store(context.Background(), "test", strings.NewReader("test"), nil))

	_, _, err := store.Load(context.Background(), "test")
	assert.NoError(t, err)
	assert.NoError(t, store.Delete(context.Background(), "test"))
	assert.NoError(t, store.Copy(context.Background(), "test", "test2"))
}
