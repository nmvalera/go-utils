package multistore

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/kkrt-labs/go-utils/store/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// use mock store to test
func TestMultiStoreMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock stores
	mockStore1 := mock.NewMockStore(ctrl)
	mockStore2 := mock.NewMockStore(ctrl)

	// Create a multiStore with the mock stores
	multiStore := New(mockStore1, mockStore2)

	t.Run("Store", func(t *testing.T) {
		ctx := context.TODO()
		mockStore1.EXPECT().Store(ctx, "test", gomock.Any(), gomock.Any()).Return(nil)
		mockStore2.EXPECT().Store(ctx, "test", gomock.Any(), gomock.Any()).Return(nil)

		err := multiStore.Store(ctx, "test", bytes.NewReader([]byte("test-store")), nil)
		assert.NoError(t, err)
	})

	t.Run("Load#Store1 returns", func(t *testing.T) {
		ctx := context.TODO()
		mockStore1.EXPECT().Load(ctx, "test").Return(io.NopCloser(bytes.NewReader([]byte("test-load"))), nil, nil)

		reader, _, err := multiStore.Load(ctx, "test")
		assert.NoError(t, err)

		body, err := io.ReadAll(reader)
		assert.NoError(t, err)
		assert.Equal(t, "test-load", string(body))
	})

	// 2.B mockStore1 returns an error
	t.Run("Load#Store1 returns an error", func(t *testing.T) {
		ctx := context.TODO()
		mockStore1.EXPECT().Load(ctx, "test").Return(nil, nil, errors.New("test-error"))
		mockStore2.EXPECT().Load(ctx, "test").Return(io.NopCloser(bytes.NewReader([]byte("test-load-2"))), nil, nil)

		reader, _, err := multiStore.Load(context.TODO(), "test")
		assert.NoError(t, err)

		body, err := io.ReadAll(reader)
		assert.NoError(t, err)
		assert.Equal(t, "test-load-2", string(body))
	})

	t.Run("Copy", func(t *testing.T) {
		ctx := context.TODO()
		mockStore1.EXPECT().Copy(ctx, "test", "test").Return(nil)
		mockStore2.EXPECT().Copy(ctx, "test", "test").Return(nil)

		err := multiStore.Copy(ctx, "test", "test")
		assert.NoError(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		ctx := context.TODO()
		mockStore1.EXPECT().Delete(ctx, "test").Return(nil)
		mockStore2.EXPECT().Delete(ctx, "test").Return(nil)

		err := multiStore.Delete(ctx, "test")
		assert.NoError(t, err)
	})
}
