package config

import (
	"testing"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/nmvalera/go-utils/common"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnv(t *testing.T) {
	type NestedConfig struct {
		Flag7 string `key:"flag7" env:"FLAG7" desc:"Flag 7"`
	}

	type TestConfig struct {
		Flag1 string `key:"flag1" env:"FLAG1" desc:"Flag 1" short:"a"`
		Flag2 int    `key:"flag2" env:"FLAG2" desc:"Flag 2" short:"b"`
		Flag3 bool   `key:"flag3" env:"FLAG3" desc:"Flag 3" short:"c"`

		Nested NestedConfig `key:"nested"`

		Flag8  []string  `key:"flag8" env:"FLAG8" desc:"Flag 8" short:"h"`
		Flag9  [2]string `key:"flag9" env:"FLAG9" desc:"Flag 9" short:"i"`
		Flag10 []*string `key:"flag10" env:"FLAG10" desc:"Flag 10" short:"j"`

		NestedWithSkip NestedConfig `key:"nested-with-skip" env:"-" flag:"-"`

		MapStringString map[string]string `key:"map_string_string" env:"MAP_STRING_STRING" desc:"Map string string" short:"k"`
		MapStringInt    map[string]int    `key:"map_string_int" env:"MAP_STRING_INT" desc:"Map string int" short:"l"`
	}

	cfg := &TestConfig{
		Flag1: "test1",
		Flag2: 10,
		Flag3: true,
		Nested: NestedConfig{
			Flag7: "test7",
		},
		Flag8:  []string{"test8-1", "test8-2"},
		Flag9:  [2]string{"test9-1", "test9-2"},
		Flag10: []*string{common.Ptr("test10-1"), nil, common.Ptr("test10-2")},
		NestedWithSkip: NestedConfig{
			Flag7: "test7",
		},
		MapStringString: map[string]string{"test1": "test2"},
		MapStringInt:    map[string]int{"test1": 1, "test2": 2},
	}

	m, err := Env(cfg, nil)
	require.NoError(t, err)

	expected := map[string]string{
		"FLAG1":             "test1",
		"FLAG2":             "10",
		"FLAG3":             "true",
		"NESTED_FLAG7":      "test7",
		"FLAG8":             "test8-1 test8-2",
		"FLAG9":             "test9-1 test9-2",
		"FLAG10":            "test10-1 <nil> test10-2",
		"FLAG7":             "test7",
		"MAP_STRING_STRING": "test1:test2",
		"MAP_STRING_INT":    "test1:1 test2:2",
	}
	assert.Equal(t, expected, m)
}

type NestedConfig struct {
	String    string  `key:"String,omitempty"`
	Uint64Ptr *uint64 `key:"Uint64Ptr,omitempty"`
}

