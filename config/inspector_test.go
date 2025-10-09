package config

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/nmvalera/go-utils/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type CustomInt int

func (i CustomInt) String() string {
	return fmt.Sprintf("custom-int-%d", i)
}

func CustomIntDecodeHook(f, t reflect.Type, data any) (any, error) {
	if t == reflect.TypeOf(CustomInt(0)) && f.Kind() == reflect.String {
		i, err := strconv.Atoi(strings.TrimPrefix(data.(string), "custom-int-"))
		return CustomInt(i), err
	}
	return data, nil
}

type Interface interface {
	ToString() string
}

type CustomInterface struct{}

func (i *CustomInterface) ToString() string {
	return "custom interface"
}

func TestEncoder(t *testing.T) {
	inspector := NewInspector(&InspectorConfig{
		TagNames: []string{"foo"},
		EncodeHook: func(val reflect.Value) (reflect.Value, error) {
			if val.Type() == reflect.TypeOf(CustomInt(0)) {
				return reflect.ValueOf(val.Interface().(CustomInt).String()), nil
			}
			if val.Type().Implements(reflect.TypeOf((*Interface)(nil)).Elem()) {
				if val.IsNil() {
					return reflect.ValueOf("<nil>"), nil
				}
				return reflect.ValueOf(val.Interface().(Interface).ToString()), nil
			}
			return val, nil
		},
		IncludeNil: true,
	})

	type IntsConfig struct {
		Int   int   `foo:"int"`
		Int8  int8  `foo:"int8"`
		Int16 int16 `foo:"int16"`
		Int32 int32 `foo:"int32"`
		Int64 int64 `foo:"int64"`
	}

	type UintsConfig struct {
		Uint   uint   `foo:"uint"`
		Uint8  uint8  `foo:"uint8"`
		Uint16 uint16 `foo:"uint16"`
		Uint32 uint32 `foo:"uint32"`
		Uint64 uint64 `foo:"uint64"`
	}

	type FloatsConfig struct {
		Float32 float32 `foo:"float32"`
		Float64 float64 `foo:"float64"`
	}

	type PtrConfig struct {
		String *string `foo:"string"`
	}

	type TestConfig struct {
		Bool               bool             `foo:"bool"`
		String             string           `foo:"string"`
		NilStringPtr       *string          `foo:"nil_string_ptr"`
		Ints               IntsConfig       `foo:"ints"`
		Uints              UintsConfig      `foo:"uints"`
		Floats             FloatsConfig     `foo:"floats"`
		Ptr                *PtrConfig       `foo:"ptr"`
		Array              [2]string        `foo:"array"`
		Slice              []string         `foo:"slice"`
		Custom             CustomInt        `foo:"custom"`
		CustomIntNil       *CustomInt       `foo:"custom_int_nil"`
		ArrayCustom        [2]CustomInt     `foo:"array_custom"`
		SliceCustom        []*CustomInt     `foo:"slice_custom"`
		Interface          Interface        `foo:"interface"`
		CustomInterface    *CustomInterface `foo:"custom_interface"`
		CustomInterfaceNil *CustomInterface `foo:"custom_interface_nil"`
	}

	cfg := TestConfig{
		Bool:   true,
		String: "test",
		Ints: IntsConfig{
			Int:   1,
			Int8:  2,
			Int16: 3,
			Int32: 4,
			Int64: 5,
		},
		Floats: FloatsConfig{
			Float32: 1.1,
			Float64: 2.2,
		},
		Uints: UintsConfig{
			Uint:   1,
			Uint8:  2,
			Uint16: 3,
			Uint32: 4,
			Uint64: 5,
		},
		Ptr: &PtrConfig{
			String: common.Ptr("test"),
		},
		Array: [2]string{"test1", "test2"},
		Slice: []string{"test1", "test2"},
		ArrayCustom: [2]CustomInt{
			CustomInt(1),
			CustomInt(2),
		},
		SliceCustom: []*CustomInt{
			common.Ptr(CustomInt(1)),
			nil,
			common.Ptr(CustomInt(2)),
		},
		Interface:       &CustomInterface{},
		CustomInterface: &CustomInterface{},
	}

	enc, err := inspector.Inspect(cfg)
	require.NoError(t, err)

	expected := map[string]*FieldInfo{
		"Bool": {
			FieldParts: []string{"Bool"},
			TagParts:   map[string][]string{"foo": {"bool"}},
			Value:      reflect.ValueOf(true),
			Processed:  reflect.ValueOf(true),
			Encoded:    "true",
		},
		"String": {
			FieldParts: []string{"String"},
			TagParts:   map[string][]string{"foo": {"string"}},
			Value:      reflect.ValueOf("test"),
			Processed:  reflect.ValueOf("test"),
			Encoded:    "test",
		},
		"NilStringPtr": {
			FieldParts: []string{"NilStringPtr"},
			TagParts:   map[string][]string{"foo": {"nil_string_ptr"}},
			Value:      reflect.ValueOf((*string)(nil)),
			Processed:  reflect.ValueOf(""),
			Encoded:    "",
			IsNil:      true,
		},
		"Ints.Int": {
			FieldParts: []string{"Ints", "Int"},
			TagParts:   map[string][]string{"foo": {"ints", "int"}},
			Value:      reflect.ValueOf(int(1)),
			Processed:  reflect.ValueOf(int(1)),
			Encoded:    "1",
		},
		"Ints.Int8": {
			FieldParts: []string{"Ints", "Int8"},
			TagParts:   map[string][]string{"foo": {"ints", "int8"}},
			Value:      reflect.ValueOf(int8(2)),
			Processed:  reflect.ValueOf(int8(2)),
			Encoded:    "2",
		},
		"Ints.Int16": {
			FieldParts: []string{"Ints", "Int16"},
			TagParts:   map[string][]string{"foo": {"ints", "int16"}},
			Value:      reflect.ValueOf(int16(3)),
			Processed:  reflect.ValueOf(int16(3)),
			Encoded:    "3",
		},
		"Ints.Int32": {
			FieldParts: []string{"Ints", "Int32"},
			TagParts:   map[string][]string{"foo": {"ints", "int32"}},
			Value:      reflect.ValueOf(int32(4)),
			Processed:  reflect.ValueOf(int32(4)),
			Encoded:    "4",
		},
		"Ints.Int64": {
			FieldParts: []string{"Ints", "Int64"},
			TagParts:   map[string][]string{"foo": {"ints", "int64"}},
			Value:      reflect.ValueOf(int64(5)),
			Processed:  reflect.ValueOf(int64(5)),
			Encoded:    "5",
		},
		"Floats.Float32": {
			FieldParts: []string{"Floats", "Float32"},
			TagParts:   map[string][]string{"foo": {"floats", "float32"}},
			Value:      reflect.ValueOf(float32(1.1)),
			Processed:  reflect.ValueOf(float32(1.1)),
			Encoded:    "1.1",
		},
		"Floats.Float64": {
			FieldParts: []string{"Floats", "Float64"},
			TagParts:   map[string][]string{"foo": {"floats", "float64"}},
			Value:      reflect.ValueOf(float64(2.2)),
			Processed:  reflect.ValueOf(float64(2.2)),
			Encoded:    "2.2",
		},
		"Ptr.String": {
			FieldParts: []string{"Ptr", "String"},
			TagParts:   map[string][]string{"foo": {"ptr", "string"}},
			Value:      reflect.ValueOf(common.Ptr("test")),
			Processed:  reflect.ValueOf("test"),
			Encoded:    "test",
		},
		"Array": {
			FieldParts: []string{"Array"},
			TagParts:   map[string][]string{"foo": {"array"}},
			Value:      reflect.ValueOf([2]string{"test1", "test2"}),
			Processed:  reflect.ValueOf([2]string{"test1", "test2"}),
			Encoded:    "test1 test2",
		},
		"Slice": {
			FieldParts: []string{"Slice"},
			TagParts:   map[string][]string{"foo": {"slice"}},
			Value:      reflect.ValueOf([]string{"test1", "test2"}),
			Processed:  reflect.ValueOf([]string{"test1", "test2"}),
			Encoded:    "test1 test2",
		},
		"Custom": {
			FieldParts: []string{"Custom"},
			TagParts:   map[string][]string{"foo": {"custom"}},
			Value:      reflect.ValueOf(CustomInt(0)),
			Processed:  reflect.ValueOf("custom-int-0"),
			Encoded:    "custom-int-0",
		},
		"CustomIntNil": {
			FieldParts: []string{"CustomIntNil"},
			TagParts:   map[string][]string{"foo": {"custom_int_nil"}},
			Value:      reflect.ValueOf((*CustomInt)(nil)),
			Processed:  reflect.ValueOf("custom-int-0"),
			Encoded:    "custom-int-0",
		},
		"Uints.Uint": {
			FieldParts: []string{"Uints", "Uint"},
			TagParts:   map[string][]string{"foo": {"uints", "uint"}},
			Value:      reflect.ValueOf(uint(1)),
			Processed:  reflect.ValueOf(uint(1)),
			Encoded:    "1",
		},
		"Uints.Uint8": {
			FieldParts: []string{"Uints", "Uint8"},
			TagParts:   map[string][]string{"foo": {"uints", "uint8"}},
			Value:      reflect.ValueOf(uint8(2)),
			Processed:  reflect.ValueOf(uint8(2)),
			Encoded:    "2",
		},
		"Uints.Uint16": {
			FieldParts: []string{"Uints", "Uint16"},
			TagParts:   map[string][]string{"foo": {"uints", "uint16"}},
			Value:      reflect.ValueOf(uint16(3)),
			Processed:  reflect.ValueOf(uint16(3)),
			Encoded:    "3",
		},
		"Uints.Uint32": {
			FieldParts: []string{"Uints", "Uint32"},
			TagParts:   map[string][]string{"foo": {"uints", "uint32"}},
			Value:      reflect.ValueOf(uint32(4)),
			Processed:  reflect.ValueOf(uint32(4)),
			Encoded:    "4",
		},
		"Uints.Uint64": {
			FieldParts: []string{"Uints", "Uint64"},
			TagParts:   map[string][]string{"foo": {"uints", "uint64"}},
			Value:      reflect.ValueOf(uint64(5)),
			Processed:  reflect.ValueOf(uint64(5)),
			Encoded:    "5",
		},
		"ArrayCustom": {
			FieldParts: []string{"ArrayCustom"},
			TagParts:   map[string][]string{"foo": {"array_custom"}},
			Value:      reflect.ValueOf([2]CustomInt{CustomInt(1), CustomInt(2)}),
			Processed:  reflect.ValueOf([2]string{"custom-int-1", "custom-int-2"}),
			Encoded:    "custom-int-1 custom-int-2",
		},
		"SliceCustom": {
			FieldParts: []string{"SliceCustom"},
			TagParts:   map[string][]string{"foo": {"slice_custom"}},
			Value:      reflect.ValueOf([]*CustomInt{common.Ptr(CustomInt(1)), nil, common.Ptr(CustomInt(2))}),
			Processed:  reflect.ValueOf([]string{"custom-int-1", "<nil>", "custom-int-2"}),
			Encoded:    "custom-int-1 <nil> custom-int-2",
		},
		"Interface": {
			FieldParts: []string{"Interface"},
			TagParts:   map[string][]string{"foo": {"interface"}},
			Value:      reflect.ValueOf(&CustomInterface{}),
			Processed:  reflect.ValueOf("custom interface"),
			Encoded:    "custom interface",
		},
		"CustomInterface": {
			FieldParts: []string{"CustomInterface"},
			TagParts:   map[string][]string{"foo": {"custom_interface"}},
			Value:      reflect.ValueOf(&CustomInterface{}),
			Processed:  reflect.ValueOf("custom interface"),
			Encoded:    "custom interface",
		},
		"CustomInterfaceNil": {
			FieldParts: []string{"CustomInterfaceNil"},
			TagParts:   map[string][]string{"foo": {"custom_interface_nil"}},
			Value:      reflect.ValueOf((*CustomInterface)(nil)),
			Processed:  reflect.ValueOf("<nil>"),
			Encoded:    "<nil>",
		},
	}
	require.Len(t, enc, len(expected))
	for k, expectedInfo := range expected {
		actualInfo, ok := enc[k]
		require.Truef(t, ok, "expected key %s to be present", k)
		assert.Equalf(t, expectedInfo.FieldParts, actualInfo.FieldParts, "expected field parts for key %s to be %v, but got %v", k, expectedInfo.FieldParts, actualInfo.FieldParts)
		assert.Equalf(t, expectedInfo.TagParts, actualInfo.TagParts, "expected tag parts for key %s to be %v, but got %v", k, expectedInfo.TagParts, actualInfo.TagParts)
		assert.Equalf(t, expectedInfo.Value.Interface(), actualInfo.Value.Interface(), "expected value for key %s to be %v, but got %v", k, expectedInfo.Value, actualInfo.Value)
		assert.Equalf(t, expectedInfo.Processed.Interface(), actualInfo.Processed.Interface(), "expected processed value for key %s to be %v, but got %v", k, expectedInfo.Processed, actualInfo.Processed)
		assert.Equalf(t, expectedInfo.Encoded, actualInfo.Encoded, "expected encoded value for key %s to be %v, but got %v", k, expectedInfo.Encoded, actualInfo.Encoded)
	}
}
