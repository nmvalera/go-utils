package http

import (
	"testing"
	"time"

	"github.com/nmvalera/go-utils/common"
	"github.com/nmvalera/go-utils/config"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntrypointConfig(t *testing.T) {
	v := config.NewViper()
	v.Set("ep.addr", "localhost:8080")
	v.Set("ep.http.readTimeout", "1s")
	v.Set("ep.http.readHeaderTimeout", "2s")
	v.Set("ep.http.writeTimeout", "3s")
	v.Set("ep.http.idleTimeout", "4s")
	v.Set("ep.http.maxHeaderBytes", 1024)
	v.Set("ep.net.keepAlive", "-10s")
	v.Set("ep.net.keepAliveProbe.enable", true)
	v.Set("ep.net.keepAliveProbe.idle", "5s")
	v.Set("ep.net.keepAliveProbe.interval", "6s")
	v.Set("ep.net.keepAliveProbe.count", 3)
	v.Set("ep.tls.certFile", "test-cert.pem")
	v.Set("ep.tls.keyFile", "test-key.pem")

	cfg := new(EntrypointConfig)
	err := cfg.Unmarshal(v)
	require.NoError(t, err)

	expectedCfg := &EntrypointConfig{
		Addr: common.Ptr("localhost:8080"),
		HTTP: &ServerConfig{
			ReadTimeout:       common.Ptr(1 * time.Second),
			ReadHeaderTimeout: common.Ptr(2 * time.Second),
			WriteTimeout:      common.Ptr(3 * time.Second),
			IdleTimeout:       common.Ptr(4 * time.Second),
			MaxHeaderBytes:    common.Ptr(1024),
		},
		Net: &ListenConfig{
			KeepAlive: common.Ptr(-10 * time.Second),
			KeepAliveProbe: &KeepAliveProbeConfig{
				Enable:   common.Ptr(true),
				Idle:     common.Ptr(5 * time.Second),
				Interval: common.Ptr(6 * time.Second),
				Count:    common.Ptr(3),
			},
		},
		TLS: &TLSCertConfig{
			CertFile: common.Ptr("test-cert.pem"),
			KeyFile:  common.Ptr("test-key.pem"),
		},
	}
	assert.Equal(t, expectedCfg, cfg)
}

func TestEnv(t *testing.T) {
	env, err := (&EntrypointConfig{
		Addr: common.Ptr("localhost:8080"),
		HTTP: &ServerConfig{
			ReadTimeout:       common.Ptr(45 * time.Second),
			ReadHeaderTimeout: common.Ptr(46 * time.Second),
			WriteTimeout:      common.Ptr(47 * time.Second),
			IdleTimeout:       common.Ptr(48 * time.Second),
			MaxHeaderBytes:    common.Ptr(1000),
		},
		Net: &ListenConfig{
			KeepAlive: common.Ptr(-10 * time.Second),
			KeepAliveProbe: &KeepAliveProbeConfig{
				Enable:   common.Ptr(true),
				Idle:     common.Ptr(5 * time.Second),
				Interval: common.Ptr(6 * time.Second),
				Count:    common.Ptr(3),
			},
		},
		TLS: &TLSCertConfig{
			CertFile: common.Ptr("test-cert.pem"),
			KeyFile:  common.Ptr("test-key.pem"),
		},
	}).Env()
	require.NoError(t, err)
	expected := map[string]string{
		"EP_ADDR":                          "localhost:8080",
		"EP_HTTP_READ_TIMEOUT":             "45s",
		"EP_HTTP_READ_HEADER_TIMEOUT":      "46s",
		"EP_HTTP_WRITE_TIMEOUT":            "47s",
		"EP_HTTP_IDLE_TIMEOUT":             "48s",
		"EP_HTTP_MAX_HEADER_BYTES":         "1000",
		"EP_NET_KEEP_ALIVE":                "-10s",
		"EP_NET_KEEP_ALIVE_PROBE_ENABLE":   "true",
		"EP_NET_KEEP_ALIVE_PROBE_IDLE":     "5s",
		"EP_NET_KEEP_ALIVE_PROBE_INTERVAL": "6s",
		"EP_NET_KEEP_ALIVE_PROBE_COUNT":    "3",
		"EP_TLS_CERT_FILE":                 "test-cert.pem",
		"EP_TLS_KEY_FILE":                  "test-key.pem",
	}
	assert.Equal(t, expected, env)
}

func TestAddFlagsAndLoadEnv(t *testing.T) {
	cfg := &EntrypointConfig{
		Addr: common.Ptr("localhost:8080"),
		HTTP: &ServerConfig{
			ReadTimeout:       common.Ptr(45 * time.Second),
			ReadHeaderTimeout: common.Ptr(46 * time.Second),
			WriteTimeout:      common.Ptr(47 * time.Second),
			IdleTimeout:       common.Ptr(48 * time.Second),
			MaxHeaderBytes:    common.Ptr(1000),
		},
		Net: &ListenConfig{
			KeepAlive: common.Ptr(-10 * time.Second),
			KeepAliveProbe: &KeepAliveProbeConfig{
				Enable:   common.Ptr(true),
				Idle:     common.Ptr(5 * time.Second),
				Interval: common.Ptr(6 * time.Second),
				Count:    common.Ptr(3),
			},
		},
		TLS: &TLSCertConfig{
			CertFile: common.Ptr("test-cert.pem"),
			KeyFile:  common.Ptr("test-key.pem"),
		},
	}

	v := config.NewViper()
	set := pflag.NewFlagSet("test", pflag.ContinueOnError)
	set.SortFlags = true
	err := AddFlags(v, set)
	require.NoError(t, err)

	expectedUsage := "      --ep-addr string                            TCP Address to listen on [env: EP_ADDR]\n      --ep-http-idle-timeout string               Maximum duration to wait for the next request when keep-alives are enabled (zero uses the value of read timeout) [env: EP_HTTP_IDLE_TIMEOUT] (default \"30s\")\n      --ep-http-max-header-bytes int              Maximum number of bytes the server will read parsing the request header's keys and values [env: EP_HTTP_MAX_HEADER_BYTES] (default 1048576)\n      --ep-http-read-header-timeout string        Maximum duration for reading request headers (zero uses the value of read timeout) [env: EP_HTTP_READ_HEADER_TIMEOUT] (default \"30s\")\n      --ep-http-read-timeout string               Maximum duration for reading the entire request including the body (zero means no timeout) [env: EP_HTTP_READ_TIMEOUT] (default \"30s\")\n      --ep-http-write-timeout string              Maximum duration before timing out writes of the response (zero means no timeout) [env: EP_HTTP_WRITE_TIMEOUT] (default \"30s\")\n      --ep-net-keep-alive string                  Keep alive period for network connections accepted by this entrypoint [env: EP_NET_KEEP_ALIVE] (default \"-1s\")\n      --ep-net-keep-alive-probe-count int         Maximum number of keep-alive probes that can go unanswered before dropping a connection [env: EP_NET_KEEP_ALIVE_PROBE_COUNT] (default 9)\n      --ep-net-keep-alive-probe-enable            Enable keep alive probes [env: EP_NET_KEEP_ALIVE_PROBE_ENABLE]\n      --ep-net-keep-alive-probe-idle string       Time that the connection must be idle before the first keep-alive probe is sent [env: EP_NET_KEEP_ALIVE_PROBE_IDLE] (default \"15s\")\n      --ep-net-keep-alive-probe-interval string   Time between keep-alive probes [env: EP_NET_KEEP_ALIVE_PROBE_INTERVAL] (default \"15s\")\n      --ep-tls-certfile string                    Path to the certificate file [env: EP_TLS_CERT_FILE]\n      --ep-tls-keyfile string                     Path to the key file [env: EP_TLS_KEY_FILE]\n"
	assert.Equal(t, expectedUsage, set.FlagUsages())

	// Generate the environment variables
	env, err := cfg.Env()
	require.NoError(t, err)
	for k, v := range env {
		t.Setenv(k, v)
	}

	loadedCfg := new(EntrypointConfig)
	err = loadedCfg.Unmarshal(v)
	require.NoError(t, err)

	assert.Equal(t, cfg, loadedCfg)
}

func TestUnmarshalFromDefaults(t *testing.T) {
	v := config.NewViper()
	err := AddFlags(v, pflag.NewFlagSet("test", pflag.ContinueOnError))
	require.NoError(t, err)

	cfg := new(EntrypointConfig)
	err = cfg.Unmarshal(v)
	require.NoError(t, err)

	assert.Equal(t, DefaultEntrypointConfig(), cfg)
}
