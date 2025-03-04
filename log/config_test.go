package log

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestParseConfig(t *testing.T) {
	cfg := &Config{
		Level:            "info",
		Format:           "json",
		EnableStacktrace: true,
		EnableCaller:     true,
		Encoder: EncoderConfig{
			MessageKey:       "msg",
			LevelKey:         "level",
			TimeKey:          "time",
			NameKey:          "logger",
			CallerKey:        "caller",
			FunctionKey:      "function",
			StacktraceKey:    "stacktrace",
			SkipLineEnding:   false,
			LineEnding:       "\n",
			LevelEncoder:     "capitalColor",
			TimeEncoder:      "rfc3339",
			DurationEncoder:  "s",
			CallerEncoder:    "short",
			NameEncoder:      "full",
			ConsoleSeparator: "\t",
		},
		Sampling: SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	zapCfg, err := ParseConfig(cfg)
	require.NoError(t, err)

	assert.Equal(t, zap.InfoLevel, zapCfg.Level.Level())
	assert.Equal(t, "json", zapCfg.Encoding)
	assert.False(t, zapCfg.DisableStacktrace)
	assert.False(t, zapCfg.DisableCaller)

	encoderCfg := zapCfg.EncoderConfig
	assert.Equal(t, "msg", encoderCfg.MessageKey)
	assert.Equal(t, "level", encoderCfg.LevelKey)
	assert.Equal(t, "time", encoderCfg.TimeKey)
	assert.Equal(t, "logger", encoderCfg.NameKey)
	assert.Equal(t, "caller", encoderCfg.CallerKey)
	assert.Equal(t, "function", encoderCfg.FunctionKey)
	assert.Equal(t, "stacktrace", encoderCfg.StacktraceKey)
	assert.False(t, encoderCfg.SkipLineEnding)
	assert.Equal(t, "\n", encoderCfg.LineEnding)
	assert.Equal(t, "\t", encoderCfg.ConsoleSeparator)

	samplingCfg := zapCfg.Sampling
	assert.Equal(t, 100, samplingCfg.Initial)
	assert.Equal(t, 100, samplingCfg.Thereafter)

	outputPaths := zapCfg.OutputPaths
	assert.Equal(t, []string{"stderr"}, outputPaths)

	errorOutputPaths := zapCfg.ErrorOutputPaths
	assert.Equal(t, []string{"stderr"}, errorOutputPaths)
}

func TestViperConfig(t *testing.T) {
	type TestConfig struct {
		Log Config `mapstructure:"log"`
	}

	viper.Set("log.level", "debug")
	viper.Set("log.format", "json")
	viper.Set("log.enable-stacktrace", true)
	viper.Set("log.enable-caller", true)
	viper.Set("log.encoder.message-key", "msg")
	viper.Set("log.encoder.level-key", "level")
	viper.Set("log.encoder.time-key", "time")
	viper.Set("log.encoder.name-key", "logger")
	viper.Set("log.encoder.caller-key", "caller")
	viper.Set("log.encoder.function-key", "function")
	viper.Set("log.encoder.stacktrace-key", "stacktrace")
	viper.Set("log.encoder.skip-line-ending", false)
	viper.Set("log.encoder.line-ending", "\n")
	viper.Set("log.encoder.level-encoder", "capitalColor")
	viper.Set("log.encoder.time-encoder", "rfc3339")
	viper.Set("log.encoder.duration-encoder", "s")
	viper.Set("log.encoder.caller-encoder", "short")
	viper.Set("log.encoder.name-encoder", "full")
	viper.Set("log.encoder.console-separator", "\t")
	viper.Set("log.sampling.initial", 100)
	viper.Set("log.sampling.thereafter", 100)
	viper.Set("log.output-paths", []string{"stderr"})
	viper.Set("log.error-output-paths", []string{"stderr"})

	var cfg TestConfig
	err := viper.Unmarshal(&cfg)
	require.NoError(t, err)

	assert.Equal(t, "debug", cfg.Log.Level)
	assert.Equal(t, "json", cfg.Log.Format)
	assert.True(t, cfg.Log.EnableStacktrace)
	assert.True(t, cfg.Log.EnableCaller)

	encoderCfg := cfg.Log.Encoder
	assert.Equal(t, "msg", encoderCfg.MessageKey)
	assert.Equal(t, "level", encoderCfg.LevelKey)
	assert.Equal(t, "time", encoderCfg.TimeKey)
	assert.Equal(t, "logger", encoderCfg.NameKey)
	assert.Equal(t, "caller", encoderCfg.CallerKey)
	assert.Equal(t, "function", encoderCfg.FunctionKey)
	assert.Equal(t, "stacktrace", encoderCfg.StacktraceKey)
	assert.False(t, encoderCfg.SkipLineEnding)
	assert.Equal(t, "\n", encoderCfg.LineEnding)
	assert.Equal(t, "\t", encoderCfg.ConsoleSeparator)

	samplingCfg := cfg.Log.Sampling
	assert.Equal(t, 100, samplingCfg.Initial)
	assert.Equal(t, 100, samplingCfg.Thereafter)

	assert.Equal(t, []string{"stderr"}, cfg.Log.OutputPaths)
	assert.Equal(t, []string{"stderr"}, cfg.Log.ErrorOutputPaths)
}
