package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-viper/mapstructure/v2"
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
	skipTag      = "-"
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
//
// It uses the following struct field tags
// - `flag` struct field tag to get the flag name, environment variable name, description, and short flag.
// - `env` struct field tag to get the environment variable name.
// - `desc` struct field tag to get the description.
// - `short` struct field tag to get the short flag.
// - `key` struct field tag to get the key.
//
// It sets the default values as set on the defaultConfig object.
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

	for field, info := range infos {
		viperKey := prepareViperKey(info.TagParts[keyTagName], info.FieldParts)
		flagName := prepareFlagName(info.TagParts[flagTagName], info.FieldParts)
		envName := prepareEnvName(info.TagParts[envTagName], info.FieldParts)

		desc := prepareDesc(info.TagParts[descTagName])
		usage := FlagDesc(desc, envName)

		short, err := prepareShort(info.TagParts[shortTagName])
		if err != nil {
			return err
		}

		switch info.Processed.Kind() {
		case reflect.String:
			f.StringP(flagName, short, info.Processed.Interface().(string), usage)
		case reflect.Bool:
			f.BoolP(flagName, short, info.Processed.Interface().(bool), usage)
		case reflect.Int:
			f.IntP(flagName, short, info.Processed.Interface().(int), usage)
		case reflect.Int8:
			f.Int8P(flagName, short, info.Processed.Interface().(int8), usage)
		case reflect.Int16:
			f.Int16P(flagName, short, info.Processed.Interface().(int16), usage)
		case reflect.Int32:
			f.Int32P(flagName, short, info.Processed.Interface().(int32), usage)
		case reflect.Int64:
			f.Int64P(flagName, short, info.Processed.Interface().(int64), usage)
		case reflect.Uint:
			f.UintP(flagName, short, info.Processed.Interface().(uint), usage)
		case reflect.Uint8:
			f.Uint8P(flagName, short, info.Processed.Interface().(uint8), usage)
		case reflect.Uint16:
			f.Uint16P(flagName, short, info.Processed.Interface().(uint16), usage)
		case reflect.Uint32:
			f.Uint32P(flagName, short, info.Processed.Interface().(uint32), usage)
		case reflect.Uint64:
			f.Uint64P(flagName, short, info.Processed.Interface().(uint64), usage)
		case reflect.Float32:
			f.Float32P(flagName, short, info.Processed.Interface().(float32), usage)
		case reflect.Float64:
			f.Float64P(flagName, short, info.Processed.Interface().(float64), usage)
		case reflect.Slice:
			elemKind := getKind(info.Processed.Type().Elem().Kind())
			switch elemKind {
			case reflect.Bool:
				f.BoolSliceP(flagName, short, info.Processed.Interface().([]bool), usage)
			case reflect.String:
				f.StringSliceP(flagName, short, info.Processed.Interface().([]string), usage)
			case reflect.Uint:
				slice := make([]uint, info.Processed.Len())
				for i := range info.Processed.Len() {
					slice[i] = uint(info.Processed.Index(i).Uint())
				}
				f.UintSliceP(flagName, short, slice, usage)
			case reflect.Int:
				slice := make([]int, info.Processed.Len())
				for i := range info.Processed.Len() {
					slice[i] = int(info.Processed.Index(i).Int())
				}
				f.IntSliceP(flagName, short, slice, usage)
			case reflect.Float32:
				slice := make([]float64, info.Processed.Len())
				for i := range info.Processed.Len() {
					slice[i] = info.Processed.Index(i).Float()
				}
				f.Float64SliceP(flagName, short, slice, usage)
			default:
				return fmt.Errorf("%v: unsupported slice element type %s", field, elemKind)
			}
		case reflect.Array:
			value := reflect.MakeSlice(reflect.SliceOf(info.Processed.Type().Elem()), info.Processed.Len(), info.Processed.Len())
			reflect.Copy(value, info.Processed)
			elemKind := getKind(info.Processed.Type().Elem().Kind())
			switch elemKind {
			case reflect.Bool:
				f.BoolSlice(flagName, value.Interface().([]bool), usage)
			case reflect.String:
				f.StringSlice(flagName, value.Interface().([]string), usage)
			case reflect.Uint:
				slice := make([]uint, info.Processed.Len())
				for i := range info.Processed.Len() {
					slice[i] = uint(info.Processed.Index(i).Uint())
				}
				f.UintSlice(flagName, slice, usage)
			case reflect.Int:
				slice := make([]int, info.Processed.Len())
				for i := range info.Processed.Len() {
					slice[i] = int(info.Processed.Index(i).Int())
				}
				f.IntSlice(flagName, slice, usage)
			case reflect.Float32:
				slice := make([]float64, info.Processed.Len())
				for i := range info.Processed.Len() {
					slice[i] = info.Processed.Index(i).Float()
				}
				f.Float64Slice(flagName, slice, usage)
			default:
				return fmt.Errorf("%v: unsupported array element type %s", field, elemKind)
			}
		default:
			return fmt.Errorf("%v: unsupported type %s", field, info.Processed.Kind())
		}

		_ = v.BindPFlag(viperKey, f.Lookup(flagName))
		if envName != "" {
			_ = v.BindEnv(viperKey, envName)
		}
		v.SetDefault(viperKey, info.Value.Interface())
	}

	return nil
}

func prepareViperKey(parts, fieldParts []string) string {
	finalParts := make([]string, 0, len(parts))
	for i, viperKeyPart := range parts {
		switch viperKeyPart {
		case "":
			finalParts = append(finalParts, fieldParts[i])
		case skipTag:
			continue
		default:
			finalParts = append(finalParts, viperKeyPart)
		}
	}
	viperKey := strings.Join(finalParts, ".")
	viperKey = strings.ReplaceAll(viperKey, "_", ".")
	viperKey = strings.ReplaceAll(viperKey, " ", ".")
	return viperKey
}

func prepareFlagName(parts, fieldParts []string) string {
	finalParts := make([]string, 0, len(parts))
	for i, flagNamePart := range parts {
		switch flagNamePart {
		case "":
			finalParts = append(finalParts, fieldParts[i])
		case skipTag:
			continue
		default:
			finalParts = append(finalParts, flagNamePart)
		}
	}
	flagName := strings.Join(finalParts, "-")
	flagName = strings.ReplaceAll(flagName, ".", "-")
	flagName = strings.ReplaceAll(flagName, "_", "-")
	flagName = strings.ReplaceAll(flagName, " ", "-")
	flagName = strings.ToLower(flagName)
	return flagName
}

func prepareEnvName(envNameParts, fieldParts []string) string {
	finalParts := make([]string, 0, len(envNameParts))
	for i, envNamePart := range envNameParts {
		switch envNamePart {
		case "":
			finalParts = append(finalParts, fieldParts[i])
		case skipTag:
			continue
		default:
			finalParts = append(finalParts, envNamePart)
		}
	}
	envName := strings.Join(finalParts, "_")
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
