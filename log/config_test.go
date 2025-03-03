package log

import (
	"testing"

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
		Encoder: &EncoderConfig{
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
}
