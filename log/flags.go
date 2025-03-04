package log

import (
	"fmt"

	"github.com/kkrt-labs/go-utils/common"
	"github.com/kkrt-labs/go-utils/spf13"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

var (
	logLevelFlag = &spf13.StringFlag{
		ViperKey:     "log.level",
		Name:         "log-level",
		Env:          "LOG_LEVEL",
		Description:  fmt.Sprintf("Log level (one of %q)", levelsStr),
		DefaultValue: common.Ptr(levelsStr[InfoLevel]),
	}
	formatFlag = &spf13.StringFlag{
		ViperKey:     "log.format",
		Name:         "log-format",
		Env:          "LOG_FORMAT",
		Description:  fmt.Sprintf("Log formatter (one of %q)", formatsStr),
		DefaultValue: common.Ptr(formatsStr[TextFormat]),
	}
	logEnableStacktraceFlag = &spf13.BoolFlag{
		ViperKey:     "log.enable-stacktrace",
		Name:         "log-enable-stacktrace",
		Env:          "LOG_ENABLE_STACKTRACE",
		Description:  "Enable extending log messages with stacktrace",
		DefaultValue: common.Ptr(false),
	}
	logEnableCallerFlag = &spf13.BoolFlag{
		ViperKey:     "log.enable-caller",
		Name:         "log-enable-caller",
		Env:          "LOG_ENABLE_CALLER",
		Description:  "Enable annotating logs with the calling function's file name and line number caller",
		DefaultValue: common.Ptr(false),
	}
	logEncoderMessageKeyFlag = &spf13.StringFlag{
		ViperKey:     "log.encoder.message-key",
		Name:         "log-encoder-message-key",
		Env:          "LOG_ENCODER_MESSAGE_KEY",
		Description:  "Log encoder message key",
		DefaultValue: common.Ptr("msg"),
	}
	logEncoderLevelKeyFlag = &spf13.StringFlag{
		ViperKey:     "log.encoder.level-key",
		Name:         "log-encoder-level-key",
		Env:          "LOG_ENCODER_LEVEL_KEY",
		Description:  "Log encoder level key",
		DefaultValue: common.Ptr("level"),
	}
	logEncoderTimeKeyFlag = &spf13.StringFlag{
		ViperKey:     "log.encoder.time-key",
		Name:         "log-encoder-time-key",
		Env:          "LOG_ENCODER_TIME_KEY",
		Description:  "Log encoder time key",
		DefaultValue: common.Ptr("ts"),
	}
	logEncoderNameKeyFlag = &spf13.StringFlag{
		ViperKey:     "log.encoder.name-key",
		Name:         "log-encoder-name-key",
		Env:          "LOG_ENCODER_NAME_KEY",
		Description:  "Log encoder name key",
		DefaultValue: common.Ptr("logger"),
	}
	logEncoderCallerKeyFlag = &spf13.StringFlag{
		ViperKey:     "log.encoder.caller-key",
		Name:         "log-encoder-caller-key",
		Env:          "LOG_ENCODER_CALLER_KEY",
		Description:  "Log encoder caller key",
		DefaultValue: common.Ptr("caller"),
	}
	logEncoderFunctionKeyFlag = &spf13.StringFlag{
		ViperKey:     "log.encoder.function-key",
		Name:         "log-encoder-function-key",
		Env:          "LOG_ENCODER_FUNCTION_KEY",
		Description:  "Log encoder function key",
		DefaultValue: common.Ptr(""),
	}
	logEncoderStacktraceKeyFlag = &spf13.StringFlag{
		ViperKey:     "log.encoder.stacktrace-key",
		Name:         "log-encoder-stacktrace-key",
		Env:          "LOG_ENCODER_STACKTRACE_KEY",
		Description:  "Log encoder stacktrace key",
		DefaultValue: common.Ptr("stacktrace"),
	}
	logEncoderSkipLineEndingFlag = &spf13.BoolFlag{
		ViperKey:     "log.encoder.skip-line-ending",
		Name:         "log-encoder-skip-line-ending",
		Env:          "LOG_ENCODER_SKIP_LINE_ENDING",
		Description:  "Log encoder skip line ending",
		DefaultValue: common.Ptr(false),
	}
	logEncoderLineEndingFlag = &spf13.StringFlag{
		ViperKey:     "log.encoder.line-ending",
		Name:         "log-encoder-line-ending",
		Env:          "LOG_ENCODER_LINE_ENDING",
		Description:  "Log encoder line ending",
		DefaultValue: common.Ptr(zapcore.DefaultLineEnding),
	}
	logEncoderLevelEncoderFlag = &spf13.StringFlag{
		ViperKey:     "log.encoder.level-encoder",
		Name:         "log-encoder-level-encoder",
		Env:          "LOG_ENCODER_LEVEL_ENCODER",
		Description:  fmt.Sprintf("Log level encoder (one of %q)", levelEncodersStr),
		DefaultValue: common.Ptr("capitalColor"),
	}
	logEncoderTimeEncoderFlag = &spf13.StringFlag{
		ViperKey:     "log.encoder.time-encoder",
		Name:         "log-encoder-time-encoder",
		Env:          "LOG_ENCODER_TIME_ENCODER",
		Description:  fmt.Sprintf("Log time encoder (one of %q)", timeEncodersStr),
		DefaultValue: common.Ptr("time"),
	}
	logEncoderDurationEncoderFlag = &spf13.StringFlag{
		ViperKey:     "log.encoder.duration-encoder",
		Name:         "log-encoder-duration-encoder",
		Env:          "LOG_ENCODER_DURATION_ENCODER",
		Description:  fmt.Sprintf("Log duration encoder (one of %q)", durationEncodersStr),
		DefaultValue: common.Ptr("s"),
	}
	logEncoderCallerEncoderFlag = &spf13.StringFlag{
		ViperKey:     "log.encoder.caller-encoder",
		Name:         "log-encoder-caller-encoder",
		Env:          "LOG_ENCODER_CALLER_ENCODER",
		Description:  fmt.Sprintf("Log caller encoder (one of %q)", callerEncodersStr),
		DefaultValue: common.Ptr("short"),
	}
	logEncoderNameEncoderFlag = &spf13.StringFlag{
		ViperKey:     "log.encoder.name-encoder",
		Name:         "log-encoder-name-encoder",
		Env:          "LOG_ENCODER_NAME_ENCODER",
		Description:  "Log name encoder",
		DefaultValue: common.Ptr("full"),
	}
	logEncoderConsoleSeparatorFlag = &spf13.StringFlag{
		ViperKey:     "log.encoder.console-separator",
		Name:         "log-encoder-console-separator",
		Env:          "LOG_ENCODER_CONSOLE_SEPARATOR",
		Description:  "Log console separator",
		DefaultValue: common.Ptr("\t"),
	}
	logSamplingInitialFlag = &spf13.IntFlag{
		ViperKey:     "log.sampling.initial",
		Name:         "log-sampling-initial",
		Env:          "LOG_SAMPLING_INITIAL",
		Description:  "Log sampling initial",
		DefaultValue: common.Ptr(100),
	}
	logSamplingThereafterFlag = &spf13.IntFlag{
		ViperKey:     "log.sampling.thereafter",
		Name:         "log-sampling-thereafter",
		Env:          "LOG_SAMPLING_THEREAFTER",
		Description:  "Log sampling thereafter",
		DefaultValue: common.Ptr(100),
	}
	logOutputPathsFlag = &spf13.StringArrayFlag{
		ViperKey:     "log.output-paths",
		Name:         "log-output-paths",
		Env:          "LOG_OUTPUT_PATHS",
		Description:  "Log output paths",
		DefaultValue: []string{"stderr"},
	}
	logErrorOutputPathsFlag = &spf13.StringArrayFlag{
		ViperKey:     "log.error-output-paths",
		Name:         "log-error-output-paths",
		Env:          "LOG_ERROR_OUTPUT_PATHS",
		Description:  "Log error output paths",
		DefaultValue: []string{"stderr"},
	}
)

func AddFlags(v *viper.Viper, f *pflag.FlagSet) {
	logLevelFlag.Add(v, f)
	formatFlag.Add(v, f)
	logEnableStacktraceFlag.Add(v, f)
	logEnableCallerFlag.Add(v, f)
	logEncoderMessageKeyFlag.Add(v, f)
	logEncoderLevelKeyFlag.Add(v, f)
	logEncoderNameKeyFlag.Add(v, f)
	logEncoderTimeKeyFlag.Add(v, f)
	logEncoderCallerKeyFlag.Add(v, f)
	logEncoderFunctionKeyFlag.Add(v, f)
	logEncoderStacktraceKeyFlag.Add(v, f)
	logEncoderSkipLineEndingFlag.Add(v, f)
	logEncoderLineEndingFlag.Add(v, f)
	logEncoderLevelEncoderFlag.Add(v, f)
	logEncoderTimeEncoderFlag.Add(v, f)
	logEncoderDurationEncoderFlag.Add(v, f)
	logEncoderCallerEncoderFlag.Add(v, f)
	logEncoderNameEncoderFlag.Add(v, f)
	logEncoderConsoleSeparatorFlag.Add(v, f)
	logSamplingInitialFlag.Add(v, f)
	logSamplingThereafterFlag.Add(v, f)
	logOutputPathsFlag.Add(v, f)
	logErrorOutputPathsFlag.Add(v, f)
}
