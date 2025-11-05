package config

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/go-viper/mapstructure/v2"
)

var decodeHooks = []mapstructure.DecodeHookFunc{
	NilStrToNilDecodeHookFunc(nilStr),
	PtrToValueDecodeHookFunc(),
	StringToWeakSliceDecodeHookFunc(envSliceSep),
	StringToMapDecodeHookFunc(envSliceSep),
	mapstructure.StringToTimeDurationHookFunc(),
}

// RegisterGlobalDecodeHooks registers the given decode hooks to the global decode hooks
func RegisterGlobalDecodeHooks(hks ...mapstructure.DecodeHookFunc) {
	decodeHooks = append(decodeHooks, hks...)
}

// GlobalDecodeHook returns the global decode hook
func GlobalDecodeHook() mapstructure.DecodeHookFunc {
	return mapstructure.ComposeDecodeHookFunc(decodeHooks...)
}

// StringToWeakSliceHookFunc is a mapstructure decode hook that converts a string to a slice of strings.
// It is set by default in the viper.Viper package so we declare it here for custom viper.Viper instances.
func StringToWeakSliceDecodeHookFunc(sep string) mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data any,
	) (any, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
			return strings.Split(data.(string), sep), nil
		}

		return data, nil
	}
}

// StringToMapDecodeHookFunc is a mapstructure decode hook that converts a string to a map.
// The string format is "key1:value1 key2:value2" where sep is the separator between pairs (default " ").
// It also supports comma-separated pairs for pflag compatibility: "key1:value1,key2:value2".
func StringToMapDecodeHookFunc(sep string) mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data any,
	) (any, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t.Kind() != reflect.Map {
			return data, nil
		}

		str := data.(string)
		if str == "" {
			return reflect.MakeMap(t).Interface(), nil
		}

		// Split by separator to get key:value pairs
		pairs := strings.Split(str, sep)
		result := reflect.MakeMap(t)

		for _, pair := range pairs {
			// Split each pair by ":"
			parts := strings.SplitN(pair, ":", 2)
			if len(parts) != 2 {
				continue
			}

			key := reflect.ValueOf(parts[0])

			// Convert value to target type if needed
			var targetVal reflect.Value
			switch getKind(t.Elem().Kind()) {
			case reflect.String:
				targetVal = reflect.ValueOf(parts[1])
			case reflect.Int:
				intVal, err := strconv.ParseInt(parts[1], 10, 64)
				if err != nil {
					return data, err
				}
				// Convert to the exact int type
				switch t.Elem().Kind() {
				case reflect.Int:
					targetVal = reflect.ValueOf(int(intVal))
				case reflect.Int8:
					targetVal = reflect.ValueOf(int8(intVal))
				case reflect.Int16:
					targetVal = reflect.ValueOf(int16(intVal))
				case reflect.Int32:
					targetVal = reflect.ValueOf(int32(intVal))
				case reflect.Int64:
					targetVal = reflect.ValueOf(intVal)
				}
			case reflect.Uint:
				uintVal, err := strconv.ParseUint(parts[1], 10, 64)
				if err != nil {
					return data, err
				}
				// Convert to the exact uint type
				switch t.Elem().Kind() {
				case reflect.Uint:
					targetVal = reflect.ValueOf(uint(uintVal))
				case reflect.Uint8:
					targetVal = reflect.ValueOf(uint8(uintVal))
				case reflect.Uint16:
					targetVal = reflect.ValueOf(uint16(uintVal))
				case reflect.Uint32:
					targetVal = reflect.ValueOf(uint32(uintVal))
				case reflect.Uint64:
					targetVal = reflect.ValueOf(uintVal)
				}
			case reflect.Float32:
				floatVal, err := strconv.ParseFloat(parts[1], 64)
				if err != nil {
					return data, err
				}
				if t.Elem().Kind() == reflect.Float32 {
					targetVal = reflect.ValueOf(float32(floatVal))
				} else {
					targetVal = reflect.ValueOf(floatVal)
				}
			case reflect.Bool:
				boolVal, err := strconv.ParseBool(parts[1])
				if err != nil {
					return data, err
				}
				targetVal = reflect.ValueOf(boolVal)
			default:
				// For other types, just use the string value
				targetVal = reflect.ValueOf(parts[1])
			}

			result.SetMapIndex(key, targetVal)
		}

		return result.Interface(), nil
	}
}

// PtrToValueDecodeHookFunc is a mapstructure decode hook that converts a pointer to a value
func PtrToValueDecodeHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		_ reflect.Type,
		data any,
	) (any, error) {
		if f.Kind() != reflect.Ptr {
			return data, nil
		}

		kind := getKind(f.Elem().Kind())
		switch kind {
		case reflect.String,
			reflect.Bool,
			reflect.Int,
			reflect.Uint,
			reflect.Float32:
			if reflect.ValueOf(data).Elem().IsZero() {
				return data, nil
			}
			return reflect.ValueOf(data).Elem().Interface(), nil
		}
		return data, nil
	}
}

// NilStrToNilDecodeHookFunc is a mapstructure decode hook that converts a string to a nil
func NilStrToNilDecodeHookFunc(nilStr string) mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data any,
	) (any, error) {
		if f.Kind() == reflect.String && data.(string) == nilStr {
			return reflect.Zero(reflect.PointerTo(t)).Interface(), nil
		}

		return data, nil
	}
}
