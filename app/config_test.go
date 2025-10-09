package app

import (
	"testing"
	"time"

	"github.com/nmvalera/go-utils/common"
	"github.com/nmvalera/go-utils/config"
	"github.com/nmvalera/go-utils/log"
	kkrthttp "github.com/nmvalera/go-utils/net/http"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestViperConfig(t *testing.T) {
	v := config.NewViper()
	v.Set("main-ep.addr", "localhost:8881")
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
	v.Set("healthz-ep.addr", "localhost:8882")
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
			Addr: common.Ptr("localhost:8881"),
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
			Addr: common.Ptr("localhost:8882"),
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
			Addr: common.Ptr("localhost:8883"),
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
			Addr: common.Ptr("localhost:8884"),
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
		"MAIN_EP_ADDR":                             "localhost:8883",
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
		"HEALTHZ_EP_ADDR":                          "localhost:8884",
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
			Addr: common.Ptr("localhost:8885"),
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
			TLS: &kkrthttp.TLSCertConfig{},
		},
		HealthzEntrypoint: &kkrthttp.EntrypointConfig{
			Addr: common.Ptr("localhost:8886"),
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
			TLS: &kkrthttp.TLSCertConfig{},
		},
		HealthzServer: &HealthzServerConfig{
			LivenessPath:  common.Ptr("/live"),
			ReadinessPath: common.Ptr("/ready"),
			MetricsPath:   common.Ptr("/metrics"),
		},
		Log:          log.DefaultConfig(),
		StartTimeout: common.Ptr("10s"),
		StopTimeout:  common.Ptr("20s"),
	}

	v := config.NewViper()
	set := pflag.NewFlagSet("test", pflag.ContinueOnError)
	set.SortFlags = true
	err := AddFlags(v, set)
	require.NoError(t, err)

	expectedUsage := "      --healthz-api-liveness-path string                  healthz API configurationPath on which the liveness probe will be served [env: HEALTHZ_API_LIVENESS_PATH] (default \"/live\")\n      --healthz-api-metrics-path string                   healthz API configurationPath on which the metrics will be served [env: HEALTHZ_API_METRICS_PATH] (default \"/metrics\")\n      --healthz-api-readiness-path string                 healthz API configurationPath on which the readiness probe will be served [env: HEALTHZ_API_READINESS_PATH] (default \"/ready\")\n      --healthz-ep-addr string                            healthz entrypoint: TCP Address to listen on [env: HEALTHZ_EP_ADDR] (default \":8081\")\n      --healthz-ep-http-idle-timeout string               healthz entrypoint: Maximum duration to wait for the next request when keep-alives are enabled (zero uses the value of read timeout) [env: HEALTHZ_EP_HTTP_IDLE_TIMEOUT] (default \"30s\")\n      --healthz-ep-http-max-header-bytes int              healthz entrypoint: Maximum number of bytes the server will read parsing the request header's keys and values [env: HEALTHZ_EP_HTTP_MAX_HEADER_BYTES] (default 1048576)\n      --healthz-ep-http-read-header-timeout string        healthz entrypoint: Maximum duration for reading request headers (zero uses the value of read timeout) [env: HEALTHZ_EP_HTTP_READ_HEADER_TIMEOUT] (default \"30s\")\n      --healthz-ep-http-read-timeout string               healthz entrypoint: Maximum duration for reading the entire request including the body (zero means no timeout) [env: HEALTHZ_EP_HTTP_READ_TIMEOUT] (default \"30s\")\n      --healthz-ep-http-write-timeout string              healthz entrypoint: Maximum duration before timing out writes of the response (zero means no timeout) [env: HEALTHZ_EP_HTTP_WRITE_TIMEOUT] (default \"30s\")\n      --healthz-ep-net-keep-alive string                  healthz entrypoint: Keep alive period for network connections accepted by this entrypoint [env: HEALTHZ_EP_NET_KEEP_ALIVE] (default \"-1s\")\n      --healthz-ep-net-keep-alive-probe-count int         healthz entrypoint: Maximum number of keep-alive probes that can go unanswered before dropping a connection [env: HEALTHZ_EP_NET_KEEP_ALIVE_PROBE_COUNT] (default 9)\n      --healthz-ep-net-keep-alive-probe-enable            healthz entrypoint: Enable keep alive probes [env: HEALTHZ_EP_NET_KEEP_ALIVE_PROBE_ENABLE]\n      --healthz-ep-net-keep-alive-probe-idle string       healthz entrypoint: Time that the connection must be idle before the first keep-alive probe is sent [env: HEALTHZ_EP_NET_KEEP_ALIVE_PROBE_IDLE] (default \"15s\")\n      --healthz-ep-net-keep-alive-probe-interval string   healthz entrypoint: Time between keep-alive probes [env: HEALTHZ_EP_NET_KEEP_ALIVE_PROBE_INTERVAL] (default \"15s\")\n      --healthz-ep-tls-certfile string                    healthz entrypoint: Path to the certificate file [env: HEALTHZ_EP_TLS_CERT_FILE]\n      --healthz-ep-tls-keyfile string                     healthz entrypoint: Path to the key file [env: HEALTHZ_EP_TLS_KEY_FILE]\n      --log-enable-caller                                 Enable caller [env: LOG_ENABLE_CALLER]\n      --log-enable-stacktrace                             Enable automatic stacktrace capturing [env: LOG_ENABLE_STACKTRACE]\n      --log-encoding-caller-encoder string                Encoding: Primitive representation for the log caller (e.g. 'full' [env: LOG_ENCODING_CALLER_ENCODER] (default \"short\")\n      --log-encoding-caller-key string                    Encoding: Key for the log caller (if empty [env: LOG_ENCODING_CALLER_KEY] (default \"caller\")\n      --log-encoding-console-separator string             Encoding: Field separator used by the console encoder [env: LOG_ENCODING_CONSOLE_SEPARATOR] (default \"\\t\")\n      --log-encoding-duration-encoder string              Encoding: Primitive representation for the log duration (e.g. 'string' [env: LOG_ENCODING_DURATION_ENCODER] (default \"s\")\n      --log-encoding-function-key string                  Encoding: Key for the log function (if empty [env: LOG_ENCODING_FUNCTION_KEY]\n      --log-encoding-level-encoder string                 Encoding: Primitive representation for the log level (e.g. 'capital' [env: LOG_ENCODING_LEVEL_ENCODER] (default \"capitalColor\")\n      --log-encoding-level-key string                     Encoding: Key for the log level (if empty [env: LOG_ENCODING_LEVEL_KEY] (default \"level\")\n      --log-encoding-line-ending string                   Encoding: Line ending [env: LOG_ENCODING_LINE_ENDING] (default \"\\n\")\n      --log-encoding-message-key string                   Encoding: Key for the log message (if empty [env: LOG_ENCODING_MESSAGE_KEY] (default \"msg\")\n      --log-encoding-name-encoder string                  Encoding: Primitive representation for the log logger name (e.g. 'full' [env: LOG_ENCODING_NAME_ENCODER] (default \"full\")\n      --log-encoding-name-key string                      Encoding: Key for the log logger name (if empty [env: LOG_ENCODING_NAME_KEY] (default \"logger\")\n      --log-encoding-skip-line-ending                     Encoding: Skip the line ending [env: LOG_ENCODING_SKIP_LINE_ENDING]\n      --log-encoding-stacktrace-key string                Encoding: Key for the log stacktrace (if empty [env: LOG_ENCODING_STACKTRACE_KEY] (default \"stacktrace\")\n      --log-encoding-time-encoder string                  Encoding: Primitive representation for the log timestamp (e.g. 'rfc3339nano' [env: LOG_ENCODING_TIME_ENCODER] (default \"rfc3339\")\n      --log-encoding-time-key string                      Encoding: Key for the log timestamp (if empty [env: LOG_ENCODING_TIME_KEY] (default \"ts\")\n      --log-err-output strings                            List of URLs to write internal logger errors to [env: LOG_ERROR_OUTPUT_PATHS] (default [stderr])\n      --log-format string                                 Log format [env: LOG_FORMAT] (default \"text\")\n      --log-level string                                  Minimum enabled logging level [env: LOG_LEVEL] (default \"info\")\n      --log-output strings                                List of URLs or file paths to write logging output to [env: LOG_OUTPUT_PATHS] (default [stderr])\n      --log-sampling-initial int                          Sampling: Number of log entries with the same level and message to log before dropping entries [env: LOG_SAMPLING_INITIAL] (default 100)\n      --log-sampling-thereafter int                       Sampling: After the initial number of entries [env: LOG_SAMPLING_THEREAFTER] (default 100)\n      --main-ep-addr string                               main entrypoint: TCP Address to listen on [env: MAIN_EP_ADDR] (default \":8080\")\n      --main-ep-http-idle-timeout string                  main entrypoint: Maximum duration to wait for the next request when keep-alives are enabled (zero uses the value of read timeout) [env: MAIN_EP_HTTP_IDLE_TIMEOUT] (default \"30s\")\n      --main-ep-http-max-header-bytes int                 main entrypoint: Maximum number of bytes the server will read parsing the request header's keys and values [env: MAIN_EP_HTTP_MAX_HEADER_BYTES] (default 1048576)\n      --main-ep-http-read-header-timeout string           main entrypoint: Maximum duration for reading request headers (zero uses the value of read timeout) [env: MAIN_EP_HTTP_READ_HEADER_TIMEOUT] (default \"30s\")\n      --main-ep-http-read-timeout string                  main entrypoint: Maximum duration for reading the entire request including the body (zero means no timeout) [env: MAIN_EP_HTTP_READ_TIMEOUT] (default \"30s\")\n      --main-ep-http-write-timeout string                 main entrypoint: Maximum duration before timing out writes of the response (zero means no timeout) [env: MAIN_EP_HTTP_WRITE_TIMEOUT] (default \"30s\")\n      --main-ep-net-keep-alive string                     main entrypoint: Keep alive period for network connections accepted by this entrypoint [env: MAIN_EP_NET_KEEP_ALIVE] (default \"-1s\")\n      --main-ep-net-keep-alive-probe-count int            main entrypoint: Maximum number of keep-alive probes that can go unanswered before dropping a connection [env: MAIN_EP_NET_KEEP_ALIVE_PROBE_COUNT] (default 9)\n      --main-ep-net-keep-alive-probe-enable               main entrypoint: Enable keep alive probes [env: MAIN_EP_NET_KEEP_ALIVE_PROBE_ENABLE]\n      --main-ep-net-keep-alive-probe-idle string          main entrypoint: Time that the connection must be idle before the first keep-alive probe is sent [env: MAIN_EP_NET_KEEP_ALIVE_PROBE_IDLE] (default \"15s\")\n      --main-ep-net-keep-alive-probe-interval string      main entrypoint: Time between keep-alive probes [env: MAIN_EP_NET_KEEP_ALIVE_PROBE_INTERVAL] (default \"15s\")\n      --main-ep-tls-certfile string                       main entrypoint: Path to the certificate file [env: MAIN_EP_TLS_CERT_FILE]\n      --main-ep-tls-keyfile string                        main entrypoint: Path to the key file [env: MAIN_EP_TLS_KEY_FILE]\n      --start-timeout string                              Start timeout [env: START_TIMEOUT] (default \"10s\")\n      --stop-timeout string                               Stop timeout [env: STOP_TIMEOUT] (default \"10s\")\n"
	assert.Equal(t, expectedUsage, set.FlagUsages())

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