type TestConfig struct {
	StringNonEmpty         string
	StringEmpty            string
	StringPtrNonEmpty      *string
	StringPtrEmpty         *string
	StringSliceNonEmpty    []string
	StringSliceEmpty       []string
	StringSlicePtrNonEmpty *[]string
	StringSlicePtrEmpty    *[]string
	StringArray            [2]string
	StringPtrSliceNonEmpty []*string
	StringPtrSliceEmpty    []*string
	StringPtrArray         [2]*string

	IntNonEmpty      int
	IntEmpty         int
	IntPtrNonEmpty   *int
	IntPtrEmpty      *int
	IntSliceNonEmpty []int
	IntSliceEmpty    []int
	IntArray         [2]int

	Int8NonEmpty      int8
	Int8Empty         int8
	Int8PtrNonEmpty   *int8
	Int8PtrEmpty      *int8
	Int8SliceNonEmpty []int8
	Int8SliceEmpty    []int8
	Int8Array         [2]int8

	BoolNonEmpty      bool
	BoolEmpty         bool
	BoolPtrNonEmpty   *bool
	BoolPtrEmpty      *bool
	BoolSliceNonEmpty []bool
	BoolSliceEmpty    []bool
	BoolArray         [2]bool

	Uint64NonEmpty      uint64
	Uint64Empty         uint64
	Uint64PtrNonEmpty   *uint64
	Uint64PtrEmpty      *uint64
	Uint64SliceNonEmpty []uint64
	Uint64SliceEmpty    []uint64
	Uint64Array         [2]uint64

	Float32NonEmpty      float32
	Float32Empty         float32
	Float32PtrNonEmpty   *float32
	Float32PtrEmpty      *float32
	Float32SliceNonEmpty []float32
	Float32SliceEmpty    []float32
	Float32Array         [2]float32

	Float64NonEmpty      float64
	Float64Empty         float64
	Float64PtrNonEmpty   *float64
	Float64PtrEmpty      *float64
	Float64SliceNonEmpty []float64
	Float64SliceEmpty    []float64
	Float64Array         [2]float64

	DurationNonEmpty            time.Duration
	DurationEmpty               time.Duration
	DurationPtrNonEmpty         *time.Duration
	DurationPtrEmpty            *time.Duration
	DurationSliceNonEmpty       []time.Duration
	DurationSliceEmpty          []time.Duration
	DurationSlicePtrNonEmpty    *[]time.Duration
	DurationSlicePtrEmpty       *[]time.Duration
	DurationPtrSliceNonEmpty    []*time.Duration
	DurationPtrSliceEmpty       []*time.Duration
	DurationPtrSlicePtrNonEmpty *[]*time.Duration
	DurationPtrSlicePtrEmpty    *[]*time.Duration
	DurationArray               [2]time.Duration

	CustomIntNonEmpty      CustomInt
	CustomIntEmpty         CustomInt
	CustomIntPtrNonEmpty   *CustomInt
	CustomIntPtrEmpty      *CustomInt
	CustomIntSliceNonEmpty []CustomInt
	CustomIntSliceEmpty    []CustomInt
	CustomIntArray         [2]CustomInt

	NestedNonEmpty    NestedConfig
	NestedEmpty       NestedConfig
	NestedPtrNonEmpty *NestedConfig

	MapStringString map[string]string
	MapStringInt    map[string]int
}

func defaultCfg() *TestConfig {
	return &TestConfig{
		StringNonEmpty:    "string-non-empty-default",
		StringPtrNonEmpty: common.Ptr("string-ptr-non-empty-default"),
		StringSliceNonEmpty: []string{
			"string-slice-non-empty-default-1",
			"string-slice-non-empty-default-2",
		},
		StringSlicePtrNonEmpty: &[]string{
			"string-slice-ptr-non-empty-default-1",
			"string-slice-ptr-non-empty-default-2",
		},
		StringArray: [2]string{
			"string-array-default-1",
			"string-array-default-2",
		},
		StringPtrSliceNonEmpty: []*string{
			common.Ptr("string-ptr-slice-non-empty-default-1"),
			nil,
			common.Ptr("string-ptr-slice-non-empty-default-2"),
		},
		StringPtrArray: [2]*string{
			common.Ptr("string-ptr-array-default-1"),
			nil,
		},
		IntNonEmpty:                 10,
		IntEmpty:                    0,
		IntPtrNonEmpty:              common.Ptr(11),
		IntSliceNonEmpty:            []int{12, 13},
		IntArray:                    [2]int{14, 15},
		Int8NonEmpty:                16,
		Int8Empty:                   0,
		Int8PtrNonEmpty:             common.Ptr(int8(17)),
		Int8SliceNonEmpty:           []int8{18, 19},
		Int8Array:                   [2]int8{20, 21},
		BoolNonEmpty:                true,
		BoolEmpty:                   false,
		BoolPtrNonEmpty:             common.Ptr(true),
		BoolSliceNonEmpty:           []bool{true, false},
		BoolArray:                   [2]bool{true, false},
		Uint64NonEmpty:              12,
		Uint64Empty:                 0,
		Uint64PtrNonEmpty:           common.Ptr(uint64(13)),
		Uint64SliceNonEmpty:         []uint64{14, 15},
		Uint64Array:                 [2]uint64{16, 17},
		Float32NonEmpty:             18.0,
		Float32Empty:                0.0,
		Float32PtrNonEmpty:          common.Ptr(float32(19)),
		Float32SliceNonEmpty:        []float32{20.0, 21.0},
		Float32Array:                [2]float32{22.0, 23.0},
		Float64NonEmpty:             24.0,
		Float64Empty:                0.0,
		Float64PtrNonEmpty:          common.Ptr(25.0),
		Float64SliceNonEmpty:        []float64{26.0, 27.0},
		Float64Array:                [2]float64{28.0, 29.0},
		DurationNonEmpty:            time.Second,
		DurationPtrNonEmpty:         common.Ptr(time.Second),
		DurationSliceNonEmpty:       []time.Duration{time.Second, time.Second * 2},
		DurationSlicePtrNonEmpty:    &[]time.Duration{time.Second, time.Second * 2},
		DurationPtrSliceNonEmpty:    []*time.Duration{common.Ptr(time.Second), nil, common.Ptr(time.Second * 2)},
		DurationPtrSlicePtrNonEmpty: &[]*time.Duration{common.Ptr(time.Second), nil, common.Ptr(time.Second * 2)},
		DurationArray:               [2]time.Duration{time.Second, time.Second * 2},
		CustomIntNonEmpty:           CustomInt(30),
		CustomIntPtrNonEmpty:        common.Ptr(CustomInt(31)),
		CustomIntSliceNonEmpty:      []CustomInt{32, 33},
		CustomIntArray:              [2]CustomInt{34, 35},
		NestedNonEmpty: NestedConfig{
			String:    "nested-non-empty-default",
			Uint64Ptr: common.Ptr(uint64(14)),
		},
		NestedEmpty:       NestedConfig{},
		NestedPtrNonEmpty: &NestedConfig{},
		MapStringString:   map[string]string{"test1": "test2"},
		MapStringInt:      map[string]int{"test1": 1, "test2": 2},
	}
}

