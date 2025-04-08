package app

import (
	"net/http"

	"github.com/kkrt-labs/go-utils/common"
	kkrthttp "github.com/kkrt-labs/go-utils/net/http"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	MainEntrypoint    *kkrthttp.EntrypointConfig `mapstructure:"main-entrypoint"`
	HealthzEntrypoint *kkrthttp.EntrypointConfig `mapstructure:"healthz-entrypoint"`
	StartTimeout      *string                    `mapstructure:"start-timeout"`
	StopTimeout       *string                    `mapstructure:"stop-timeout"`
}

func (cfg *Config) Load(v *viper.Viper) error {
	type embedConfig struct {
		App *Config `mapstructure:"app"`
	}
	AddFlags(v, new(pflag.FlagSet))
	return v.Unmarshal(&embedConfig{cfg})
}

func (cfg *Config) Env() map[string]*string {
	m := make(map[string]*string)
	if cfg.MainEntrypoint != nil {
		if cfg.MainEntrypoint.Addr != nil {
			m[mainEntrypointFlag.Env] = cfg.MainEntrypoint.Addr
		}
		if cfg.MainEntrypoint.HTTP != nil {
			if cfg.MainEntrypoint.HTTP.ReadTimeout != nil {
				m[mainReadTimeoutFlag.Env] = cfg.MainEntrypoint.HTTP.ReadTimeout
			}
			if cfg.MainEntrypoint.HTTP.ReadHeaderTimeout != nil {
				m[mainReadHeaderTimeoutFlag.Env] = cfg.MainEntrypoint.HTTP.ReadHeaderTimeout
			}
			if cfg.MainEntrypoint.HTTP.WriteTimeout != nil {
				m[mainWriteTimeoutFlag.Env] = cfg.MainEntrypoint.HTTP.WriteTimeout
			}
			if cfg.MainEntrypoint.HTTP.IdleTimeout != nil {
				m[mainIdleTimeoutFlag.Env] = cfg.MainEntrypoint.HTTP.IdleTimeout
			}
		}
		if cfg.MainEntrypoint.Net != nil {
			if cfg.MainEntrypoint.Net.KeepAlive != nil {
				m[mainKeepAliveFlag.Env] = cfg.MainEntrypoint.Net.KeepAlive
			}
		}
	}

	if cfg.HealthzEntrypoint != nil {
		if cfg.HealthzEntrypoint.Addr != nil {
			m[healthzEntrypointFlag.Env] = cfg.HealthzEntrypoint.Addr
		}
		if cfg.HealthzEntrypoint.HTTP != nil {
			if cfg.HealthzEntrypoint.HTTP.ReadTimeout != nil {
				m[healthzReadTimeoutFlag.Env] = cfg.HealthzEntrypoint.HTTP.ReadTimeout
			}
			if cfg.HealthzEntrypoint.HTTP.ReadHeaderTimeout != nil {
				m[healthzReadHeaderTimeoutFlag.Env] = cfg.HealthzEntrypoint.HTTP.ReadHeaderTimeout
			}
			if cfg.HealthzEntrypoint.HTTP.WriteTimeout != nil {
				m[healthzWriteTimeoutFlag.Env] = cfg.HealthzEntrypoint.HTTP.WriteTimeout
			}
			if cfg.HealthzEntrypoint.HTTP.IdleTimeout != nil {
				m[healthzIdleTimeoutFlag.Env] = cfg.HealthzEntrypoint.HTTP.IdleTimeout
			}
		}
		if cfg.HealthzEntrypoint.Net != nil {
			if cfg.HealthzEntrypoint.Net.KeepAlive != nil {
				m[healthzKeepAliveFlag.Env] = cfg.HealthzEntrypoint.Net.KeepAlive
			}
		}
	}

	return m
}

var defaultConfig = &Config{
	MainEntrypoint: &kkrthttp.EntrypointConfig{
		Addr: common.Ptr(":8080"),
		HTTP: &kkrthttp.ServerConfig{
			ReadTimeout:       common.Ptr("30s"),
			ReadHeaderTimeout: common.Ptr("30s"),
			WriteTimeout:      common.Ptr("30s"),
			IdleTimeout:       common.Ptr("30s"),
			MaxHeaderBytes:    common.Ptr(http.DefaultMaxHeaderBytes),
		},
		Net: &kkrthttp.ListenConfig{
			KeepAlive: common.Ptr("0"),
			KeepAliveProbe: &kkrthttp.KeepAliveProbeConfig{
				Enable:   common.Ptr(false),
				Idle:     common.Ptr("15s"),
				Interval: common.Ptr("15s"),
				Count:    common.Ptr(9),
			},
		},
	},
	HealthzEntrypoint: &kkrthttp.EntrypointConfig{
		Addr: common.Ptr(":8081"),
		HTTP: &kkrthttp.ServerConfig{
			ReadTimeout:       common.Ptr("30s"),
			ReadHeaderTimeout: common.Ptr("30s"),
			WriteTimeout:      common.Ptr("30s"),
			IdleTimeout:       common.Ptr("30s"),
			MaxHeaderBytes:    common.Ptr(http.DefaultMaxHeaderBytes),
		},
		Net: &kkrthttp.ListenConfig{
			KeepAlive: common.Ptr("0"),
			KeepAliveProbe: &kkrthttp.KeepAliveProbeConfig{
				Enable:   common.Ptr(false),
				Idle:     common.Ptr("15s"),
				Interval: common.Ptr("15s"),
				Count:    common.Ptr(9),
			},
		},
	},
	StartTimeout: common.Ptr("10s"),
	StopTimeout:  common.Ptr("10s"),
}
