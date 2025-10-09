package log

import (
	"testing"

	"github.com/nmvalera/go-utils/common"
	"github.com/nmvalera/go-utils/config"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestViperConfig(t *testing.T) {
	v := config.NewViper()
	v.Set("log.level", "debug")
	v.Set("log.format", "json")
	v.Set("log.enable-stacktrace", true)
	v.Set("log.enable-caller", true)
	v.Set("log.encoding.message-key", "msg")
	v.Set("log.encoding.level-key", "level")
	v.Set("log.encoding.time-key", "time")
	v.Set("log.encoding.name-key", "logger")
	v.Set("log.encoding.caller-key", "caller")
	v.Set("log.encoding.function-key", "function")
	v.Set("log.encoding.stacktrace-key", "stacktrace")
	v.Set("log.encoding.skip-line-ending", false)
	v.Set("log.encoding.line-ending", "\n")
	v.Set("log.encoding.level-encoder", "capitalColor")
	v.Set("log.encoding.time-encoder", "rfc3339")
	v.Set("log.encoding.duration-encoder", "s")
	v.Set("log.encoding.caller-encoder", "short")
	v.Set("log.encoding.name-encoder", "full")
	v.Set("log.encoding.console-separator", "\t")
	v.Set("log.sampling.initial", 100)
	v.Set("log.sampling.thereafter", 100)
	v.Set("log.output-paths", []string{"stderr"})
	v.Set("log.error-output-paths", []string{"stderr"})

	cfg := new(Config)
	err := cfg.Unmarshal(v)
	require.NoError(t, err)

	expectedCfg := &Config{
		Format:           common.Ptr(JSONFormat),
		Level:            common.Ptr(DebugLevel),
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
			LevelEncoder:     common.Ptr(LevelEncoderCapitalColor),
			TimeEncoder:      common.Ptr(TimeEncoderRFC3339),
			DurationEncoder:  common.Ptr(DurationEncoderSeconds),
			CallerEncoder:    common.Ptr(CallerEncoderShort),
			NameEncoder:      common.Ptr(NameEncoderFull),
			ConsoleSeparator: common.Ptr("\t"),
		},
		Sampling: &SamplingConfig{
			Initial:    common.Ptr(100),
			Thereafter: common.Ptr(100),
		},
		OutputPaths:      common.PtrSlice("stderr"),
		ErrorOutputPaths: common.PtrSlice("stderr"),
	}
	assert.Equal(t, expectedCfg, cfg)
}

func TestEnv(t *testing.T) {
	env, err := (&Config{
		Level:            common.Ptr(InfoLevel),
		Format:           common.Ptr(TextFormat),
		EnableStacktrace: common.Ptr(false),
		EnableCaller:     common.Ptr(false),
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
			LevelEncoder:     common.Ptr(LevelEncoderCapitalColor),
			TimeEncoder:      common.Ptr(TimeEncoderRFC3339),
			DurationEncoder:  common.Ptr(DurationEncoderSeconds),
			CallerEncoder:    common.Ptr(CallerEncoderShort),
			NameEncoder:      common.Ptr(NameEncoderFull),
			ConsoleSeparator: common.Ptr("separator-test"),
		},
		Sampling: &SamplingConfig{
			Initial:    common.Ptr(1000),
			Thereafter: common.Ptr(1000),
		},
		OutputPaths:      common.PtrSlice("output-path-test#1,output-path-test#2"),
		ErrorOutputPaths: common.PtrSlice("error-output-path-test#1,error-output-path-test#2"),
	}).Env()
	require.NoError(t, err)
	expected := map[string]string{
		"LOG_LEVEL":                      "info",
		"LOG_FORMAT":                     "text",
		"LOG_ENABLE_STACKTRACE":          "false",
		"LOG_ENABLE_CALLER":              "false",
		"LOG_ENCODING_MESSAGE_KEY":       "msg-test",
		"LOG_ENCODING_LEVEL_KEY":         "level-test",
		"LOG_ENCODING_TIME_KEY":          "time-test",
		"LOG_ENCODING_NAME_KEY":          "logger-test",
		"LOG_ENCODING_CALLER_KEY":        "caller-test",
		"LOG_ENCODING_FUNCTION_KEY":      "function-test",
		"LOG_ENCODING_STACKTRACE_KEY":    "stacktrace-test",
		"LOG_ENCODING_SKIP_LINE_ENDING":  "true",
		"LOG_ENCODING_LINE_ENDING":       "line-ending-test",
		"LOG_ENCODING_LEVEL_ENCODER":     "capitalColor",
		"LOG_ENCODING_TIME_ENCODER":      "rfc3339",
		"LOG_ENCODING_DURATION_ENCODER":  "s",
		"LOG_ENCODING_CALLER_ENCODER":    "short",
		"LOG_ENCODING_NAME_ENCODER":      "full",
		"LOG_ENCODING_CONSOLE_SEPARATOR": "separator-test",
		"LOG_SAMPLING_INITIAL":           "1000",
		"LOG_SAMPLING_THEREAFTER":        "1000",
		"LOG_OUTPUT_PATHS":               "output-path-test#1,output-path-test#2",
		"LOG_ERROR_OUTPUT_PATHS":         "error-output-path-test#1,error-output-path-test#2",
	}
	assert.Equal(t, expected, env)
}

