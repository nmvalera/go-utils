package net

import (
	"crypto/tls"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func TestEntrypoint(t *testing.T) {
	t.Run("NoTLS", func(t *testing.T) {
		ep := NewEntrypoint(&EntrypointConfig{
			Network: "tcp",
		})

		l, err := ep.Listen(context.TODO())
		require.NoError(t, err)
		require.NotNil(t, l)

		addr := l.Addr()
		require.NotNil(t, addr)
		require.Equal(t, "tcp", addr.Network())

		err = l.Close()
		require.NoError(t, err)
	})

	t.Run("TLS", func(t *testing.T) {
		ep := NewEntrypoint(&EntrypointConfig{
			Network: "tcp",
			TLS: &tls.Config{
				Certificates: []tls.Certificate{
					{
						Certificate: [][]byte{[]byte("cert")},
						PrivateKey:  [][]byte{[]byte("key")},
					},
				},
			},
		})

		l, err := ep.Listen(context.TODO())
		require.NoError(t, err)
		require.NotNil(t, l)

		addr := l.Addr()
		require.NotNil(t, addr)
		require.Equal(t, "tcp", addr.Network())

		err = l.Close()
		require.NoError(t, err)
	})
}
