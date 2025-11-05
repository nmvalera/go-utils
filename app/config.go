package app

import (
	"github.com/nmvalera/go-utils/common"
	"github.com/nmvalera/go-utils/config"
	"github.com/nmvalera/go-utils/log"
	kkrthttp "github.com/nmvalera/go-utils/net/http"
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
		HealthzServer: &HealthzServerConfig{
			LivenessPath:  common.Ptr("/live"),
			ReadinessPath: common.Ptr("/ready"),
			MetricsPath:   common.Ptr("/metrics"),
		},
		Log:          log.DefaultConfig(),
		StartTimeout: common.Ptr("10s"),
		StopTimeout:  common.Ptr("10s"),
	}
}

// Config is the configuration for the application.
type Config struct {
	MainEntrypoint    *kkrthttp.EntrypointConfig `key:"mainEp" env:"MAIN_EP" flag:"main-ep" desc:"main entrypoint: "`
	HealthzEntrypoint *kkrthttp.EntrypointConfig `key:"healthzEp" env:"HEALTHZ_EP" flag:"healthz-ep" desc:"healthz entrypoint: "`
	HealthzServer     *HealthzServerConfig       `key:"healthzApi" env:"HEALTHZ_API" flag:"healthz-api" desc:"healthz API: "`
	Log               *log.Config                `key:"log"`
	StartTimeout      *string                    `key:"startTimeout" env:"START_TIMEOUT" flag:"start-timeout" desc:"Start timeout"`
	StopTimeout       *string                    `key:"stopTimeout" env:"STOP_TIMEOUT" flag:"stop-timeout" desc:"Stop timeout"`
}

type HealthzServerConfig struct {
	LivenessPath  *string `key:"livenessPath" env:"LIVENESS_PATH" flag:"liveness-path" desc:"Path on which the liveness probe will be served"`
	ReadinessPath *string `key:"readinessPath" env:"READINESS_PATH" flag:"readiness-path" desc:"Path on which the readiness probe will be served"`
	MetricsPath   *string `key:"metricsPath" env:"METRICS_PATH" flag:"metrics-path" desc:"Path on which the metrics will be served"`
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
