package log

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestAddFlags(t *testing.T) {
	v := viper.New()
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	AddFlags(v, f)

	assert.Equal(t, "info", v.GetString("log.level"))
	assert.Equal(t, "text", v.GetString("log.format"))
	assert.Equal(t, false, v.GetBool("log.enable-stacktrace"))
	assert.Equal(t, false, v.GetBool("log.enable-caller"))
	assert.Equal(t, "msg", v.GetString("log.encoder.message-key"))
	assert.Equal(t, "level", v.GetString("log.encoder.level-key"))
	assert.Equal(t, "ts", v.GetString("log.encoder.time-key"))
	assert.Equal(t, "logger", v.GetString("log.encoder.name-key"))
	assert.Equal(t, "caller", v.GetString("log.encoder.caller-key"))
	assert.Equal(t, "stacktrace", v.GetString("log.encoder.stacktrace-key"))
	assert.Equal(t, false, v.GetBool("log.encoder.skip-line-ending"))
	assert.Equal(t, "capitalColor", v.GetString("log.encoder.level-encoder"))
	assert.Equal(t, "time", v.GetString("log.encoder.time-encoder"))
	assert.Equal(t, "s", v.GetString("log.encoder.duration-encoder"))
	assert.Equal(t, "short", v.GetString("log.encoder.caller-encoder"))
	assert.Equal(t, "full", v.GetString("log.encoder.name-encoder"))
	assert.Equal(t, "\t", v.GetString("log.encoder.console-separator"))
}