func nonDefaultCfg() *TestConfig {
	return &TestConfig{
		StringNonEmpty:    "string-non-empty-non-default",
		StringEmpty:       "string-empty-non-default",
		StringPtrNonEmpty: common.Ptr("string-ptr-non-empty-non-default"),
		StringSliceNonEmpty: []string{
			"string-slice-non-empty-non-default-1",
			"string-slice-non-empty-non-default-2",
		},
		StringSlicePtrNonEmpty: &[]string{
			"string-slice-ptr-non-empty-non-default-1",
			"string-slice-ptr-non-empty-non-default-2",
		},
		StringArray: [2]string{
			"string-array-non-default-1",
			"string-array-non-default-2",
		},
		StringPtrSliceNonEmpty: []*string{
			common.Ptr("string-ptr-slice-non-empty-non-default-1"),
			nil,
			common.Ptr("string-ptr-slice-non-empty-non-default-2"),
		},
		StringPtrArray: [2]*string{
			common.Ptr("string-ptr-array-non-default-1"),
			nil,
		},
		IntNonEmpty:                 20,
		IntEmpty:                    0,
		IntPtrNonEmpty:              common.Ptr(21),
		IntSliceNonEmpty:            []int{22, 23},
		IntArray:                    [2]int{24, 25},
		Int8NonEmpty:                26,
		Int8Empty:                   0,
		Int8PtrNonEmpty:             common.Ptr(int8(27)),
		Int8SliceNonEmpty:           []int8{28, 29},
		Int8Array:                   [2]int8{30, 31},
		BoolNonEmpty:                false,
		BoolEmpty:                   true,
		BoolPtrNonEmpty:             common.Ptr(false),
		BoolSliceNonEmpty:           []bool{false, true},
		BoolArray:                   [2]bool{false, true},
		Uint64NonEmpty:              32,
		Uint64Empty:                 0,
		Uint64PtrNonEmpty:           common.Ptr(uint64(33)),
		Uint64SliceNonEmpty:         []uint64{34, 35},
		Uint64Array:                 [2]uint64{36, 37},
		Float32NonEmpty:             38.0,
		Float32Empty:                0.0,
		Float32PtrNonEmpty:          common.Ptr(float32(39)),
		Float32SliceNonEmpty:        []float32{40.0, 41.0},
		Float32Array:                [2]float32{42.0, 43.0},
		Float64NonEmpty:             44.0,
		Float64Empty:                0.0,
		Float64PtrNonEmpty:          common.Ptr(45.0),
		Float64SliceNonEmpty:        []float64{46.0, 47.0},
		Float64Array:                [2]float64{48.0, 49.0},
		DurationNonEmpty:            time.Second * 2,
		DurationPtrNonEmpty:         common.Ptr(time.Second * 3),
		DurationSliceNonEmpty:       []time.Duration{time.Second * 4, time.Second * 5},
		DurationSlicePtrNonEmpty:    &[]time.Duration{time.Second * 6, time.Second * 7},
		DurationPtrSliceNonEmpty:    []*time.Duration{common.Ptr(time.Second * 8), nil, common.Ptr(time.Second * 9)},
		DurationPtrSlicePtrNonEmpty: &[]*time.Duration{common.Ptr(time.Second * 10), nil, common.Ptr(time.Second * 11)},
		DurationArray:               [2]time.Duration{time.Second * 12, time.Second * 13},
		CustomIntNonEmpty:           CustomInt(50),
		CustomIntPtrNonEmpty:        common.Ptr(CustomInt(51)),
		CustomIntSliceNonEmpty:      []CustomInt{52, 53},
		CustomIntArray:              [2]CustomInt{54, 55},
		NestedNonEmpty: NestedConfig{
			String:    "nested-non-empty-non-default",
			Uint64Ptr: common.Ptr(uint64(56)),
		},
		NestedEmpty:       NestedConfig{},
		NestedPtrNonEmpty: &NestedConfig{},
		MapStringString:   map[string]string{"test1-non-default": "test2-non-default"},
		MapStringInt:      map[string]int{"test1": 3, "test2-non-default": 2},
	}
}

