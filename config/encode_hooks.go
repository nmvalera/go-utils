package config

import (
	"fmt"
	"reflect"
)

var hooks = []EncodeHookFunc{
	StringerHook(),
}

// RegisterGlobalEncodeHooks registers the given encode hooks to the global encode hooks
func RegisterGlobalEncodeHooks(hks ...EncodeHookFunc) {
	hooks = append(hooks, hks...)
}

// GlobalEncodeHook returns the global encode hook
func GlobalEncodeHook() EncodeHookFunc {
	return CombineHooks(hooks...)
}

// CombineHooks combines the given encode hooks into a single encode hook
func CombineHooks(hooks ...EncodeHookFunc) EncodeHookFunc {
	return func(val reflect.Value) (reflect.Value, error) {
		var err error
		for _, hook := range hooks {
			val, err = hook(val)
			if err != nil {
				return reflect.Value{}, err
			}
		}
		return val, nil
	}
}

// StringerHook is a mapstructure encode hook that converts a value to a string
func StringerHook() EncodeHookFunc {
	return func(val reflect.Value) (reflect.Value, error) {
		if val.Type().Implements(reflect.TypeOf((*fmt.Stringer)(nil)).Elem()) {
			switch {
			case val.Kind() != reflect.Ptr,
				val.Kind() == reflect.Ptr && !val.IsNil():
				return reflect.ValueOf(val.Interface().(fmt.Stringer).String()), nil
			}
		}
		return val, nil
	}
}
