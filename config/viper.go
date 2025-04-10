package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/kkrt-labs/go-utils/common"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	envTagName   = "env"
	descTagName  = "desc"
	shortTagName = "short"
	flagTagName  = "flag"
	keyTagName   = "key"
	envSliceSep  = " " // We use space as a separator for env variables because it's used by viper when loading env variables
	nilStr       = "<nil>"
)

// NewViper creates a new viper instance with the given options.
func NewViper(opts ...viper.Option) *viper.Viper {
	v := viper.NewWithOptions(
		append([]viper.Option{viper.WithDecodeHook(GlobalDecodeHook())}, opts...)...,
	)
	v.AllowEmptyEnv(true)
	return v
}

// TagNameDecoderConfigOption is a viper.DecoderConfigOption that sets the tag name for the decoder.
func TagNameDecoderConfigOption(name string) viper.DecoderConfigOption {
	return func(cfg *mapstructure.DecoderConfig) {
		cfg.TagName = name
	}
}

// UnmarshalKey takes a single key and unmarshals it into a Struct.
func UnmarshalKey(key string, rawVal any, v *viper.Viper, opts ...viper.DecoderConfigOption) error {
	return v.UnmarshalKey(key, rawVal, append([]viper.DecoderConfigOption{TagNameDecoderConfigOption(keyTagName)}, opts...)...)
}

// Unmarshal unmarshals the config into a Struct. Make sure that the tags
// on the fields of the structure are properly set.
func Unmarshal(cfg any, v *viper.Viper, opts ...viper.DecoderConfigOption) error {
	return v.Unmarshal(cfg, append([]viper.DecoderConfigOption{TagNameDecoderConfigOption(keyTagName)}, opts...)...)
}

// Env returns a map of environment variables for the given configuration object.
// It uses the `env` struct field tag to get the environment variable name.
func Env(cfg any, hook func(reflect.Value) (reflect.Value, error)) (map[string]string, error) {
	if hook == nil {
		hook = GlobalEncodeHook()
	}

	inspector := NewInspector(&InspectorConfig{
		TagNames:   []string{envTagName},
		EncodeHook: hook,
	})

	infos, err := inspector.Inspect(cfg)
	if err != nil {
		return nil, fmt.Errorf("invalid configuration object: %w", err)
	}

	m := make(map[string]string)

	for _, info := range infos {
		envName := prepareEnvName(info.TagParts[envTagName], info.FieldParts)
		encoded := strings.TrimRight(strings.TrimLeft(info.Encoded, envSliceSep), envSliceSep)
		if encoded != "" {
			m[envName] = encoded
		}
	}

	return m, nil
}

