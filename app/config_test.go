package app

import (
	"testing"

	"github.com/kkrt-labs/go-utils/common"
	kkrthttp "github.com/kkrt-labs/go-utils/net/http"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestViperConfig(t *testing.T) {
	v := viper.New()
	v.Set("app.main-entrypoint.addr", "main-ep-addr-test")
	v.Set("app.main-entrypoint.http.read-timeout", "main-ep-http-read-timeout-test")
	v.Set("app.main-entrypoint.http.write-timeout", "main-ep-http-write-timeout-test")
	v.Set("app.main-entrypoint.http.idle-timeout", "main-ep-http-idle-timeout-test")
	v.Set("app.main-entrypoint.net.keep-alive", "main-ep-net-keep-alive-test")
	v.Set("app.main-entrypoint.net.keep-alive-probe.enable", "true")
	v.Set("app.healthz-entrypoint.addr", "healthz-ep-addr-test")
	v.Set("app.healthz-entrypoint.http.read-timeout", "healthz-ep-http-read-timeout-test")
	v.Set("app.healthz-entrypoint.http.write-timeout", "healthz-ep-http-write-timeout-test")
	v.Set("app.healthz-entrypoint.http.idle-timeout", "healthz-ep-http-idle-timeout-test")
	v.Set("app.healthz-entrypoint.net.keep-alive", "healthz-ep-net-keep-alive-test")
	v.Set("app.healthz-entrypoint.net.keep-alive-probe.enable", "true")

	type TestConfig struct {
		App Config `mapstructure:"app"`
	}

	cfg := &TestConfig{}
	err := v.Unmarshal(cfg)
	require.NoError(t, err)

	assert.Equal(t, "main-ep-addr-test", *cfg.App.MainEntrypoint.Addr)
	assert.Equal(t, "main-ep-http-read-timeout-test", *cfg.App.MainEntrypoint.HTTP.ReadTimeout)
	assert.Equal(t, "main-ep-http-write-timeout-test", *cfg.App.MainEntrypoint.HTTP.WriteTimeout)
	assert.Equal(t, "main-ep-http-idle-timeout-test", *cfg.App.MainEntrypoint.HTTP.IdleTimeout)
	assert.Equal(t, "main-ep-net-keep-alive-test", *cfg.App.MainEntrypoint.Net.KeepAlive)
	assert.True(t, *cfg.App.MainEntrypoint.Net.KeepAliveProbe.Enable)
	assert.Equal(t, "healthz-ep-addr-test", *cfg.App.HealthzEntrypoint.Addr)
	assert.Equal(t, "healthz-ep-http-read-timeout-test", *cfg.App.HealthzEntrypoint.HTTP.ReadTimeout)
	assert.Equal(t, "healthz-ep-http-write-timeout-test", *cfg.App.HealthzEntrypoint.HTTP.WriteTimeout)
	assert.Equal(t, "healthz-ep-http-idle-timeout-test", *cfg.App.HealthzEntrypoint.HTTP.IdleTimeout)
	assert.Equal(t, "healthz-ep-net-keep-alive-test", *cfg.App.HealthzEntrypoint.Net.KeepAlive)
	assert.True(t, *cfg.App.HealthzEntrypoint.Net.KeepAliveProbe.Enable)
}

func TestLoadEnv(t *testing.T) {
	cfg := &Config{
		MainEntrypoint: &kkrthttp.EntrypointConfig{
			Addr: common.Ptr("main-ep-addr-test"),
			HTTP: &kkrthttp.ServerConfig{
				ReadTimeout:       common.Ptr("main-ep-http-read-timeout-test"),
				ReadHeaderTimeout: common.Ptr("main-ep-http-read-header-timeout-test"),
				WriteTimeout:      common.Ptr("main-ep-http-write-timeout-test"),
				IdleTimeout:       common.Ptr("main-ep-http-idle-timeout-test"),
			},
			Net: &kkrthttp.ListenConfig{
				KeepAlive: common.Ptr("main-ep-net-keep-alive-test"),
			},
		},
		HealthzEntrypoint: &kkrthttp.EntrypointConfig{
			Addr: common.Ptr("healthz-ep-addr-test"),
			HTTP: &kkrthttp.ServerConfig{
				ReadTimeout:       common.Ptr("healthz-ep-http-read-timeout-test"),
				ReadHeaderTimeout: common.Ptr("healthz-ep-http-read-header-timeout-test"),
				WriteTimeout:      common.Ptr("healthz-ep-http-write-timeout-test"),
				IdleTimeout:       common.Ptr("healthz-ep-http-idle-timeout-test"),
			},
			Net: &kkrthttp.ListenConfig{
				KeepAlive: common.Ptr("healthz-ep-net-keep-alive-test"),
			},
		},
	}

	// Generate the environment variables
	env := cfg.Env()
	for k, v := range env {
		t.Setenv(k, *v)
	}

	// Load configuration from generated environment variables
	v := viper.New()
	loadedCfg := new(Config)
	err := loadedCfg.Load(v)
	require.NoError(t, err)

	// Assert the loaded config is equal to the original config
	assert.Equal(t, cfg, loadedCfg)
}
