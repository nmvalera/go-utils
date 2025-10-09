package jsonrpc_test

import (
	"context"
	"testing"
	"time"

	"github.com/nmvalera/go-utils/jsonrpc"
	jsonrpcmock "github.com/nmvalera/go-utils/jsonrpc/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestWithVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCli := jsonrpcmock.NewMockClient(ctrl)
	c := jsonrpc.WithVersion("2.0")(mockCli)

	mockCli.EXPECT().Call(
		gomock.Any(),
		jsonrpcmock.HasVersion("2.0"),
		gomock.Any(),
	).Return(nil)
	err := c.Call(context.Background(), &jsonrpc.Request{}, nil)
	require.NoError(t, err)
}

func TestWithIncrementalID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCli := jsonrpcmock.NewMockClient(ctrl)
	c := jsonrpc.WithIncrementalID()(mockCli)

	mockCli.EXPECT().Call(
		gomock.Any(),
		jsonrpcmock.HasID(uint32(0)),
		gomock.Any(),
	).Return(nil)
	err := c.Call(context.Background(), &jsonrpc.Request{}, nil)
	require.NoError(t, err)

	mockCli.EXPECT().Call(
		gomock.Any(),
		jsonrpcmock.HasID(uint32(1)),
		gomock.Any())
	err = c.Call(context.Background(), &jsonrpc.Request{}, nil)
	require.NoError(t, err)
}

func TestWithTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCli := jsonrpcmock.NewMockClient(ctrl)
	c := jsonrpc.WithTimeout(100 * time.Millisecond)(mockCli)

	// Does not timeout
	mockCli.EXPECT().Call(gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(nil)
	err := c.Call(context.Background(), &jsonrpc.Request{}, nil)
	require.NoError(t, err)

	// Does timeout
	mockCli.EXPECT().Call(gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).DoAndReturn(func(ctx context.Context, _ *jsonrpc.Request, _ any) error {
		select {
		case <-time.After(1 * time.Second):
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})
	err = c.Call(context.Background(), &jsonrpc.Request{}, nil)
	require.Error(t, err)
	assert.Equal(t, "jsonrpc: call timed out after \"100ms\": context deadline exceeded", err.Error())
}
