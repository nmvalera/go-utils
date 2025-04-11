package config

import (
	"testing"

	"github.com/kkrt-labs/go-utils/common"
	"github.com/spf13/cast"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type NestedConfig struct {
	Flag7 string `key:"flag7" env:"FLAG7" desc:"Flag 7"`
}

type TestConfig struct {
	Flag1 string `key:"flag1" env:"FLAG1" desc:"Flag 1" short:"a"`
	Flag2 int    `key:"flag2" env:"FLAG2" desc:"Flag 2" short:"b"`
	Flag3 bool   `key:"flag3" env:"FLAG3" desc:"Flag 3" short:"c"`

	Flag4 *string `key:"flag4,omitempty" env:"FLAG4" desc:"Flag 4" short:"d"`
	Flag5 *int    `key:"flag5,omitempty" env:"FLAG5" desc:"Flag 5" short:"e"`
	Flag6 *bool   `key:"flag6,omitempty" env:"FLAG6" desc:"Flag 6" short:"f"`

	Nested NestedConfig `key:"nested"`

	Flag8  []string  `key:"flag8" env:"FLAG8" desc:"Flag 8" short:"h"`
	Flag9  [2]string `key:"flag9" env:"FLAG9" desc:"Flag 9" short:"i"`
	Flag10 []*string `key:"flag10" env:"FLAG10" desc:"Flag 10" short:"j"`

	NestedWithSkip NestedConfig `key:"nested-with-skip" env:"-" flag:"-"`
}

func TestAddFlags(t *testing.T) {
	defaultCfg := &TestConfig{
		Flag1: "flag1-default",
		Flag2: 10,
		Flag3: true,

		Flag4: common.Ptr("flag4-default"),
		Flag5: common.Ptr(100),

		Nested: NestedConfig{
			Flag7: "flag7-default",
		},

		Flag8:  []string{"flag8-default-1", "flag8-default-2"},
		Flag9:  [2]string{"flag9-default-1", "flag9-default-2"},
		Flag10: []*string{common.Ptr("flag10-default-1"), nil, common.Ptr("flag10-default-2")},
	}
	v := NewViper()
	err := AddFlags(defaultCfg, v, pflag.NewFlagSet("test", pflag.ContinueOnError), nil)
	require.NoError(t, err)

	expected := map[string]interface{}{
		"flag1":  "flag1-default",
		"flag2":  10,
		"flag3":  true,
		"flag4":  common.Ptr("flag4-default"),
		"flag5":  common.Ptr(100),
		"flag6":  (*bool)(nil),
		"flag8":  []string{"flag8-default-1", "flag8-default-2"},
		"flag9":  []string{"flag9-default-1", "flag9-default-2"},
		"flag10": []string{"flag10-default-1", "", "flag10-default-2"},
		"nested": map[string]interface{}{
			"flag7": "flag7-default",
		},
		"nested-with-skip": map[string]interface{}{
			"flag7": "",
		},
	}

	assert.Equal(t, expected, v.AllSettings())
}

func TestEnv(t *testing.T) {
	cfg := &TestConfig{
		Flag1: "test1",
		Flag2: 10,
		Flag3: true,
		Nested: NestedConfig{
			Flag7: "test7",
		},
		NestedWithSkip: NestedConfig{
			Flag7: "test7",
		},
	}

	m, err := Env(cfg, nil)
	require.NoError(t, err)

	expected := map[string]string{
		"FLAG1":        "test1",
		"FLAG2":        "10",
		"FLAG3":        "true",
		"NESTED_FLAG7": "test7",
		"FLAG7":        "test7",
	}
	assert.Equal(t, expected, m)
}

func TestAddFlagsAndLoadEnv(t *testing.T) {
	cfg := &TestConfig{
		Flag1: "flag1-default",
		Flag2: 10,
		Flag3: true,

		Flag4: common.Ptr("flag4-default"),
		Flag5: common.Ptr(100),

		Nested: NestedConfig{
			Flag7: "flag7-default",
		},

		Flag8:  []string{"flag8-default-1", "flag8-default-2"},
		Flag9:  [2]string{"flag9-default-1", "flag9-default-2"},
		Flag10: []*string{common.Ptr("flag10-default-1"), nil, common.Ptr("flag10-default-2")},
	}
	v := NewViper()
	err := AddFlags(cfg, v, pflag.NewFlagSet("test", pflag.ContinueOnError), nil)
	require.NoError(t, err)

	env, err := Env(cfg, nil)
	t.Logf("env: %+v\n", env)
	require.NoError(t, err)
	cast.ToStringSlice(env["FLAG8"])
	for k, v := range env {
		t.Setenv(k, v)
	}

	loadedCfg := new(TestConfig)
	err = Unmarshal(loadedCfg, v)
	require.NoError(t, err)

	assert.Equal(t, cfg, loadedCfg)
}