// AddFlags adds flags to the given viper and flag set.
// It uses the `flag` struct field tag to get the flag name, environment variable name, description, and short flag.
func AddFlags(defaultConfig any, v *viper.Viper, f *pflag.FlagSet, hook func(reflect.Value) (reflect.Value, error)) error {
	if hook == nil {
		hook = GlobalEncodeHook()
	}

	inspector := NewInspector(&InspectorConfig{
		TagNames:   []string{keyTagName, flagTagName, envTagName, descTagName, shortTagName},
		EncodeHook: hook,
		IncludeNil: true,
	})

	infos, err := inspector.Inspect(defaultConfig)
	if err != nil {
		return fmt.Errorf("invalid configuration object: %w", err)
	}

loop:
	for field, info := range infos {
		viperKey := prepareViperKey(info.TagParts[keyTagName], info.FieldParts)
		flagName := prepareFlagName(info.TagParts[flagTagName], info.FieldParts)
		envName := prepareEnvName(info.TagParts[envTagName], info.FieldParts)

		desc := prepareDesc(info.TagParts[descTagName])

		short, err := prepareShort(info.TagParts[shortTagName])
		if err != nil {
			return err
		}

		kind := getKind(info.Processed.Kind())
		typ := info.Processed.Type()

		switch kind {
		case reflect.String:
			stringFlag := &StringFlag{
				ViperKey:    viperKey,
				Name:        flagName,
				Shorthand:   short,
				Env:         envName,
				Description: desc,
			}

			switch {
			case info.Value.Kind() == reflect.Ptr && info.Value.IsNil():
				stringFlag.DefaultValue = (*string)(nil)
			case info.Value.Kind() == reflect.Ptr:
				stringFlag.DefaultValue = common.Ptr(info.Processed.Interface().(string))
			default:
				stringFlag.DefaultValue = info.Processed.Interface().(string)
			}

			stringFlag.Add(v, f)
			continue loop
		case reflect.Bool:
			boolFlag := &BoolFlag{
				ViperKey:    viperKey,
				Name:        flagName,
				Shorthand:   short,
				Env:         envName,
				Description: desc,
			}

			switch {
			case info.Value.Kind() == reflect.Ptr && info.Value.IsNil():
				boolFlag.DefaultValue = (*bool)(nil)
			case info.Value.Kind() == reflect.Ptr:
				boolFlag.DefaultValue = common.Ptr(info.Processed.Interface().(bool))
			default:
				boolFlag.DefaultValue = info.Processed.Interface().(bool)
			}

			boolFlag.Add(v, f)
			continue loop
		case reflect.Int:
			intFlag := &IntFlag{
				ViperKey:    viperKey,
				Name:        flagName,
				Shorthand:   short,
				Env:         envName,
				Description: desc,
			}

			switch {
			case info.Value.Kind() == reflect.Ptr && info.Value.IsNil():
				intFlag.DefaultValue = (*int)(nil)
			case info.Value.Kind() == reflect.Ptr:
				intFlag.DefaultValue = common.Ptr(info.Processed.Interface().(int))
			default:
				intFlag.DefaultValue = info.Processed.Interface().(int)
			}

			intFlag.Add(v, f)
			continue loop
		case reflect.Array, reflect.Slice:
			elemKind := typ.Elem().Kind()
			switch elemKind {
			case reflect.String:
				stringArrayFlag := (&StringArrayFlag{
					ViperKey:    viperKey,
					Name:        flagName,
					Shorthand:   short,
					Env:         envName,
					Description: desc,
				})

				switch {
				case info.Value.Kind() == reflect.Ptr && info.Value.IsNil():
					stringArrayFlag.DefaultValue = (*[]string)(nil)
				case info.Value.Kind() == reflect.Ptr:
					defaultValue := make([]string, info.Processed.Len())
					reflect.Copy(reflect.ValueOf(defaultValue), info.Processed)
					stringArrayFlag.DefaultValue = common.Ptr(defaultValue)
				default:
					defaultValue := make([]string, info.Processed.Len())
					reflect.Copy(reflect.ValueOf(defaultValue), info.Processed)
					stringArrayFlag.DefaultValue = defaultValue
				}

				stringArrayFlag.Add(v, f)

				continue loop
			case reflect.Ptr:
				if typ.Elem().Elem().Kind() == reflect.String {
					stringArrayFlag := (&StringArrayFlag{
						ViperKey:    viperKey,
						Name:        flagName,
						Shorthand:   short,
						Env:         envName,
						Description: desc,
					})

					switch {
					case info.Value.Kind() == reflect.Ptr && info.Value.IsNil():
						stringArrayFlag.DefaultValue = (*[]string)(nil)
					case info.Value.Kind() == reflect.Ptr:
						stringArrayFlag.DefaultValue = common.Ptr(common.ValSlice(info.Processed.Interface().([]*string)...))
					default:
						stringArrayFlag.DefaultValue = common.ValSlice(info.Processed.Interface().([]*string)...)
					}

					stringArrayFlag.Add(v, f)

					continue loop
				}
			}
		}

		return fmt.Errorf("%s: unsupported type %s", field, typ)
	}

	return nil
}

func prepareViperKey(parts, fieldParts []string) string {
	for i, viperKeyPart := range parts {
		if viperKeyPart == "" {
			parts[i] = fieldParts[i]
		}
	}
	viperKey := strings.Join(parts, ".")
	viperKey = strings.ReplaceAll(viperKey, "_", ".")
	viperKey = strings.ReplaceAll(viperKey, " ", ".")
	return viperKey
}

func prepareFlagName(parts, fieldParts []string) string {
	for i, flagNamePart := range parts {
		if flagNamePart == "" {
			parts[i] = fieldParts[i]
		}
	}
	flagName := strings.Join(parts, "-")
	flagName = strings.ReplaceAll(flagName, ".", "-")
	flagName = strings.ReplaceAll(flagName, "_", "-")
	flagName = strings.ReplaceAll(flagName, " ", "-")
	flagName = strings.ToLower(flagName)
	return flagName
}

func prepareEnvName(envNameParts, fieldParts []string) string {
	for i, envNamePart := range envNameParts {
		if envNamePart == "" {
			envNameParts[i] = fieldParts[i]
		}
	}
	envName := strings.Join(envNameParts, "_")
	envName = strings.ReplaceAll(envName, ".", "_")
	envName = strings.ReplaceAll(envName, " ", "_")
	envName = strings.ToUpper(envName)
	return envName
}

func prepareDesc(parts []string) string {
	desc := strings.Join(parts, "")
	return desc
}

func prepareShort(parts []string) (string, error) {
	short := strings.Join(parts, "")
	if len(short) > 1 {
		return "", fmt.Errorf("invalid short flag must be a single character but got %q", short)
	}
	return short, nil
}
