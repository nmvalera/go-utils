package log

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/kkrt-labs/go-utils/common"
	"github.com/kkrt-labs/go-utils/config"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	// Register decode hooks for the Config types declared in this file.
	// We do not need to register encoders because default config.StringerEncodeHook works for all types in this file.
	config.RegisterGlobalDecodeHooks(
		func(f reflect.Type, t reflect.Type, data any) (any, error) {
			// Parse Format
			if f.Kind() != reflect.String {
				return data, nil
			}

			if t == reflect.TypeOf(Format(0)) {
				return ParseFormat(data.(string))
			}

			// Parse Level
			if t == reflect.TypeOf(Level(0)) {
				return ParseLevel(data.(string))
			}

			// Parse LevelEncoder
			if t == reflect.TypeOf(LevelEncoder(0)) {
				return ParseLevelEncoder(data.(string))
			}

			// Parse TimeEncoder
			if t == reflect.TypeOf(TimeEncoder(0)) {
				return ParseTimeEncoder(data.(string))
			}

			// Parse DurationEncoder
			if t == reflect.TypeOf(DurationEncoder(0)) {
				return ParseDurationEncoder(data.(string))
			}

			// Parse CallerEncoder
			if t == reflect.TypeOf(CallerEncoder(0)) {
				return ParseCallerEncoder(data.(string))
			}

			// Parse NameEncoder
			if t == reflect.TypeOf(NameEncoder(0)) {
				return ParseNameEncoder(data.(string))
			}

			return data, nil
		},
	)
}