func TestAddFlagsAndLoadEnv(t *testing.T) {
	cfg := &Config{
		Level:            common.Ptr(InfoLevel),
		Format:           common.Ptr(TextFormat),
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
			LevelEncoder:     common.Ptr(LevelEncoderCapitalColor),
			TimeEncoder:      common.Ptr(TimeEncoderRFC3339),
			DurationEncoder:  common.Ptr(DurationEncoderSeconds),
			CallerEncoder:    common.Ptr(CallerEncoderShort),
			NameEncoder:      common.Ptr(NameEncoderFull),
			ConsoleSeparator: common.Ptr("console-separator-test"),
		},
		Sampling: &SamplingConfig{
			Initial:    common.Ptr(1000),
			Thereafter: common.Ptr(1000),
		},
		OutputPaths:      common.PtrSlice("output-path-test"),
		ErrorOutputPaths: common.PtrSlice("error-output-path-test"),
	}

	// Add the flags
	v := config.NewViper()
	set := pflag.NewFlagSet("test", pflag.ContinueOnError)
	set.SortFlags = true
	err := AddFlags(v, set)
	require.NoError(t, err)

	expectedUsage := "      --log-enable-caller                       Enable caller [env: LOG_ENABLE_CALLER]\n      --log-enable-stacktrace                   Enable automatic stacktrace capturing [env: LOG_ENABLE_STACKTRACE]\n      --log-encoding-caller-encoder string      Encoding: Primitive representation for the log caller (e.g. 'full' [env: LOG_ENCODING_CALLER_ENCODER] (default \"short\")\n      --log-encoding-caller-key string          Encoding: Key for the log caller (if empty [env: LOG_ENCODING_CALLER_KEY] (default \"caller\")\n      --log-encoding-console-separator string   Encoding: Field separator used by the console encoder [env: LOG_ENCODING_CONSOLE_SEPARATOR] (default \"\\t\")\n      --log-encoding-duration-encoder string    Encoding: Primitive representation for the log duration (e.g. 'string' [env: LOG_ENCODING_DURATION_ENCODER] (default \"s\")\n      --log-encoding-function-key string        Encoding: Key for the log function (if empty [env: LOG_ENCODING_FUNCTION_KEY]\n      --log-encoding-level-encoder string       Encoding: Primitive representation for the log level (e.g. 'capital' [env: LOG_ENCODING_LEVEL_ENCODER] (default \"capitalColor\")\n      --log-encoding-level-key string           Encoding: Key for the log level (if empty [env: LOG_ENCODING_LEVEL_KEY] (default \"level\")\n      --log-encoding-line-ending string         Encoding: Line ending [env: LOG_ENCODING_LINE_ENDING] (default \"\\n\")\n      --log-encoding-message-key string         Encoding: Key for the log message (if empty [env: LOG_ENCODING_MESSAGE_KEY] (default \"msg\")\n      --log-encoding-name-encoder string        Encoding: Primitive representation for the log logger name (e.g. 'full' [env: LOG_ENCODING_NAME_ENCODER] (default \"full\")\n      --log-encoding-name-key string            Encoding: Key for the log logger name (if empty [env: LOG_ENCODING_NAME_KEY] (default \"logger\")\n      --log-encoding-skip-line-ending           Encoding: Skip the line ending [env: LOG_ENCODING_SKIP_LINE_ENDING]\n      --log-encoding-stacktrace-key string      Encoding: Key for the log stacktrace (if empty [env: LOG_ENCODING_STACKTRACE_KEY] (default \"stacktrace\")\n      --log-encoding-time-encoder string        Encoding: Primitive representation for the log timestamp (e.g. 'rfc3339nano' [env: LOG_ENCODING_TIME_ENCODER] (default \"rfc3339\")\n      --log-encoding-time-key string            Encoding: Key for the log timestamp (if empty [env: LOG_ENCODING_TIME_KEY] (default \"ts\")\n      --log-err-output strings                  List of URLs to write internal logger errors to [env: LOG_ERROR_OUTPUT_PATHS] (default [stderr])\n      --log-format string                       Log format [env: LOG_FORMAT] (default \"text\")\n      --log-level string                        Minimum enabled logging level [env: LOG_LEVEL] (default \"info\")\n      --log-output strings                      List of URLs or file paths to write logging output to [env: LOG_OUTPUT_PATHS] (default [stderr])\n      --log-sampling-initial int                Sampling: Number of log entries with the same level and message to log before dropping entries [env: LOG_SAMPLING_INITIAL] (default 100)\n      --log-sampling-thereafter int             Sampling: After the initial number of entries [env: LOG_SAMPLING_THEREAFTER] (default 100)\n"
	assert.Equal(t, expectedUsage, set.FlagUsages())

	// Generate the environment variables
	env, err := cfg.Env()
	require.NoError(t, err)
	for k, v := range env {
		t.Setenv(k, v)
	}

	// Load the environment variables
	loadedCfg := new(Config)
	err = loadedCfg.Unmarshal(v)
	require.NoError(t, err)

	// Assert the loaded config is equal to the original config
	assert.Equal(t, *cfg, *loadedCfg)
}

func TestUnmarshalFromDefaults(t *testing.T) {
	v := config.NewViper()
	err := AddFlags(v, pflag.NewFlagSet("test", pflag.ContinueOnError))
	require.NoError(t, err)

	cfg := new(Config)
	err = cfg.Unmarshal(v)
	require.NoError(t, err)

	assert.Equal(t, DefaultConfig(), cfg)
}
