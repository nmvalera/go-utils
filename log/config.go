package log

import (
	"fmt"
	"strings"

	"github.com/kkrt-labs/go-utils/common"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level int

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

type EncoderConfig struct {
	MessageKey       *string `mapstructure:"message-key,omitempty"`
	LevelKey         *string `mapstructure:"level-key,omitempty"`
	TimeKey          *string `mapstructure:"time-key,omitempty"`
	NameKey          *string `mapstructure:"name-key,omitempty"`
	CallerKey        *string `mapstructure:"caller-key,omitempty"`
	FunctionKey      *string `mapstructure:"function-key,omitempty"`
	StacktraceKey    *string `mapstructure:"stacktrace-key,omitempty"`
	SkipLineEnding   *bool   `mapstructure:"skip-line-ending,omitempty"`
	LineEnding       *string `mapstructure:"line-ending,omitempty"`
	LevelEncoder     *string `mapstructure:"level-encoder,omitempty"`
	TimeEncoder      *string `mapstructure:"time-encoder,omitempty"`
	DurationEncoder  *string `mapstructure:"duration-encoder,omitempty"`
	CallerEncoder    *string `mapstructure:"caller-encoder,omitempty"`
	NameEncoder      *string `mapstructure:"name-encoder,omitempty"`
	ConsoleSeparator *string `mapstructure:"console-separator,omitempty"`
}

func (cfg *EncoderConfig) SetDefaults() *EncoderConfig {
	if cfg.MessageKey == nil {
		cfg.MessageKey = common.Copy(defaultConfig.Encoder.MessageKey)
	}
	if cfg.LevelKey == nil {
		cfg.LevelKey = common.Copy(defaultConfig.Encoder.LevelKey)
	}
	if cfg.TimeKey == nil {
		cfg.TimeKey = common.Copy(defaultConfig.Encoder.TimeKey)
	}
	if cfg.NameKey == nil {
		cfg.NameKey = common.Copy(defaultConfig.Encoder.NameKey)
	}
	if cfg.CallerKey == nil {
		cfg.CallerKey = common.Copy(defaultConfig.Encoder.CallerKey)
	}
	if cfg.FunctionKey == nil {
		cfg.FunctionKey = common.Copy(defaultConfig.Encoder.FunctionKey)
	}
	if cfg.StacktraceKey == nil {
		cfg.StacktraceKey = common.Copy(defaultConfig.Encoder.StacktraceKey)
	}
	if cfg.SkipLineEnding == nil {
		cfg.SkipLineEnding = common.Copy(defaultConfig.Encoder.SkipLineEnding)
	}
	if cfg.LineEnding == nil {
		cfg.LineEnding = common.Copy(defaultConfig.Encoder.LineEnding)
	}
	if cfg.LevelEncoder == nil {
		cfg.LevelEncoder = common.Copy(defaultConfig.Encoder.LevelEncoder)
	}
	if cfg.TimeEncoder == nil {
		cfg.TimeEncoder = common.Copy(defaultConfig.Encoder.TimeEncoder)
	}
	if cfg.DurationEncoder == nil {
		cfg.DurationEncoder = common.Copy(defaultConfig.Encoder.DurationEncoder)
	}
	if cfg.CallerEncoder == nil {
		cfg.CallerEncoder = common.Copy(defaultConfig.Encoder.CallerEncoder)
	}
	if cfg.NameEncoder == nil {
		cfg.NameEncoder = common.Copy(defaultConfig.Encoder.NameEncoder)
	}
	if cfg.ConsoleSeparator == nil {
		cfg.ConsoleSeparator = common.Copy(defaultConfig.Encoder.ConsoleSeparator)
	}
	return cfg
}

func ParseEncoderConfig(cfg *EncoderConfig) (*zapcore.EncoderConfig, error) {
	zapCfg := zapcore.EncoderConfig{
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
	}

	// LevelEncoder
	levelEncoder, err := ParseLevelEncoder(common.Val(cfg.LevelEncoder))
	if err != nil {
		return nil, err
	}
	zapCfg.EncodeLevel = zapLevelEncoders[levelEncoder]

	// TimeEncoder
	timeEncoder, err := ParseTimeEncoder(common.Val(cfg.TimeEncoder))
	if err != nil {
		return nil, err
	}
	zapCfg.EncodeTime = zapTimeEncoders[timeEncoder]

	// DurationEncoder
	durationEncoder, err := ParseDurationEncoder(common.Val(cfg.DurationEncoder))
	if err != nil {
		return nil, err
	}
	zapCfg.EncodeDuration = zapDurationEncoders[durationEncoder]

	// CallerEncoder
	callerEncoder, err := ParseCallerEncoder(common.Val(cfg.CallerEncoder))
	if err != nil {
		return nil, err
	}
	zapCfg.EncodeCaller = zapCallerEncoders[callerEncoder]

	// NameEncoder
	nameEncoder, err := ParseNameEncoder(common.Val(cfg.NameEncoder))
	if err != nil {
		return nil, err
	}
	zapCfg.EncodeName = zapNameEncoders[nameEncoder]

	return &zapCfg, nil
}

type SamplingConfig struct {
	Initial    *int `mapstructure:"initial,omitempty"`
	Thereafter *int `mapstructure:"thereafter,omitempty"`
}

func (cfg *SamplingConfig) SetDefaults() *SamplingConfig {
	if cfg.Initial == nil {
		cfg.Initial = common.Copy(defaultConfig.Sampling.Initial)
	}
	if cfg.Thereafter == nil {
		cfg.Thereafter = common.Copy(defaultConfig.Sampling.Thereafter)
	}
	return cfg
}

// Config is the configuration for the logger.
// It can be used in conjunction with viper.
//
// Example:
//
//	 import "github.com/kkrt-labs/go-utils/log"
//
//	 cfg := &log.Config{
//			Log log.Config `mapstructure:"log"`
//	 }
//
//	 viper.Set("log.format", "json")
//	 viper.Set("log.level", "info")
//	 viper.Set("log.enable-stacktrace", true)
//	 viper.Set("log.enable-caller", true)
//	 viper.Set("log.encoder.message-key", "msg")
//
//	 var cfg Config
//	 err := viper.Unmarshal(&cfg)
//	 if err != nil {
//	 	fmt.Printf("Failed to unmarshal config: %v", err)
//	 } else {
//	 	fmt.Printf("Config: %+v", cfg)
//	 }
type Config struct {
	Format           *string         `mapstructure:"format,omitempty"`
	Level            *string         `mapstructure:"level,omitempty"`
	EnableStacktrace *bool           `mapstructure:"enable-stacktrace,omitempty"`
	EnableCaller     *bool           `mapstructure:"enable-caller,omitempty"`
	Encoder          *EncoderConfig  `mapstructure:"encoder,omitempty"`
	Sampling         *SamplingConfig `mapstructure:"sampling,omitempty"`
	OutputPaths      *[]*string      `mapstructure:"output-paths,omitempty"`
	ErrorOutputPaths *[]*string      `mapstructure:"error-output-paths,omitempty"`
}

func (cfg *Config) SetDefaults() *Config {
	if cfg.Format == nil {
		cfg.Format = common.Copy(defaultConfig.Format)
	}
	if cfg.Level == nil {
		cfg.Level = common.Copy(defaultConfig.Level)
	}
	if cfg.EnableStacktrace == nil {
		cfg.EnableStacktrace = common.Copy(defaultConfig.EnableStacktrace)
	}
	if cfg.EnableCaller == nil {
		cfg.EnableCaller = common.Copy(defaultConfig.EnableCaller)
	}
	if cfg.Encoder == nil {
		cfg.Encoder = &EncoderConfig{}
	}
	cfg.Encoder.SetDefaults()

	if cfg.Sampling == nil {
		cfg.Sampling = new(SamplingConfig)
	}
	cfg.Sampling.SetDefaults()

	if cfg.OutputPaths == nil {
		cfg.OutputPaths = common.Ptr(common.CopySlice(*defaultConfig.OutputPaths...))
	}

	if cfg.ErrorOutputPaths == nil {
		cfg.ErrorOutputPaths = common.Ptr(common.CopySlice(*defaultConfig.ErrorOutputPaths...))
	}

	return cfg
}

func (cfg *Config) Load(v *viper.Viper) error {
	type embedConfig struct {
		Log *Config `mapstructure:"log"`
	}
	AddFlags(v, new(pflag.FlagSet))
	return v.Unmarshal(&embedConfig{cfg})
}

func (cfg *Config) Env() map[string]*string {
	m := make(map[string]*string)
	if cfg.Level != nil {
		m[logLevelFlag.Env] = cfg.Level
	}
	if cfg.Format != nil {
		m[logFormatFlag.Env] = cfg.Format
	}
	if cfg.EnableStacktrace != nil {
		m[logEnableStacktraceFlag.Env] = common.Ptr(fmt.Sprintf("%t", *cfg.EnableStacktrace))
	}
	if cfg.EnableCaller != nil {
		m[logEnableCallerFlag.Env] = common.Ptr(fmt.Sprintf("%t", *cfg.EnableCaller))
	}
	if cfg.Encoder != nil {
		if cfg.Encoder.MessageKey != nil {
			m[logEncoderMessageKeyFlag.Env] = cfg.Encoder.MessageKey
		}
		if cfg.Encoder.LevelKey != nil {
			m[logEncoderLevelKeyFlag.Env] = cfg.Encoder.LevelKey
		}
		if cfg.Encoder.TimeKey != nil {
			m[logEncoderTimeKeyFlag.Env] = cfg.Encoder.TimeKey
		}
		if cfg.Encoder.NameKey != nil {
			m[logEncoderNameKeyFlag.Env] = cfg.Encoder.NameKey
		}
		if cfg.Encoder.CallerKey != nil {
			m[logEncoderCallerKeyFlag.Env] = cfg.Encoder.CallerKey
		}
		if cfg.Encoder.FunctionKey != nil {
			m[logEncoderFunctionKeyFlag.Env] = cfg.Encoder.FunctionKey
		}
		if cfg.Encoder.StacktraceKey != nil {
			m[logEncoderStacktraceKeyFlag.Env] = cfg.Encoder.StacktraceKey
		}
		if cfg.Encoder.SkipLineEnding != nil {
			m[logEncoderSkipLineEndingFlag.Env] = common.Ptr(fmt.Sprintf("%t", *cfg.Encoder.SkipLineEnding))
		}
		if cfg.Encoder.LineEnding != nil {
			m[logEncoderLineEndingFlag.Env] = cfg.Encoder.LineEnding
		}
		if cfg.Encoder.LevelEncoder != nil {
			m[logEncoderLevelEncoderFlag.Env] = cfg.Encoder.LevelEncoder
		}
		if cfg.Encoder.TimeEncoder != nil {
			m[logEncoderTimeEncoderFlag.Env] = cfg.Encoder.TimeEncoder
		}
		if cfg.Encoder.DurationEncoder != nil {
			m[logEncoderDurationEncoderFlag.Env] = cfg.Encoder.DurationEncoder
		}
		if cfg.Encoder.CallerEncoder != nil {
			m[logEncoderCallerEncoderFlag.Env] = cfg.Encoder.CallerEncoder
		}
		if cfg.Encoder.NameEncoder != nil {
			m[logEncoderNameEncoderFlag.Env] = cfg.Encoder.NameEncoder
		}
		if cfg.Encoder.ConsoleSeparator != nil {
			m[logEncoderConsoleSeparatorFlag.Env] = cfg.Encoder.ConsoleSeparator
		}
	}

	if cfg.Sampling != nil {
		if cfg.Sampling.Initial != nil {
			m[logSamplingInitialFlag.Env] = common.Ptr(fmt.Sprintf("%d", *cfg.Sampling.Initial))
		}
		if cfg.Sampling.Thereafter != nil {
			m[logSamplingThereafterFlag.Env] = common.Ptr(fmt.Sprintf("%d", *cfg.Sampling.Thereafter))
		}
	}

	if cfg.OutputPaths != nil {
		m[logOutputPathsFlag.Env] = common.Ptr(strings.Join(common.ValSlice(*cfg.OutputPaths...), ","))
	}

	if cfg.ErrorOutputPaths != nil {
		m[logErrorOutputPathsFlag.Env] = common.Ptr(strings.Join(common.ValSlice(*cfg.ErrorOutputPaths...), ","))
	}

	return m
}

func ParseConfig(cfg *Config) (*zap.Config, error) {
	zapCfg := &zap.Config{
		DisableStacktrace: !common.Val(cfg.EnableStacktrace),
		DisableCaller:     !common.Val(cfg.EnableCaller),
		Level:             zap.NewAtomicLevel(),
		Sampling: &zap.SamplingConfig{
			Initial:    common.Val(cfg.Sampling.Initial),
			Thereafter: common.Val(cfg.Sampling.Thereafter),
		},
		OutputPaths:      common.ValSlice(common.Val(cfg.OutputPaths)...),
		ErrorOutputPaths: common.ValSlice(common.Val(cfg.ErrorOutputPaths)...),
	}

	// Log Level
	level, err := ParseLevel(common.Val(cfg.Level))
	if err != nil {
		return nil, err
	}
	zapCfg.Level.SetLevel(zapLevels[level])

	// Log Format
	format, err := ParseFormat(common.Val(cfg.Format))
	if err != nil {
		return nil, err
	}
	zapCfg.Encoding = zapFormats[format]

	// Encoder Config
	encoderCfg, err := ParseEncoderConfig(cfg.Encoder)
	if err != nil {
		return nil, err
	}
	zapCfg.EncoderConfig = *encoderCfg

	return zapCfg, nil
}

var defaultConfig = &Config{
	Level:            common.Ptr(levelsStr[InfoLevel]),
	Format:           common.Ptr(formatsStr[TextFormat]),
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
		LevelEncoder:     common.Ptr(levelEncodersStr[LevelEncoderCapitalColor]),
		TimeEncoder:      common.Ptr(timeEncodersStr[TimeEncoderRFC3339]),
		DurationEncoder:  common.Ptr(durationEncodersStr[DurationEncoderSeconds]),
		CallerEncoder:    common.Ptr(callerEncodersStr[CallerEncoderShort]),
		NameEncoder:      common.Ptr(nameEncodersStr[NameEncoderFull]),
		ConsoleSeparator: common.Ptr("\t"),
	},
	Sampling: &SamplingConfig{
		Initial:    common.Ptr(100),
		Thereafter: common.Ptr(100),
	},
	OutputPaths:      common.PtrSlice("stderr"),
	ErrorOutputPaths: common.PtrSlice("stderr"),
}
