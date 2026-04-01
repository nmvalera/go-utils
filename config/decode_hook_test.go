package config

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringToMapDecodeHookFunc(t *testing.T) {
	hook := StringToMapDecodeHookFunc(" ")

	decode := func(t *testing.T, input any, target any) any {
		t.Helper()
		f := reflect.TypeOf(input)
		to := reflect.TypeOf(target)
		result, err := hook.(func(reflect.Type, reflect.Type, any) (any, error))(f, to, input)
		require.NoError(t, err)
		return result
	}

	decodeErr := func(t *testing.T, input any, target any) error {
		t.Helper()
		f := reflect.TypeOf(input)
		to := reflect.TypeOf(target)
		_, err := hook.(func(reflect.Type, reflect.Type, any) (any, error))(f, to, input)
		return err
	}

	t.Run("non-string source returns data unchanged", func(t *testing.T) {
		result := decode(t, 42, map[string]string{})
		assert.Equal(t, 42, result)
	})

	t.Run("non-map target returns data unchanged", func(t *testing.T) {
		result := decode(t, "key:value", "")
		assert.Equal(t, "key:value", result)
	})

	t.Run("empty string returns empty map", func(t *testing.T) {
		result := decode(t, "", map[string]string{})
		assert.Equal(t, map[string]string{}, result)
	})

	t.Run("map[string]string single pair", func(t *testing.T) {
		result := decode(t, "key1:value1", map[string]string{})
		assert.Equal(t, map[string]string{"key1": "value1"}, result)
	})

	t.Run("map[string]string multiple pairs", func(t *testing.T) {
		result := decode(t, "key1:value1 key2:value2", map[string]string{})
		assert.Equal(t, map[string]string{"key1": "value1", "key2": "value2"}, result)
	})

	t.Run("map[string]string value with colon creates nesting error", func(t *testing.T) {
		// "url:http://example.com" splits into ["url", "http", "//example.com"]
		// which creates a nested map — incompatible with map[string]string target
		err := decodeErr(t, "url:http://example.com", map[string]string{})
		assert.Error(t, err)
	})

	t.Run("map[string]int", func(t *testing.T) {
		result := decode(t, "a:1 b:2", map[string]int{})
		assert.Equal(t, map[string]int{"a": 1, "b": 2}, result)
	})

	t.Run("map[string]int8", func(t *testing.T) {
		result := decode(t, "a:1", map[string]int8{})
		assert.Equal(t, map[string]int8{"a": 1}, result)
	})

	t.Run("map[string]int16", func(t *testing.T) {
		result := decode(t, "a:1", map[string]int16{})
		assert.Equal(t, map[string]int16{"a": 1}, result)
	})

	t.Run("map[string]int32", func(t *testing.T) {
		result := decode(t, "a:1", map[string]int32{})
		assert.Equal(t, map[string]int32{"a": 1}, result)
	})

	t.Run("map[string]int64", func(t *testing.T) {
		result := decode(t, "a:42", map[string]int64{})
		assert.Equal(t, map[string]int64{"a": 42}, result)
	})

	t.Run("map[string]int invalid value", func(t *testing.T) {
		err := decodeErr(t, "a:notanint", map[string]int{})
		assert.Error(t, err)
	})

	t.Run("map[string]uint", func(t *testing.T) {
		result := decode(t, "a:1", map[string]uint{})
		assert.Equal(t, map[string]uint{"a": 1}, result)
	})

	t.Run("map[string]uint8", func(t *testing.T) {
		result := decode(t, "a:1", map[string]uint8{})
		assert.Equal(t, map[string]uint8{"a": 1}, result)
	})

	t.Run("map[string]uint16", func(t *testing.T) {
		result := decode(t, "a:1", map[string]uint16{})
		assert.Equal(t, map[string]uint16{"a": 1}, result)
	})

	t.Run("map[string]uint32", func(t *testing.T) {
		result := decode(t, "a:1", map[string]uint32{})
		assert.Equal(t, map[string]uint32{"a": 1}, result)
	})

	t.Run("map[string]uint64", func(t *testing.T) {
		result := decode(t, "a:1", map[string]uint64{})
		assert.Equal(t, map[string]uint64{"a": 1}, result)
	})

	t.Run("map[string]uint invalid value", func(t *testing.T) {
		err := decodeErr(t, "a:notauint", map[string]uint{})
		assert.Error(t, err)
	})

	t.Run("map[string]float32", func(t *testing.T) {
		result := decode(t, "a:1.5", map[string]float32{})
		assert.Equal(t, map[string]float32{"a": 1.5}, result)
	})

	t.Run("map[string]float64", func(t *testing.T) {
		result := decode(t, "a:2.5", map[string]float64{})
		assert.Equal(t, map[string]float64{"a": 2.5}, result)
	})

	t.Run("map[string]float invalid value", func(t *testing.T) {
		err := decodeErr(t, "a:notafloat", map[string]float64{})
		assert.Error(t, err)
	})

	t.Run("map[string]bool", func(t *testing.T) {
		result := decode(t, "a:true b:false", map[string]bool{})
		assert.Equal(t, map[string]bool{"a": true, "b": false}, result)
	})

	t.Run("map[string]bool invalid value", func(t *testing.T) {
		err := decodeErr(t, "a:notabool", map[string]bool{})
		assert.Error(t, err)
	})

	t.Run("pair without colon returns error", func(t *testing.T) {
		err := decodeErr(t, "valid:pair nocolon", map[string]string{})
		assert.Error(t, err)
	})

	t.Run("empty value", func(t *testing.T) {
		result := decode(t, "key:", map[string]string{})
		assert.Equal(t, map[string]string{"key": ""}, result)
	})

	t.Run("comma separator", func(t *testing.T) {
		commaHook := StringToMapDecodeHookFunc(",")
		f := reflect.TypeFor[string]()
		to := reflect.TypeFor[map[string]string]()
		result, err := commaHook.(func(reflect.Type, reflect.Type, any) (any, error))(f, to, "a:1,b:2")
		require.NoError(t, err)
		assert.Equal(t, map[string]string{"a": "1", "b": "2"}, result)
	})

	t.Run("map[string]string returns empty map for empty string", func(t *testing.T) {
		result := decode(t, "", map[string]int{})
		assert.Equal(t, map[string]int{}, result)
	})

	t.Run("duplicate keys last wins", func(t *testing.T) {
		result := decode(t, "key:first key:second", map[string]string{})
		assert.Equal(t, map[string]string{"key": "second"}, result)
	})

	// Nested map tests
	t.Run("nested map[string]map[string]string", func(t *testing.T) {
		result := decode(t, "outer:inner:value", map[string]map[string]string{})
		assert.Equal(t, map[string]map[string]string{"outer": {"inner": "value"}}, result)
	})

	t.Run("nested map multiple outer keys", func(t *testing.T) {
		result := decode(t, "a:x:1 b:y:2", map[string]map[string]string{})
		assert.Equal(t, map[string]map[string]string{"a": {"x": "1"}, "b": {"y": "2"}}, result)
	})

	t.Run("nested map multiple inner keys same outer", func(t *testing.T) {
		result := decode(t, "a:x:1 a:y:2", map[string]map[string]string{})
		assert.Equal(t, map[string]map[string]string{"a": {"x": "1", "y": "2"}}, result)
	})

	t.Run("nested map[string]map[string]int", func(t *testing.T) {
		result := decode(t, "group:count:42", map[string]map[string]int{})
		assert.Equal(t, map[string]map[string]int{"group": {"count": 42}}, result)
	})

	t.Run("nested map[string]map[string]bool", func(t *testing.T) {
		result := decode(t, "feat:enabled:true feat:visible:false", map[string]map[string]bool{})
		assert.Equal(t, map[string]map[string]bool{"feat": {"enabled": true, "visible": false}}, result)
	})

	t.Run("deeply nested map depth 3", func(t *testing.T) {
		result := decode(t, "a:b:c:leaf", map[string]map[string]map[string]string{})
		assert.Equal(t, map[string]map[string]map[string]string{"a": {"b": {"c": "leaf"}}}, result)
	})

	t.Run("nested map value with colon creates deeper nesting error", func(t *testing.T) {
		// "url:host:http://example.com" splits into ["url", "host", "http", "//example.com"]
		// which creates 3 levels of nesting — incompatible with map[string]map[string]string target
		err := decodeErr(t, "url:host:http://example.com", map[string]map[string]string{})
		assert.Error(t, err)
	})

	t.Run("nested map empty string returns empty map", func(t *testing.T) {
		result := decode(t, "", map[string]map[string]string{})
		assert.Equal(t, map[string]map[string]string{}, result)
	})

	t.Run("nested map pair without enough colons returns error", func(t *testing.T) {
		err := decodeErr(t, "valid:inner:val nocolon", map[string]map[string]string{})
		assert.Error(t, err)
	})

	// map[string]any tests
	t.Run("map[string]any flat key:value", func(t *testing.T) {
		result := decode(t, "a:hello b:world", map[string]any{})
		assert.Equal(t, map[string]any{"a": "hello", "b": "world"}, result)
	})

	t.Run("map[string]any nested key:nestedKey:value", func(t *testing.T) {
		result := decode(t, "a:b:value", map[string]any{})
		assert.Equal(t, map[string]any{"a": map[string]any{"b": "value"}}, result)
	})

	t.Run("map[string]any deeply nested", func(t *testing.T) {
		result := decode(t, "a:b:c:leaf", map[string]any{})
		assert.Equal(t, map[string]any{"a": map[string]any{"b": map[string]any{"c": "leaf"}}}, result)
	})

	t.Run("map[string]any mixed flat and nested", func(t *testing.T) {
		result := decode(t, "flat:val nested:inner:deep", map[string]any{})
		assert.Equal(t, map[string]any{"flat": "val", "nested": map[string]any{"inner": "deep"}}, result)
	})

	t.Run("map[string]any multiple inner keys same outer", func(t *testing.T) {
		result := decode(t, "a:x:1 a:y:2", map[string]any{})
		assert.Equal(t, map[string]any{"a": map[string]any{"x": "1", "y": "2"}}, result)
	})

	t.Run("map[string]any empty string", func(t *testing.T) {
		result := decode(t, "", map[string]any{})
		assert.Equal(t, map[string]any{}, result)
	})

	t.Run("map[string]any pair without colon returns error", func(t *testing.T) {
		err := decodeErr(t, "nocolon", map[string]any{})
		assert.Error(t, err)
	})
}
