package http

import (
	"context"
	"net"
	"net/http"
	"testing"

	"github.com/nmvalera/go-utils/common"
	"github.com/nmvalera/go-utils/tag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntrypoint(t *testing.T) {
	ep, err := NewEntrypoint("")
	require.NoError(t, err)
	err = ep.Start(context.Background())
	require.NoError(t, err)

	addr := ep.Addr()
	require.NotEmpty(t, addr)

	// Test that the entrypoint accepts connections
	conn, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	err = conn.Close()
	require.NoError(t, err)

	// Test that the entrypoint accept HTTP requests
	ep.SetHandler(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req, err := http.NewRequest("GET", "http://"+addr, http.NoBody)
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	err = ep.Stop(context.Background())
	require.NoError(t, err)
}

func TestOptions(t *testing.T) {
	t.Run("WithServer", func(t *testing.T) {
		srv := &http.Server{}
		ep, err := NewEntrypoint("", WithServer(srv))
		require.NoError(t, err)
		assert.Equal(t, srv, ep.server)
	})

	t.Run("WithListenConfig", func(t *testing.T) {
		lCfg := &net.ListenConfig{}
		ep, err := NewEntrypoint("", WithListenConfig(lCfg))
		require.NoError(t, err)
		assert.Equal(t, lCfg, ep.lCfg)
	})

	t.Run("WithTags", func(t *testing.T) {
		ep, err := NewEntrypoint("", WithTags(tag.Key("test-key").String("test-value")))
		require.NoError(t, err)
		require.Len(t, ep.tagged.Set, 1)
		assert.Equal(t, tag.Key("test-key").String("test-value"), ep.tagged.Set[0])
	})

	t.Run("WithTLSConfig", func(t *testing.T) {
		tlsCfg := &TLSCertConfig{
			CertFile: common.Ptr("test-cert.pem"),
			KeyFile:  common.Ptr("test-key.pem"),
		}
		ep, err := NewEntrypoint("", WithTLSConfig(tlsCfg))
		require.NoError(t, err)
		assert.Equal(t, tlsCfg, ep.tlsCfg)
	})
}
