package config

import (
	"fmt"
	"reflect"
)

var hooks = []EncodeHookFunc{
	StringerHook(),
}

func RegisterGlobalEncodeHooks(hks ...EncodeHookFunc) {
	hooks = append(hooks, hks...)
}

func GlobalEncodeHook() EncodeHookFunc {
	return CombineHooks(hooks...)
}

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

func StringerHook() EncodeHookFunc {
	return func(val reflect.Value) (reflect.Value, error) {
		if val.Type().Implements(reflect.TypeOf((*fmt.Stringer)(nil)).Elem()) && !val.IsNil() {
			return reflect.ValueOf(val.Interface().(fmt.Stringer).String()), nil
		}
		return val, nil
	}
}
