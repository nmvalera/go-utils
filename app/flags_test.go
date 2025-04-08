package app

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddFlags(t *testing.T) {
	v := viper.New()
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	AddFlags(v, f)

	require.Len(t, v.AllKeys(), 12) // If adding a new flag, add it below
	assert.Equal(t, ":8080", v.GetString("app.main-entrypoint.addr"))
	assert.Equal(t, "0", v.GetString("app.main-entrypoint.net.keep-alive"))
	assert.Equal(t, "30s", v.GetString("app.main-entrypoint.http.read-timeout"))
	assert.Equal(t, "30s", v.GetString("app.main-entrypoint.http.read-header-timeout"))
	assert.Equal(t, "30s", v.GetString("app.main-entrypoint.http.write-timeout"))
	assert.Equal(t, "30s", v.GetString("app.main-entrypoint.http.idle-timeout"))
	assert.Equal(t, ":8081", v.GetString("app.healthz-entrypoint.addr"))
	assert.Equal(t, "0", v.GetString("app.healthz-entrypoint.net.keep-alive"))
	assert.Equal(t, "30s", v.GetString("app.healthz-entrypoint.http.read-timeout"))
	assert.Equal(t, "30s", v.GetString("app.healthz-entrypoint.http.read-header-timeout"))
	assert.Equal(t, "30s", v.GetString("app.healthz-entrypoint.http.write-timeout"))
	assert.Equal(t, "30s", v.GetString("app.healthz-entrypoint.http.idle-timeout"))

	type testConfig struct {
		App Config `mapstructure:"app"`
	}

	cfg := &testConfig{}
	err := v.Unmarshal(cfg)
	require.NoError(t, err)

	assert.Equal(t, ":8080", *cfg.App.MainEntrypoint.Addr)
	assert.Equal(t, "0", *cfg.App.MainEntrypoint.Net.KeepAlive)
	assert.Equal(t, "30s", *cfg.App.MainEntrypoint.HTTP.ReadTimeout)
	assert.Equal(t, "30s", *cfg.App.MainEntrypoint.HTTP.ReadHeaderTimeout)
	assert.Equal(t, "30s", *cfg.App.MainEntrypoint.HTTP.WriteTimeout)
	assert.Equal(t, "30s", *cfg.App.MainEntrypoint.HTTP.IdleTimeout)
	assert.Equal(t, ":8081", *cfg.App.HealthzEntrypoint.Addr)
	assert.Equal(t, "0", *cfg.App.HealthzEntrypoint.Net.KeepAlive)
	assert.Equal(t, "30s", *cfg.App.HealthzEntrypoint.HTTP.ReadTimeout)
	assert.Equal(t, "30s", *cfg.App.HealthzEntrypoint.HTTP.ReadHeaderTimeout)
	assert.Equal(t, "30s", *cfg.App.HealthzEntrypoint.HTTP.WriteTimeout)
	assert.Equal(t, "30s", *cfg.App.HealthzEntrypoint.HTTP.IdleTimeout)
}
