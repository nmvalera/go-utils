package config

import (
	"reflect"
	"strings"

	"github.com/go-viper/mapstructure/v2"
)

var decodeHooks = []mapstructure.DecodeHookFunc{
	NilStrToNilDecodeHookFunc(nilStr),
	PtrToValueDecodeHookFunc(),
	StringToWeakSliceDecodeHookFunc(envSliceSep),
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
