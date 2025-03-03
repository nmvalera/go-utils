# Logging Configuration

This package provides a flexible logging configuration system built on top of [Uber's zap logging library](https://github.com/uber-go/zap).

## Configuration Options

The logging can be configured through flags or environment variables. Here are all available options:

### Basic Configuration

#### Log Level

Controls the minimum logging level.

- Flag: `--log-level`
- Environment: `LOG_LEVEL`
- Values: `debug`, `info`, `warn`, `error`
- Default: `info`

#### Log Format 

Controls the output format of the logs.

- Flag: `--log-format`
- Environment: `LOG_FORMAT` 
- Values: `text`, `json`
- Default: `text`

#### Stack Traces

Controls whether stack traces are included for errors.

- Flag: `--log-enable-stacktrace`
- Environment: `LOG_ENABLE_STACKTRACE`
- Values: `true`, `false`
- Default: `false`

#### Caller Information

Controls whether to include the caller (file and line) in log entries.

- Flag: `--log-enable-caller`
- Environment: `LOG_ENABLE_CALLER`
- Values: `true`, `false`
- Default: `false`

## Advanced Configuration

### Key Names

Configure the keys used in structured logging:

- Message Key
  - Flag: `--log-encoder-message-key`
  - Environment: `LOG_ENCODER_MESSAGE_KEY`
  - Default: `msg`

- Level Key
  - Flag: `--log-encoder-level-key`
  - Environment: `LOG_ENCODER_LEVEL_KEY`
  - Default: `level`

- Time Key
  - Flag: `--log-encoder-time-key`
  - Environment: `LOG_ENCODER_TIME_KEY`
  - Default: `ts`

- Name Key
  - Flag: `--log-encoder-name-key`
  - Environment: `LOG_ENCODER_NAME_KEY`
  - Default: `logger`

- Caller Key
  - Flag: `--log-encoder-caller-key`
  - Environment: `LOG_ENCODER_CALLER_KEY`
  - Default: `caller`

- Function Key
  - Flag: `--log-encoder-function-key`
  - Environment: `LOG_ENCODER_FUNCTION_KEY`
  - Default: `` (empty)

- Stacktrace Key
  - Flag: `--log-encoder-stacktrace-key`
  - Environment: `LOG_ENCODER_STACKTRACE_KEY`
  - Default: `stacktrace`

### Encoder Configuration

#### Level Encoder

Controls how the log level is formatted.

- Flag: `--log-encoder-level-encoder`
- Environment: `LOG_ENCODER_LEVEL_ENCODER`
- Values:
  - `capital` - "INFO", "ERROR"
  - `capitalColor` - "INFO" (with color)
  - `color` - "info" (with color)
  - `lowercase` - "info", "error"
- Default: `capitalColor`

#### Time Encoder

Controls how timestamps are formatted.

- Flag: `--log-encoder-time-encoder`
- Environment: `LOG_ENCODER_TIME_ENCODER`
- Values:
  - `rfc3339nano` - RFC3339 with nanoseconds
  - `rfc3339` - RFC3339
  - `iso8601` - ISO8601
  - `millis` - Milliseconds since epoch
  - `nanos` - Nanoseconds since epoch
  - `time` - Custom format "2006-01-02T15:04:05.000Z07:00"
- Default: `time`

#### Duration Encoder

Controls how durations are formatted.

- Flag: `--log-encoder-duration-encoder`
- Environment: `LOG_ENCODER_DURATION_ENCODER`
- Values:
  - `string` - Human-readable duration (e.g. "1s")
  - `nanos` - Nanoseconds as integer
  - `ms` - Milliseconds as integer
  - `s` - Seconds as float
- Default: `s`

#### Caller Encoder

Controls how caller information is formatted.

- Flag: `--log-encoder-caller-encoder`
- Environment: `LOG_ENCODER_CALLER_ENCODER`
- Values:
  - `full` - Full file path
  - `short` - Package name and file
- Default: `short`

#### Name Encoder

Controls how logger names are formatted.

- Flag: `--log-encoder-name-encoder`
- Environment: `LOG_ENCODER_NAME_ENCODER`
- Values:
  - `full` - Full logger name
- Default: `full`

### Output Formatting

#### Line Ending

- Skip Line Ending
  - Flag: `--log-encoder-skip-line-ending`
  - Environment: `LOG_ENCODER_SKIP_LINE_ENDING`
  - Values: `true`, `false`
  - Default: `false`

- Line Ending
  - Flag: `--log-encoder-line-ending`
  - Environment: `LOG_ENCODER_LINE_ENDING`
  - Default: `\n`

#### Console Separator

Controls the separator used in console output format.

- Flag: `--log-encoder-console-separator`
- Environment: `LOG_ENCODER_CONSOLE_SEPARATOR`
- Default: `\t` (tab)

## Example Usage

Using flags:

```bash
./myapp --log-level debug --log-format json --log-enable-stacktrace true --log-enable-caller true
```

Using environment variables:

```bash
LOG_LEVEL=debug LOG_FORMAT=json LOG_ENABLE_STACKTRACE=true LOG_ENABLE_CALLER=true ./myapp
```

Advanced configuration example:

```bash
./myapp \
  --log-level debug \
  --log-format json \
  --log-encoder-time-encoder rfc3339nano \
  --log-encoder-level-encoder capitalColor \
  --log-encoder-caller-encoder full \
  --log-encoder-console-separator "  "
```

