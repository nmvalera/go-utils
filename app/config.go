package app

import (
	"github.com/kkrt-labs/go-utils/common"
	"github.com/kkrt-labs/go-utils/config"
	"github.com/kkrt-labs/go-utils/log"
	kkrthttp "github.com/kkrt-labs/go-utils/net/http"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func DefaultConfig() *Config {
	mainEp := kkrthttp.DefaultEntrypointConfig()
	mainEp.Addr = common.Ptr(":8080")
	healthzEp := kkrthttp.DefaultEntrypointConfig()
	healthzEp.Addr = common.Ptr(":8081")
	return &Config{
		MainEntrypoint:    mainEp,
		HealthzEntrypoint: healthzEp,
		Log:               log.DefaultConfig(),
		StartTimeout:      common.Ptr("10s"),
		StopTimeout:       common.Ptr("10s"),
	}
}

// Config is the configuration for the application.
type Config struct {
	MainEntrypoint    *kkrthttp.EntrypointConfig `key:"main-ep" env:"MAIN_EP" flag:"main-ep" desc:"Main entrypoint"`
	HealthzEntrypoint *kkrthttp.EntrypointConfig `key:"healthz-ep" env:"HEALTHZ_EP" flag:"healthz-ep" desc:"Healthz entrypoint"`
	Log               *log.Config                `key:"log"`
	StartTimeout      *string                    `key:"start-timeout" env:"START_TIMEOUT" flag:"start-timeout" desc:"Start timeout"`
	StopTimeout       *string                    `key:"stop-timeout" env:"STOP_TIMEOUT" flag:"stop-timeout" desc:"Stop timeout"`
}

// Env returns the environment variables for the given Config.
func (cfg *Config) Env() (map[string]string, error) {
	return config.Env(cfg, nil)
}

// Unmarshal unmarshals the given viper into the Config.
func (cfg *Config) Unmarshal(v *viper.Viper) error {
	return config.Unmarshal(cfg, v)
}

// AddFlags adds flags to the given viper and pflag.FlagSet.
func AddFlags(v *viper.Viper, f *pflag.FlagSet) error {
	return config.AddFlags(DefaultConfig(), v, f, nil)
}