// DefaultConfig returns a default Config.
func DefaultConfig() *Config {
	return &Config{
		Level:            common.Ptr(InfoLevel),
		Format:           common.Ptr(TextFormat),
		EnableStacktrace: common.Ptr(false),
		EnableCaller:     common.Ptr(false),
		Encoder: &EncoderConfig{
			MessageKey:       common.Ptr("msg"),
			LevelKey:         common.Ptr("level"),
			TimeKey:          common.Ptr("ts"),
			NameKey:          common.Ptr("logger"),
			CallerKey:        common.Ptr("caller"),
			FunctionKey:      common.Ptr(""),
			StacktraceKey:    common.Ptr("stacktrace"),
			SkipLineEnding:   common.Ptr(false),
			LineEnding:       common.Ptr(zapcore.DefaultLineEnding),
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

}

// Config is the configuration to create a zap.Config.
type Config struct {
	Format           *Format         `key:"format,omitempty" desc:"Log format"`
	Level            *Level          `key:"level,omitempty" desc:"Minimum enabled logging level"`
	EnableStacktrace *bool           `key:"enable-stacktrace,omitempty" env:"ENABLE_STACKTRACE" flag:"enable-stacktrace" desc:"Enable automatic stacktrace capturing"`
	EnableCaller     *bool           `key:"enable-caller,omitempty" env:"ENABLE_CALLER" flag:"enable-caller" desc:"Enable caller"`
	Encoder          *EncoderConfig  `key:"encoding,omitempty" flag:"encoding" env:"ENCODING" desc:"Encoding: "`
	Sampling         *SamplingConfig `key:"sampling,omitempty" desc:"Sampling: "`
	OutputPaths      *[]*string      `key:"output-paths,omitempty" env:"OUTPUT_PATHS" flag:"output" desc:"List of URLs or file paths to write logging output to"`
	ErrorOutputPaths *[]*string      `key:"error-output-paths,omitempty" env:"ERROR_OUTPUT_PATHS" flag:"err-output" desc:"List of URLs to write internal logger errors to"`
}

func (cfg *Config) ZapConfig() *zap.Config {
	return &zap.Config{
		DisableStacktrace: !common.Val(cfg.EnableStacktrace),
		DisableCaller:     !common.Val(cfg.EnableCaller),
		Level:             zap.NewAtomicLevelAt(zapLevels[*cfg.Level]),
		Sampling: &zap.SamplingConfig{
			Initial:    common.Val(cfg.Sampling.Initial),
			Thereafter: common.Val(cfg.Sampling.Thereafter),
		},
		Encoding:         zapFormats[*cfg.Format],
		EncoderConfig:    *cfg.Encoder.EncoderConfig(),
		OutputPaths:      common.ValSlice(common.Val(cfg.OutputPaths)...),
		ErrorOutputPaths: common.ValSlice(common.Val(cfg.ErrorOutputPaths)...),
	}
}

type embedConfig struct {
	Log *Config `key:"log"`
}

// Env returns the environment variables for the given Config.
// All environment variables are prefixed with "LOG_".
func (cfg *Config) Env() (map[string]string, error) {
	return config.Env(&embedConfig{cfg}, nil)
}

// Unmarshal unmarshals the given viper into the Config.
// Assumes
// - all viper keys are prefixed with "log."
// - all environment variables are prefixed with "LOG_".
func (cfg *Config) Unmarshal(v *viper.Viper) error {
	return config.Unmarshal(&embedConfig{cfg}, v)
}

// AddFlags adds flags to the given viper and pflag.FlagSet.
// Sets
// - all viper keys with "log." prefix
// - all environment variables with "LOG_" prefix
// - all flags with "log-" prefix
func AddFlags(v *viper.Viper, f *pflag.FlagSet) error {
	return config.AddFlags(&embedConfig{DefaultConfig()}, v, f, nil)
}

type EncoderConfig struct {
	MessageKey       *string          `key:"message-key,omitempty" env:"MESSAGE_KEY" flag:"message-key" desc:"Key for the log message (if empty, the message is omitted)"`
	LevelKey         *string          `key:"level-key,omitempty" env:"LEVEL_KEY" flag:"level-key" desc:"Key for the log level (if empty, the level is omitted)"`
	TimeKey          *string          `key:"time-key,omitempty" env:"TIME_KEY" flag:"time-key" desc:"Key for the log timestamp (if empty, the timestamp is omitted)"`
	NameKey          *string          `key:"name-key,omitempty" env:"NAME_KEY" flag:"name-key" desc:"Key for the log logger name (if empty, the logger name is omitted)"`
	CallerKey        *string          `key:"caller-key,omitempty" env:"CALLER_KEY" flag:"caller-key" desc:"Key for the log caller (if empty, the caller is omitted)"`
	FunctionKey      *string          `key:"function-key,omitempty" env:"FUNCTION_KEY" flag:"function-key" desc:"Key for the log function (if empty, the function is omitted)"`
	StacktraceKey    *string          `key:"stacktrace-key,omitempty" env:"STACKTRACE_KEY" flag:"stacktrace-key" desc:"Key for the log stacktrace (if empty, the stacktrace is omitted)"`
	SkipLineEnding   *bool            `key:"skip-line-ending,omitempty" env:"SKIP_LINE_ENDING" flag:"skip-line-ending" desc:"Skip the line ending"`
	LineEnding       *string          `key:"line-ending,omitempty" env:"LINE_ENDING" flag:"line-ending" desc:"Line ending"`
	LevelEncoder     *LevelEncoder    `key:"level-encoder,omitempty" env:"LEVEL_ENCODER" flag:"level-encoder" desc:"Primitive representation for the log level (e.g. 'capital', 'color', 'capitalColor', 'lowercase')"`
	TimeEncoder      *TimeEncoder     `key:"time-encoder,omitempty" env:"TIME_ENCODER" flag:"time-encoder" desc:"Primitive representation for the log timestamp (e.g. 'rfc3339nano', 'rfc3339', 'iso8601', 'millis', 'nanos', 'time')"`
	DurationEncoder  *DurationEncoder `key:"duration-encoder,omitempty" env:"DURATION_ENCODER" flag:"duration-encoder" desc:"Primitive representation for the log duration (e.g. 'string', 'nanos', 'ms', 's')"`
	CallerEncoder    *CallerEncoder   `key:"caller-encoder,omitempty" env:"CALLER_ENCODER" flag:"caller-encoder" desc:"Primitive representation for the log caller (e.g. 'full', 'short')"`
	NameEncoder      *NameEncoder     `key:"name-encoder,omitempty" env:"NAME_ENCODER" flag:"name-encoder" desc:"Primitive representation for the log logger name (e.g. 'full', 'short')"`
	ConsoleSeparator *string          `key:"console-separator,omitempty" env:"CONSOLE_SEPARATOR" flag:"console-separator" desc:"Field separator used by the console encoder"`
}

func (cfg *EncoderConfig) EncoderConfig() *zapcore.EncoderConfig {
	return &zapcore.EncoderConfig{
		MessageKey:       common.Val(cfg.MessageKey),
		LevelKey:         common.Val(cfg.LevelKey),
		TimeKey:          common.Val(cfg.TimeKey),
		NameKey:          common.Val(cfg.NameKey),
		CallerKey:        common.Val(cfg.CallerKey),
		FunctionKey:      common.Val(cfg.FunctionKey),
		StacktraceKey:    common.Val(cfg.StacktraceKey),
		SkipLineEnding:   common.Val(cfg.SkipLineEnding),
		LineEnding:       common.Val(cfg.LineEnding),
		ConsoleSeparator: common.Val(cfg.ConsoleSeparator),
		EncodeLevel:      zapLevelEncoders[common.Val(cfg.LevelEncoder)],
		EncodeTime:       zapTimeEncoders[common.Val(cfg.TimeEncoder)],
		EncodeDuration:   zapDurationEncoders[common.Val(cfg.DurationEncoder)],
		EncodeCaller:     zapCallerEncoders[common.Val(cfg.CallerEncoder)],
		EncodeName:       zapNameEncoders[common.Val(cfg.NameEncoder)],
	}
}

type SamplingConfig struct {
	Initial    *int `key:"initial,omitempty" desc:"Number of log entries with the same level and message to log before dropping entries"`
	Thereafter *int `key:"thereafter,omitempty" desc:"After the initial number of entries, every Mth entry is logged and the rest are dropped"`
}

var unknown = "unknown"

type Level int

func (l Level) String() string {
	if l >= 0 && l < Level(len(levelsStr)) {
		return levelsStr[l]
	}
	return unknown
}

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

var levelsStr = []string{
	"debug",
	"info",
	"warn",
	"error",
}

var levels = map[string]Level{
	levelsStr[DebugLevel]: DebugLevel,
	levelsStr[InfoLevel]:  InfoLevel,
	levelsStr[WarnLevel]:  WarnLevel,
	levelsStr[ErrorLevel]: ErrorLevel,
}

func ParseLevel(level string) (Level, error) {
	if l, ok := levels[strings.ToLower(level)]; ok {
		return l, nil
	}
	return 0, fmt.Errorf("invalid log level %q (must be one of %q)", level, levelsStr)
}

var zapLevels = map[Level]zapcore.Level{
	DebugLevel: zap.DebugLevel,
	InfoLevel:  zap.InfoLevel,
	WarnLevel:  zap.WarnLevel,
	ErrorLevel: zap.ErrorLevel,
}

type Format int

func (f Format) String() string {
	if f >= 0 && f < Format(len(formatsStr)) {
		return formatsStr[f]
	}
	return unknown
}

const (
	TextFormat Format = iota
	JSONFormat
)

var formatsStr = []string{
	"text",
	"json",
}

var formats = map[string]Format{
	formatsStr[TextFormat]: TextFormat,
	formatsStr[JSONFormat]: JSONFormat,
}

func ParseFormat(format string) (Format, error) {
	if f, ok := formats[strings.ToLower(format)]; ok {
		return f, nil
	}
	return 0, fmt.Errorf("invalid log format %q (must be one of %q)", format, formatsStr)
}

var zapFormats = map[Format]string{
	TextFormat: "console",
	JSONFormat: "json",
}

type LevelEncoder int

func (l LevelEncoder) String() string {
	if l >= 0 && l < LevelEncoder(len(levelEncodersStr)) {
		return levelEncodersStr[l]
	}
	return unknown
}

const (
	LevelEncoderCapital LevelEncoder = iota
	LevelEncoderCapitalColor
	LevelEncoderColor
	LevelEncoderLowercase
)

var levelEncodersStr = []string{
	"capital",
	"capitalColor",
	"color",
	"lowercase",
}

var levelEncoders = map[string]LevelEncoder{
	levelEncodersStr[LevelEncoderCapital]:      LevelEncoderCapital,
	levelEncodersStr[LevelEncoderCapitalColor]: LevelEncoderCapitalColor,
	levelEncodersStr[LevelEncoderColor]:        LevelEncoderColor,
	levelEncodersStr[LevelEncoderLowercase]:    LevelEncoderLowercase,
}

func ParseLevelEncoder(encoder string) (LevelEncoder, error) {
	if e, ok := levelEncoders[encoder]; ok {
		return e, nil
	}
	return 0, fmt.Errorf("invalid log level encoder %q (must be one of %q)", encoder, levelEncodersStr)
}

var zapLevelEncoders = map[LevelEncoder]zapcore.LevelEncoder{
	LevelEncoderCapital:      zapcore.CapitalLevelEncoder,
	LevelEncoderCapitalColor: zapcore.CapitalColorLevelEncoder,
	LevelEncoderColor:        zapcore.LowercaseColorLevelEncoder,
	LevelEncoderLowercase:    zapcore.LowercaseLevelEncoder,
}

type TimeEncoder int

func (t TimeEncoder) String() string {
	if t >= 0 && t < TimeEncoder(len(timeEncodersStr)) {
		return timeEncodersStr[t]
	}
	return unknown
}

const (
	TimeEncoderRFC3339Nano TimeEncoder = iota
	TimeEncoderRFC3339
	TimeEncoderISO8601
	TimeEncoderMillis
	TimeEncoderNanos
	TimeEncoderTime
)

var timeEncodersStr = []string{
	"rfc3339nano",
	"rfc3339",
	"iso8601",
	"millis",
	"nanos",
	"time",
}

var timeEncoders = map[string]TimeEncoder{
	timeEncodersStr[TimeEncoderRFC3339Nano]: TimeEncoderRFC3339Nano,
	timeEncodersStr[TimeEncoderRFC3339]:     TimeEncoderRFC3339,
	timeEncodersStr[TimeEncoderISO8601]:     TimeEncoderISO8601,
	timeEncodersStr[TimeEncoderMillis]:      TimeEncoderMillis,
	timeEncodersStr[TimeEncoderNanos]:       TimeEncoderNanos,
	timeEncodersStr[TimeEncoderTime]:        TimeEncoderTime,
}

func ParseTimeEncoder(encoder string) (TimeEncoder, error) {
	if e, ok := timeEncoders[strings.ToLower(encoder)]; ok {
		return e, nil
	}
	return 0, fmt.Errorf("invalid log time encoder %q (must be one of %q)", encoder, timeEncodersStr)
}

var zapTimeEncoders = map[TimeEncoder]zapcore.TimeEncoder{
	TimeEncoderRFC3339Nano: zapcore.RFC3339NanoTimeEncoder,
	TimeEncoderRFC3339:     zapcore.RFC3339TimeEncoder,
	TimeEncoderISO8601:     zapcore.ISO8601TimeEncoder,
	TimeEncoderMillis:      zapcore.EpochMillisTimeEncoder,
	TimeEncoderNanos:       zapcore.EpochNanosTimeEncoder,
	TimeEncoderTime:        zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000Z07:00"),
}

type DurationEncoder int

func (d DurationEncoder) String() string {
	if d >= 0 && d < DurationEncoder(len(durationEncodersStr)) {
		return durationEncodersStr[d]
	}
	return unknown
}

const (
	DurationEncoderString DurationEncoder = iota
	DurationEncoderNanos
	DurationEncoderMillis
	DurationEncoderSeconds
)

var durationEncodersStr = []string{
	"string",
	"nanos",
	"ms",
	"s",
}

var durationEncoders = map[string]DurationEncoder{
	durationEncodersStr[DurationEncoderString]:  DurationEncoderString,
	durationEncodersStr[DurationEncoderNanos]:   DurationEncoderNanos,
	durationEncodersStr[DurationEncoderMillis]:  DurationEncoderMillis,
	durationEncodersStr[DurationEncoderSeconds]: DurationEncoderSeconds,
}

func ParseDurationEncoder(encoder string) (DurationEncoder, error) {
	if e, ok := durationEncoders[strings.ToLower(encoder)]; ok {
		return e, nil
	}
	return 0, fmt.Errorf("invalid log duration encoder %q (must be one of %q)", encoder, durationEncodersStr)
}

var zapDurationEncoders = map[DurationEncoder]zapcore.DurationEncoder{
	DurationEncoderString:  zapcore.StringDurationEncoder,
	DurationEncoderNanos:   zapcore.NanosDurationEncoder,
	DurationEncoderMillis:  zapcore.MillisDurationEncoder,
	DurationEncoderSeconds: zapcore.SecondsDurationEncoder,
}

type CallerEncoder int

func (c CallerEncoder) String() string {
	if c >= 0 && c < CallerEncoder(len(callerEncodersStr)) {
		return callerEncodersStr[c]
	}
	return unknown
}

const (
	CallerEncoderFull CallerEncoder = iota
	CallerEncoderShort
)

var callerEncodersStr = []string{
	"full",
	"short",
}

var callerEncoders = map[string]CallerEncoder{
	callerEncodersStr[CallerEncoderFull]:  CallerEncoderFull,
	callerEncodersStr[CallerEncoderShort]: CallerEncoderShort,
}

func ParseCallerEncoder(encoder string) (CallerEncoder, error) {
	if e, ok := callerEncoders[strings.ToLower(encoder)]; ok {
		return e, nil
	}
	return 0, fmt.Errorf("invalid log caller encoder %q (must be one of %q)", encoder, callerEncodersStr)
}

var zapCallerEncoders = map[CallerEncoder]zapcore.CallerEncoder{
	CallerEncoderFull:  zapcore.FullCallerEncoder,
	CallerEncoderShort: zapcore.ShortCallerEncoder,
}

type NameEncoder int

func (n NameEncoder) String() string {
	if n >= 0 && n < NameEncoder(len(nameEncodersStr)) {
		return nameEncodersStr[n]
	}
	return unknown
}

const (
	NameEncoderFull NameEncoder = iota
	NameEncoderShort
)

var nameEncodersStr = []string{
	"full",
}

var nameEncoders = map[string]NameEncoder{
	nameEncodersStr[NameEncoderFull]: NameEncoderFull,
}

func ParseNameEncoder(encoder string) (NameEncoder, error) {
	if e, ok := nameEncoders[strings.ToLower(encoder)]; ok {
		return e, nil
	}
	return 0, fmt.Errorf("invalid log name encoder %q (must be one of %q)", encoder, nameEncodersStr)
}

var zapNameEncoders = map[NameEncoder]zapcore.NameEncoder{
	NameEncoderFull: zapcore.FullNameEncoder,
}