func TestAddFlagsAndLoad(t *testing.T) {
	t.Run("FromDefault", func(t *testing.T) {
		v := NewViper()
		err := AddFlags(defaultCfg(), v, pflag.NewFlagSet("test", pflag.ContinueOnError), nil)
		require.NoError(t, err)

		loadedCfg := new(TestConfig)
		err = Unmarshal(loadedCfg, v)
		require.NoError(t, err)
		assert.Equal(t, defaultCfg(), loadedCfg)
	})

	t.Run("FromEnv", func(t *testing.T) {
		v := NewViper(
			viper.WithDecodeHook(mapstructure.ComposeDecodeHookFunc(
				PtrToValueDecodeHookFunc(),
				StringToWeakSliceDecodeHookFunc(envSliceSep),
				StringToMapDecodeHookFunc(envSliceSep),
				NilStrToNilDecodeHookFunc(nilStr),
				mapstructure.StringToTimeDurationHookFunc(),
				CustomIntDecodeHook,
			)),
		)
		err := AddFlags(defaultCfg(), v, pflag.NewFlagSet("test", pflag.ContinueOnError), nil)
		require.NoError(t, err)

		env, err := Env(nonDefaultCfg(), nil)
		require.NoError(t, err)
		for k, v := range env {
			t.Setenv(k, v)
		}

		loadedCfg := new(TestConfig)
		err = Unmarshal(loadedCfg, v)
		require.NoError(t, err)
		assert.Equal(t, nonDefaultCfg(), loadedCfg)
	})

	t.Run("FromFlags", func(t *testing.T) {
		v := NewViper(
			viper.WithDecodeHook(mapstructure.ComposeDecodeHookFunc(
				PtrToValueDecodeHookFunc(),
				StringToWeakSliceDecodeHookFunc(envSliceSep),
				StringToMapDecodeHookFunc(envSliceSep),
				NilStrToNilDecodeHookFunc(nilStr),
				mapstructure.StringToTimeDurationHookFunc(),
				CustomIntDecodeHook,
			)),
		)

		set := pflag.NewFlagSet("test", pflag.ContinueOnError)
		err := AddFlags(defaultCfg(), v, set, nil)
		require.NoError(t, err)

		require.NoError(t, set.Set("stringnonempty", "string-non-empty-flag-2"))
		require.NoError(t, set.Set("stringptrnonempty", "string-ptr-non-empty-flag-2"))
		require.NoError(t, set.Set("stringslicenonempty", "string-slice-non-empty-flag-1,string-slice-non-empty-flag-2"))
		require.NoError(t, set.Set("stringarray", "string-array-flag-1,string-array-flag-2"))
		require.NoError(t, set.Set("stringsliceptrnonempty", "string-slice-ptr-non-empty-flag-1,string-slice-ptr-non-empty-flag-2"))
		require.NoError(t, set.Set("stringptrslicenonempty", "string-ptr-slice-non-empty-flag-1,string-ptr-slice-non-empty-flag-2,<nil>"))
		require.NoError(t, set.Set("stringptrarray", "string-ptr-array-flag-1,<nil>"))
		require.NoError(t, set.Set("intnonempty", "20"))
		require.NoError(t, set.Set("intempty", "0"))
		require.NoError(t, set.Set("intptrnonempty", "21"))
		require.NoError(t, set.Set("intslicenonempty", "22,23"))
		require.NoError(t, set.Set("intarray", "26,27"))
		require.NoError(t, set.Set("int8nonempty", "28"))
		require.NoError(t, set.Set("int8ptrnonempty", "29"))
		require.NoError(t, set.Set("int8slicenonempty", "30,31"))
		require.NoError(t, set.Set("int8array", "34,35"))
		require.NoError(t, set.Set("boolnonempty", "false"))
		require.NoError(t, set.Set("boolptrnonempty", "false"))
		require.NoError(t, set.Set("boolslicenonempty", "false,true"))
		require.NoError(t, set.Set("boolarray", "true,false"))
		require.NoError(t, set.Set("uint64nonempty", "32"))
		require.NoError(t, set.Set("uint64ptrnonempty", "33"))
		require.NoError(t, set.Set("uint64slicenonempty", "34,35"))
		require.NoError(t, set.Set("uint64array", "38,39"))
		require.NoError(t, set.Set("float32nonempty", "38.0"))
		require.NoError(t, set.Set("float32ptrnonempty", "39.0"))
		require.NoError(t, set.Set("float32slicenonempty", "40.0,41.0"))
		require.NoError(t, set.Set("float32array", "44.0,45.0"))
		require.NoError(t, set.Set("float64nonempty", "46.0"))
		require.NoError(t, set.Set("float64ptrnonempty", "47.0"))
		require.NoError(t, set.Set("float64slicenonempty", "48.0,49.0"))
		require.NoError(t, set.Set("float64array", "52.0,53.0"))
		require.NoError(t, set.Set("durationnonempty", "12s"))
		require.NoError(t, set.Set("durationptrnonempty", "13s"))
		require.NoError(t, set.Set("durationslicenonempty", "14s,15s"))
		require.NoError(t, set.Set("durationsliceptrnonempty", "10s,20s"))
		require.NoError(t, set.Set("durationptrslicenonempty", "18s,<nil>,19s"))
		require.NoError(t, set.Set("durationptrsliceptrnonempty", "20s,<nil>,21s"))
		require.NoError(t, set.Set("durationarray", "22s,23s"))
		require.NoError(t, set.Set("customintnonempty", "56"))
		require.NoError(t, set.Set("customintptrnonempty", "custom-int-57"))
		require.NoError(t, set.Set("customintslicenonempty", "custom-int-58,custom-int-59"))
		require.NoError(t, set.Set("customintarray", "custom-int-62,custom-int-63"))
		require.NoError(t, set.Set("nestednonempty-string", "nested-non-empty-string-flag-2"))
		require.NoError(t, set.Set("nestednonempty-uint64ptr", "56"))
		require.NoError(t, set.Set("nestedptrnonempty-string", "nested-ptr-non-empty-string-flag-2"))
		require.NoError(t, set.Set("mapstringstring", "test1-non-default:test2-non-default"))
		require.NoError(t, set.Set("mapstringint", "test1:3,test2-non-default:2"))

		expectedCfg := &TestConfig{
			StringNonEmpty:    "string-non-empty-flag-2",
			StringPtrNonEmpty: common.Ptr("string-ptr-non-empty-flag-2"),
			StringSliceNonEmpty: []string{
				"string-slice-non-empty-flag-1",
				"string-slice-non-empty-flag-2",
			},
			StringSlicePtrNonEmpty: &[]string{
				"string-slice-ptr-non-empty-flag-1",
				"string-slice-ptr-non-empty-flag-2",
			},
			StringArray: [2]string{
				"string-array-flag-1",
				"string-array-flag-2",
			},
			StringPtrSliceNonEmpty: []*string{
				common.Ptr("string-ptr-slice-non-empty-flag-1"),
				common.Ptr("string-ptr-slice-non-empty-flag-2"),
				nil,
			},
			StringPtrArray: [2]*string{
				common.Ptr("string-ptr-array-flag-1"),
				nil,
			},
			IntNonEmpty:                 20,
			IntPtrNonEmpty:              common.Ptr(21),
			IntSliceNonEmpty:            []int{22, 23},
			IntArray:                    [2]int{26, 27},
			Int8NonEmpty:                28,
			Int8PtrNonEmpty:             common.Ptr(int8(29)),
			Int8SliceNonEmpty:           []int8{30, 31},
			Int8Array:                   [2]int8{34, 35},
			BoolNonEmpty:                false,
			BoolPtrNonEmpty:             common.Ptr(false),
			BoolSliceNonEmpty:           []bool{false, true},
			BoolArray:                   [2]bool{true, false},
			Uint64NonEmpty:              32,
			Uint64PtrNonEmpty:           common.Ptr(uint64(33)),
			Uint64SliceNonEmpty:         []uint64{34, 35},
			Uint64Array:                 [2]uint64{38, 39},
			Float32NonEmpty:             38.0,
			Float32PtrNonEmpty:          common.Ptr(float32(39)),
			Float32SliceNonEmpty:        []float32{40.0, 41.0},
			Float32Array:                [2]float32{44.0, 45.0},
			Float64NonEmpty:             46.0,
			Float64PtrNonEmpty:          common.Ptr(47.0),
			Float64SliceNonEmpty:        []float64{48.0, 49.0},
			Float64Array:                [2]float64{52.0, 53.0},
			DurationNonEmpty:            time.Second * 12,
			DurationPtrNonEmpty:         common.Ptr(time.Second * 13),
			DurationSliceNonEmpty:       []time.Duration{time.Second * 14, time.Second * 15},
			DurationSlicePtrNonEmpty:    &[]time.Duration{time.Second * 10, time.Second * 20},
			DurationPtrSliceNonEmpty:    []*time.Duration{common.Ptr(time.Second * 18), nil, common.Ptr(time.Second * 19)},
			DurationPtrSlicePtrNonEmpty: &[]*time.Duration{common.Ptr(time.Second * 20), nil, common.Ptr(time.Second * 21)},
			DurationArray:               [2]time.Duration{time.Second * 22, time.Second * 23},
			CustomIntNonEmpty:           CustomInt(56),
			CustomIntPtrNonEmpty:        common.Ptr(CustomInt(57)),
			CustomIntSliceNonEmpty:      []CustomInt{58, 59},
			CustomIntArray:              [2]CustomInt{62, 63},
			NestedNonEmpty: NestedConfig{
				String:    "nested-non-empty-string-flag-2",
				Uint64Ptr: common.Ptr(uint64(56)),
			},
			NestedPtrNonEmpty: &NestedConfig{
				String: "nested-ptr-non-empty-string-flag-2",
			},
			MapStringString: map[string]string{"test1-non-default": "test2-non-default"},
			MapStringInt:    map[string]int{"test1": 3, "test2-non-default": 2},
		}

		loadedCfg := new(TestConfig)
		err = Unmarshal(loadedCfg, v)
		require.NoError(t, err)
		assert.Equal(t, expectedCfg, loadedCfg)
	})
}
