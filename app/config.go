package app

import kkrthttp "github.com/kkrt-labs/go-utils/net/http"

type Config struct {
	MainEntrypoint    kkrthttp.EntrypointConfig `mapstructure:"main-entrypoint"`
	HealthzEntrypoint kkrthttp.EntrypointConfig `mapstructure:"healthz-entrypoint"`
	StartTimeout      string                    `mapstructure:"start-timeout"`
	StopTimeout       string                    `mapstructure:"stop-timeout"`
}
