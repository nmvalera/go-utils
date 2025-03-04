package log

import (
	"fmt"
	"strings"

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
	MessageKey       string `mapstructure:"message-key"`
	LevelKey         string `mapstructure:"level-key"`
	TimeKey          string `mapstructure:"time-key"`
	NameKey          string `mapstructure:"name-key"`
	CallerKey        string `mapstructure:"caller-key"`
	FunctionKey      string `mapstructure:"function-key"`
	StacktraceKey    string `mapstructure:"stacktrace-key"`
	SkipLineEnding   bool   `mapstructure:"skip-line-ending"`
	LineEnding       string `mapstructure:"line-ending"`
	LevelEncoder     string `mapstructure:"level-encoder"`
	TimeEncoder      string `mapstructure:"time-encoder"`
	DurationEncoder  string `mapstructure:"duration-encoder"`
	CallerEncoder    string `mapstructure:"caller-encoder"`
	NameEncoder      string `mapstructure:"name-encoder"`
	ConsoleSeparator string `mapstructure:"console-separator"`
}

func ParseEncoderConfig(cfg *EncoderConfig) (*zapcore.EncoderConfig, error) {
	zapCfg := zapcore.EncoderConfig{
		MessageKey:       cfg.MessageKey,
		LevelKey:         cfg.LevelKey,
		TimeKey:          cfg.TimeKey,
		NameKey:          cfg.NameKey,
		CallerKey:        cfg.CallerKey,
		FunctionKey:      cfg.FunctionKey,
		StacktraceKey:    cfg.StacktraceKey,
		SkipLineEnding:   cfg.SkipLineEnding,
		LineEnding:       cfg.LineEnding,
		ConsoleSeparator: cfg.ConsoleSeparator,
	}

	// LevelEncoder
	if cfg.LevelEncoder != "" {
		levelEncoder, err := ParseLevelEncoder(cfg.LevelEncoder)
		if err != nil {
			return nil, err
		}
		zapCfg.EncodeLevel = zapLevelEncoders[levelEncoder]
	}

	// TimeEncoder
	if cfg.TimeEncoder != "" {
		timeEncoder, err := ParseTimeEncoder(cfg.TimeEncoder)
		if err != nil {
			return nil, err
		}
		zapCfg.EncodeTime = zapTimeEncoders[timeEncoder]
	}

	// DurationEncoder
	if cfg.DurationEncoder != "" {
		durationEncoder, err := ParseDurationEncoder(cfg.DurationEncoder)
		if err != nil {
			return nil, err
		}
		zapCfg.EncodeDuration = zapDurationEncoders[durationEncoder]
	}

	// CallerEncoder
	if cfg.CallerEncoder != "" {
		callerEncoder, err := ParseCallerEncoder(cfg.CallerEncoder)
		if err != nil {
			return nil, err
		}
		zapCfg.EncodeCaller = zapCallerEncoders[callerEncoder]
	}

	// NameEncoder
	if cfg.NameEncoder != "" {
		nameEncoder, err := ParseNameEncoder(cfg.NameEncoder)
		if err != nil {
			return nil, err
		}
		zapCfg.EncodeName = zapNameEncoders[nameEncoder]
	}

	return &zapCfg, nil
}

type SamplingConfig struct {
	Initial    int `mapstructure:"initial"`
	Thereafter int `mapstructure:"thereafter"`
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
	Format           string         `mapstructure:"format"`
	Level            string         `mapstructure:"level"`
	EnableStacktrace bool           `mapstructure:"enable-stacktrace"`
	EnableCaller     bool           `mapstructure:"enable-caller"`
	Encoder          EncoderConfig  `mapstructure:"encoder"`
	Sampling         SamplingConfig `mapstructure:"sampling"`
	OutputPaths      []string       `mapstructure:"output-paths"`
	ErrorOutputPaths []string       `mapstructure:"error-output-paths"`
}

func ParseConfig(cfg *Config) (*zap.Config, error) {
	zapCfg := &zap.Config{
		DisableStacktrace: !cfg.EnableStacktrace,
		DisableCaller:     !cfg.EnableCaller,
		Level:             zap.NewAtomicLevel(),
		Sampling: &zap.SamplingConfig{
			Initial:    cfg.Sampling.Initial,
			Thereafter: cfg.Sampling.Thereafter,
		},
	}

	// Log Level
	level, err := ParseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}
	zapCfg.Level.SetLevel(zapLevels[level])

	// Log Format
	format, err := ParseFormat(cfg.Format)
	if err != nil {
		return nil, err
	}
	zapCfg.Encoding = zapFormats[format]

	// Encoder Config
	encoderCfg, err := ParseEncoderConfig(&cfg.Encoder)
	if err != nil {
		return nil, err
	}
	zapCfg.EncoderConfig = *encoderCfg

	return zapCfg, nil
}
