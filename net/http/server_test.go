package http

import (
	"context"
	"net"
	"net/http"
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

	addr := server.Addr()
	require.NotEmpty(t, addr)

	// Test that the server accepts connections
	conn, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	err = conn.Close()
	require.NoError(t, err)

	// Test that the server accept HTTP requests
	server.Server.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	req, err := http.NewRequest("GET", "http://"+addr, http.NoBody)
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	err = server.Stop(context.Background())
	require.NoError(t, err)
}
