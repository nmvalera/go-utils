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

	// 1. Store
	ctx := context.TODO()
	mockStore1.EXPECT().Store(ctx, "test", gomock.Any(), gomock.Any()).Return(nil)
	mockStore2.EXPECT().Store(ctx, "test", gomock.Any(), gomock.Any()).Return(nil)

	err := multiStore.Store(ctx, "test", bytes.NewReader([]byte("test-store")), nil)
	assert.NoError(t, err)

	// 2. Load
	// 2.A mockStore1 returns no error
	mockStore1.EXPECT().Load(ctx, "test", gomock.Any()).Return(io.NopCloser(bytes.NewReader([]byte("test-load"))), nil)

	reader, err := multiStore.Load(context.TODO(), "test", nil)
	assert.NoError(t, err)

	body, err := io.ReadAll(reader)
	assert.NoError(t, err)
	assert.Equal(t, "test-load", string(body))

	// 2.B mockStore1 returns an error
	mockStore1.EXPECT().Load(ctx, "test", gomock.Any()).Return(nil, errors.New("test-error"))
	mockStore2.EXPECT().Load(ctx, "test", gomock.Any()).Return(io.NopCloser(bytes.NewReader([]byte("test-load-2"))), nil)

	reader, err = multiStore.Load(context.TODO(), "test", nil)
	assert.NoError(t, err)

	body, err = io.ReadAll(reader)
	assert.NoError(t, err)
	assert.Equal(t, "test-load-2", string(body))
}
