package config

import (
	"fmt"
	"reflect"
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
// It supports nested maps using additional colon-separated keys: "key:nestedKey:value".
// The nesting depth is determined by the target type (e.g. map[string]map[string]int has depth 2).
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

		// Parse string into a nested map[string]any with string leaves
		segments := strings.Split(str, sep)
		intermediate := make(map[string]any)
		for _, segment := range segments {
			parts := strings.Split(segment, ":")
			if len(parts) < 2 {
				return data, fmt.Errorf("invalid segment: %s", segment)
			}
			setNestedMapValue(intermediate, parts[:len(parts)-1], parts[len(parts)-1])
		}

		// Let mapstructure handle all type conversions
		targetPtr := reflect.New(t)
		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			WeaklyTypedInput: true,
			Result:           targetPtr.Interface(),
		})
		if err != nil {
			return data, err
		}
		if err := decoder.Decode(intermediate); err != nil {
			return data, err
		}

		return targetPtr.Elem().Interface(), nil
	}
}

// setNestedMapValue sets a value in a nested map[string]any structure.
// keys is the path of map keys, value is the leaf string value.
func setNestedMapValue(m map[string]any, keys []string, value string) {
	if len(keys) == 1 {
		m[keys[0]] = value
		return
	}
	sub, ok := m[keys[0]].(map[string]any)
	if !ok {
		sub = make(map[string]any)
		m[keys[0]] = sub
	}
	setNestedMapValue(sub, keys[1:], value)
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
