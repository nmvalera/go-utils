package http

import (
	"context"
	"testing"

	kkrtnet "github.com/kkrt-labs/go-utils/net"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	server := &Server{
		Entrypoint: kkrtnet.NewEntrypoint((&kkrtnet.EntrypointConfig{}).SetDefault()),
		Server:     NewServer((&ServerConfig{}).SetDefault()),
	}
	err := server.Start(context.Background())
	require.NoError(t, err)

	err = server.Stop(context.Background())
	require.NoError(t, err)
}
