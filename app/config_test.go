package app

import (
	"testing"
	"time"

	"github.com/kkrt-labs/go-utils/common"
	"github.com/kkrt-labs/go-utils/config"
	"github.com/kkrt-labs/go-utils/log"
	kkrthttp "github.com/kkrt-labs/go-utils/net/http"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestViperConfig(t *testing.T) {
	v := config.NewViper()
	v.Set("main-ep.addr", "localhost:8888")
	v.Set("main-ep.http.read-timeout", "40s")
	v.Set("main-ep.http.read-header-timeout", "41s")
	v.Set("main-ep.http.write-timeout", "42s")
	v.Set("main-ep.http.idle-timeout", "43s")
	v.Set("main-ep.net.keep-alive", "44s")
	v.Set("main-ep.net.keep-alive-probe.enable", "true")
	v.Set("main-ep.net.keep-alive-probe.idle", "45s")
	v.Set("main-ep.net.keep-alive-probe.interval", "46s")
	v.Set("main-ep.net.keep-alive-probe.count", "47")
	v.Set("main-ep.http.max-header-bytes", "40000")
	v.Set("healthz-ep.addr", "localhost:8889")
	v.Set("healthz-ep.http.read-timeout", "50s")
	v.Set("healthz-ep.http.read-header-timeout", "51s")
	v.Set("healthz-ep.http.write-timeout", "52s")
	v.Set("healthz-ep.http.idle-timeout", "53s")
	v.Set("healthz-ep.net.keep-alive", "54s")
	v.Set("healthz-ep.net.keep-alive-probe.enable", "true")
	v.Set("healthz-ep.net.keep-alive-probe.idle", "55s")
	v.Set("healthz-ep.net.keep-alive-probe.interval", "56s")
	v.Set("healthz-ep.net.keep-alive-probe.count", "57")
	v.Set("healthz-ep.http.max-header-bytes", "50000")
	v.Set("log.level", "info")
	v.Set("start-timeout", "10s")
	v.Set("stop-timeout", "20s")

	cfg := new(Config)
	err := cfg.Unmarshal(v)
	require.NoError(t, err)

	expectedCfg := &Config{
		MainEntrypoint: &kkrthttp.EntrypointConfig{
			Addr: common.Ptr("localhost:8888"),
			HTTP: &kkrthttp.ServerConfig{
				ReadTimeout:       common.Ptr(40 * time.Second),
				ReadHeaderTimeout: common.Ptr(41 * time.Second),
				WriteTimeout:      common.Ptr(42 * time.Second),
				IdleTimeout:       common.Ptr(43 * time.Second),
				MaxHeaderBytes:    common.Ptr(40000),
			},
			Net: &kkrthttp.ListenConfig{
				KeepAlive: common.Ptr(44 * time.Second),
				KeepAliveProbe: &kkrthttp.KeepAliveProbeConfig{
					Enable:   common.Ptr(true),
					Idle:     common.Ptr(45 * time.Second),
					Interval: common.Ptr(46 * time.Second),
					Count:    common.Ptr(47),
				},
			},
		},
		HealthzEntrypoint: &kkrthttp.EntrypointConfig{
			Addr: common.Ptr("localhost:8889"),
			HTTP: &kkrthttp.ServerConfig{
				ReadTimeout:       common.Ptr(50 * time.Second),
				ReadHeaderTimeout: common.Ptr(51 * time.Second),
				WriteTimeout:      common.Ptr(52 * time.Second),
				IdleTimeout:       common.Ptr(53 * time.Second),
				MaxHeaderBytes:    common.Ptr(50000),
			},
			Net: &kkrthttp.ListenConfig{
				KeepAlive: common.Ptr(54 * time.Second),
				KeepAliveProbe: &kkrthttp.KeepAliveProbeConfig{
					Enable:   common.Ptr(true),
					Idle:     common.Ptr(55 * time.Second),
					Interval: common.Ptr(56 * time.Second),
					Count:    common.Ptr(57),
				},
			},
		},
		Log: &log.Config{
			Level: common.Ptr(log.InfoLevel),
		},
		StartTimeout: common.Ptr("10s"),
		StopTimeout:  common.Ptr("20s"),
	}
	assert.Equal(t, expectedCfg, cfg)
}

func TestEnv(t *testing.T) {
	env, err := (&Config{
		MainEntrypoint: &kkrthttp.EntrypointConfig{
			Addr: common.Ptr("localhost:8888"),
			HTTP: &kkrthttp.ServerConfig{
				ReadTimeout:       common.Ptr(40 * time.Second),
				ReadHeaderTimeout: common.Ptr(41 * time.Second),
				WriteTimeout:      common.Ptr(42 * time.Second),
				IdleTimeout:       common.Ptr(43 * time.Second),
				MaxHeaderBytes:    common.Ptr(40000),
			},
			Net: &kkrthttp.ListenConfig{
				KeepAlive: common.Ptr(44 * time.Second),
				KeepAliveProbe: &kkrthttp.KeepAliveProbeConfig{
					Enable:   common.Ptr(true),
					Idle:     common.Ptr(45 * time.Second),
					Interval: common.Ptr(46 * time.Second),
					Count:    common.Ptr(47),
				},
			},
		},
		HealthzEntrypoint: &kkrthttp.EntrypointConfig{
			Addr: common.Ptr("localhost:8889"),
			HTTP: &kkrthttp.ServerConfig{
				ReadTimeout:       common.Ptr(50 * time.Second),
				ReadHeaderTimeout: common.Ptr(51 * time.Second),
				WriteTimeout:      common.Ptr(52 * time.Second),
				IdleTimeout:       common.Ptr(53 * time.Second),
				MaxHeaderBytes:    common.Ptr(50000),
			},
			Net: &kkrthttp.ListenConfig{
				KeepAlive: common.Ptr(54 * time.Second),
				KeepAliveProbe: &kkrthttp.KeepAliveProbeConfig{
					Enable:   common.Ptr(true),
					Idle:     common.Ptr(55 * time.Second),
					Interval: common.Ptr(56 * time.Second),
					Count:    common.Ptr(57),
				},
			},
		},
		Log: &log.Config{
			Level: common.Ptr(log.InfoLevel),
		},
		StartTimeout: common.Ptr("10s"),
		StopTimeout:  common.Ptr("20s"),
	}).Env()
	require.NoError(t, err)
	assert.Equal(t, map[string]string{
		"MAIN_EP_ADDR":                             "localhost:8888",
		"MAIN_EP_HTTP_READ_TIMEOUT":                "40s",
		"MAIN_EP_HTTP_READ_HEADER_TIMEOUT":         "41s",
		"MAIN_EP_HTTP_WRITE_TIMEOUT":               "42s",
		"MAIN_EP_HTTP_IDLE_TIMEOUT":                "43s",
		"MAIN_EP_NET_KEEP_ALIVE":                   "44s",
		"MAIN_EP_NET_KEEP_ALIVE_PROBE_ENABLE":      "true",
		"MAIN_EP_NET_KEEP_ALIVE_PROBE_IDLE":        "45s",
		"MAIN_EP_NET_KEEP_ALIVE_PROBE_INTERVAL":    "46s",
		"MAIN_EP_NET_KEEP_ALIVE_PROBE_COUNT":       "47",
		"MAIN_EP_HTTP_MAX_HEADER_BYTES":            "40000",
		"HEALTHZ_EP_ADDR":                          "localhost:8889",
		"HEALTHZ_EP_HTTP_READ_TIMEOUT":             "50s",
		"HEALTHZ_EP_HTTP_READ_HEADER_TIMEOUT":      "51s",
		"HEALTHZ_EP_HTTP_WRITE_TIMEOUT":            "52s",
		"HEALTHZ_EP_HTTP_IDLE_TIMEOUT":             "53s",
		"HEALTHZ_EP_NET_KEEP_ALIVE":                "54s",
		"HEALTHZ_EP_NET_KEEP_ALIVE_PROBE_ENABLE":   "true",
		"HEALTHZ_EP_NET_KEEP_ALIVE_PROBE_IDLE":     "55s",
		"HEALTHZ_EP_NET_KEEP_ALIVE_PROBE_INTERVAL": "56s",
		"HEALTHZ_EP_NET_KEEP_ALIVE_PROBE_COUNT":    "57",
		"HEALTHZ_EP_HTTP_MAX_HEADER_BYTES":         "50000",
		"LOG_LEVEL":                                "info",
		"START_TIMEOUT":                            "10s",
		"STOP_TIMEOUT":                             "20s",
	}, env)
}

func TestAddFlagsAndLoadEnv(t *testing.T) {
	cfg := &Config{
		MainEntrypoint: &kkrthttp.EntrypointConfig{
			Addr: common.Ptr("localhost:8888"),
			HTTP: &kkrthttp.ServerConfig{
				ReadTimeout:       common.Ptr(40 * time.Second),
				ReadHeaderTimeout: common.Ptr(41 * time.Second),
				WriteTimeout:      common.Ptr(42 * time.Second),
				IdleTimeout:       common.Ptr(43 * time.Second),
				MaxHeaderBytes:    common.Ptr(40000),
			},
			Net: &kkrthttp.ListenConfig{
				KeepAlive: common.Ptr(44 * time.Second),
				KeepAliveProbe: &kkrthttp.KeepAliveProbeConfig{
					Enable:   common.Ptr(true),
					Idle:     common.Ptr(45 * time.Second),
					Interval: common.Ptr(46 * time.Second),
					Count:    common.Ptr(47),
				},
			},
		},
		HealthzEntrypoint: &kkrthttp.EntrypointConfig{
			Addr: common.Ptr("localhost:8889"),
			HTTP: &kkrthttp.ServerConfig{
				ReadTimeout:       common.Ptr(50 * time.Second),
				ReadHeaderTimeout: common.Ptr(51 * time.Second),
				WriteTimeout:      common.Ptr(52 * time.Second),
				IdleTimeout:       common.Ptr(53 * time.Second),
				MaxHeaderBytes:    common.Ptr(50000),
			},
			Net: &kkrthttp.ListenConfig{
				KeepAlive: common.Ptr(54 * time.Second),
				KeepAliveProbe: &kkrthttp.KeepAliveProbeConfig{
					Enable:   common.Ptr(true),
					Idle:     common.Ptr(55 * time.Second),
					Interval: common.Ptr(56 * time.Second),
					Count:    common.Ptr(57),
				},
			},
		},
		Log:          log.DefaultConfig(),
		StartTimeout: common.Ptr("10s"),
		StopTimeout:  common.Ptr("20s"),
	}

	v := config.NewViper()
	err := AddFlags(v, pflag.NewFlagSet("test", pflag.ContinueOnError))
	require.NoError(t, err)

	env, err := cfg.Env()
	require.NoError(t, err)
	for k, v := range env {
		t.Setenv(k, v)
	}

	loadedCfg := new(Config)
	err = loadedCfg.Unmarshal(v)
	require.NoError(t, err)
	assert.Equal(t, cfg, loadedCfg)
}

func TestUnmarshalFromDefaults(t *testing.T) {
	v := config.NewViper()
	err := AddFlags(v, pflag.NewFlagSet("test", pflag.ContinueOnError))
	require.NoError(t, err)

	cfg := new(Config)
	err = cfg.Unmarshal(v)
	require.NoError(t, err)
	assert.Equal(t, DefaultConfig(), cfg)
}
