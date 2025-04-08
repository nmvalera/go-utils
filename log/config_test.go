package log

import (
	"testing"

	"github.com/kkrt-labs/go-utils/common"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestParseConfig(t *testing.T) {
	cfg := &Config{
		Level:            common.Ptr("info"),
		Format:           common.Ptr("json"),
		EnableStacktrace: common.Ptr(true),
		EnableCaller:     common.Ptr(true),
		Encoder: &EncoderConfig{
			MessageKey:       common.Ptr("msg"),
			LevelKey:         common.Ptr("level"),
			TimeKey:          common.Ptr("time"),
			NameKey:          common.Ptr("logger"),
			CallerKey:        common.Ptr("caller"),
			FunctionKey:      common.Ptr("function"),
			StacktraceKey:    common.Ptr("stacktrace"),
			SkipLineEnding:   common.Ptr(false),
			LineEnding:       common.Ptr("\n"),
			LevelEncoder:     common.Ptr("capitalColor"),
			TimeEncoder:      common.Ptr("rfc3339"),
			DurationEncoder:  common.Ptr("s"),
			CallerEncoder:    common.Ptr("short"),
			NameEncoder:      common.Ptr("full"),
			ConsoleSeparator: common.Ptr("\t"),
		},
		Sampling: &SamplingConfig{
			Initial:    common.Ptr(100),
			Thereafter: common.Ptr(100),
		},
		OutputPaths:      common.PtrSlice("stderr"),
		ErrorOutputPaths: common.PtrSlice("stderr"),
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
	v := viper.New()
	v.Set("log.level", "debug")
	v.Set("log.format", "json")
	v.Set("log.enable-stacktrace", true)
	v.Set("log.enable-caller", true)
	v.Set("log.encoder.message-key", "msg")
	v.Set("log.encoder.level-key", "level")
	v.Set("log.encoder.time-key", "time")
	v.Set("log.encoder.name-key", "logger")
	v.Set("log.encoder.caller-key", "caller")
	v.Set("log.encoder.function-key", "function")
	v.Set("log.encoder.stacktrace-key", "stacktrace")
	v.Set("log.encoder.skip-line-ending", false)
	v.Set("log.encoder.line-ending", "\n")
	v.Set("log.encoder.level-encoder", "capitalColor")
	v.Set("log.encoder.time-encoder", "rfc3339")
	v.Set("log.encoder.duration-encoder", "s")
	v.Set("log.encoder.caller-encoder", "short")
	v.Set("log.encoder.name-encoder", "full")
	v.Set("log.encoder.console-separator", "\t")
	v.Set("log.sampling.initial", 100)
	v.Set("log.sampling.thereafter", 100)
	v.Set("log.output-paths", []string{"stderr"})
	v.Set("log.error-output-paths", []string{"stderr"})

	type TestConfig struct {
		Log Config `mapstructure:"log"`
	}

	var cfg TestConfig
	err := v.Unmarshal(&cfg)
	require.NoError(t, err)

	assert.Equal(t, "debug", *cfg.Log.Level)
	assert.Equal(t, "json", *cfg.Log.Format)
	assert.True(t, *cfg.Log.EnableStacktrace)
	assert.True(t, *cfg.Log.EnableCaller)

	encoderCfg := cfg.Log.Encoder
	assert.Equal(t, "msg", *encoderCfg.MessageKey)
	assert.Equal(t, "level", *encoderCfg.LevelKey)
	assert.Equal(t, "time", *encoderCfg.TimeKey)
	assert.Equal(t, "logger", *encoderCfg.NameKey)
	assert.Equal(t, "caller", *encoderCfg.CallerKey)
	assert.Equal(t, "function", *encoderCfg.FunctionKey)
	assert.Equal(t, "stacktrace", *encoderCfg.StacktraceKey)
	assert.False(t, *encoderCfg.SkipLineEnding)
	assert.Equal(t, "\n", *encoderCfg.LineEnding)
	assert.Equal(t, "\t", *encoderCfg.ConsoleSeparator)

	samplingCfg := cfg.Log.Sampling
	assert.Equal(t, 100, *samplingCfg.Initial)
	assert.Equal(t, 100, *samplingCfg.Thereafter)

	assert.Equal(t, []string{"stderr"}, common.ValSlice(*cfg.Log.OutputPaths...))
	assert.Equal(t, []string{"stderr"}, common.ValSlice(*cfg.Log.ErrorOutputPaths...))
}

func TestLoadEnv(t *testing.T) {
	cfg := &Config{
		Level:            common.Ptr("level-test"),
		Format:           common.Ptr("format-test"),
		EnableStacktrace: common.Ptr(true),
		EnableCaller:     common.Ptr(true),
		Encoder: &EncoderConfig{
			MessageKey:       common.Ptr("msg-test"),
			LevelKey:         common.Ptr("level-test"),
			TimeKey:          common.Ptr("time-test"),
			NameKey:          common.Ptr("logger-test"),
			CallerKey:        common.Ptr("caller-test"),
			FunctionKey:      common.Ptr("function-test"),
			StacktraceKey:    common.Ptr("stacktrace-test"),
			SkipLineEnding:   common.Ptr(true),
			LineEnding:       common.Ptr("line-ending-test"),
			LevelEncoder:     common.Ptr("level-encoder-test"),
			TimeEncoder:      common.Ptr("time-encoder-test"),
			DurationEncoder:  common.Ptr("duration-encoder-test"),
			CallerEncoder:    common.Ptr("caller-encoder-test"),
			NameEncoder:      common.Ptr("name-encoder-test"),
			ConsoleSeparator: common.Ptr("console-separator-test"),
		},
		Sampling: &SamplingConfig{
			Initial:    common.Ptr(1000),
			Thereafter: common.Ptr(1000),
		},
		OutputPaths:      common.PtrSlice("output-path-test"),
		ErrorOutputPaths: common.PtrSlice("error-output-path-test"),
	}

	// Generate the environment variables
	env := cfg.Env()
	for k, v := range env {
		t.Setenv(k, *v)
	}

	// Load the environment variables
	v := viper.New()
	loadedCfg := new(Config)
	err := loadedCfg.Load(v)
	require.NoError(t, err)

	// Assert the loaded config is equal to the original config
	assert.Equal(t, *cfg, *loadedCfg)
}

func TestEnv(t *testing.T) {
	cfg := &Config{
		Level: common.Ptr("level-test"),
	}

	env := cfg.Env()
	assert.Len(t, env, 1)
	assert.Equal(t, *env["LOG_LEVEL"], "level-test")
}
