package config

import (
	"fmt"

	"github.com/nmvalera/go-utils/common"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type StringFlag struct {
	ViperKey     string
	Name         string
	Shorthand    string
	Env          string
	Description  string
	DefaultValue any
}

func (flag *StringFlag) Add(v *viper.Viper, f *pflag.FlagSet) {
	defaultValue, ok := flag.DefaultValue.(string)
	if !ok {
		ptr, ok := flag.DefaultValue.(*string)
		if !ok {
			panic("default value for flag is not a string")
		}
		defaultValue = common.Val(ptr)
	}

	if flag.Name != "" {
		if flag.Shorthand != "" {
			f.StringP(flag.Name, flag.Shorthand, defaultValue, FlagDesc(flag.Description, flag.Env))
		} else {
			f.String(flag.Name, defaultValue, FlagDesc(flag.Description, flag.Env))
		}
		_ = v.BindPFlag(flag.ViperKey, f.Lookup(flag.Name))
	}
	if flag.Env != "" {
		_ = v.BindEnv(flag.ViperKey, flag.Env)
	}
	v.SetDefault(flag.ViperKey, flag.DefaultValue)
}

type StringArrayFlag struct {
	ViperKey     string
	Name         string
	Shorthand    string
	Env          string
	Description  string
	DefaultValue any
}

func (flag *StringArrayFlag) Add(v *viper.Viper, f *pflag.FlagSet) {
	defaultValue, ok := flag.DefaultValue.([]string)
	if !ok {
		ptr, ok := flag.DefaultValue.(*[]string)
		if !ok {
			panic("default value for flag is not a []string")
		}
		defaultValue = common.Val(ptr)
	}

	if flag.Name != "" {
		if flag.Shorthand != "" {
			f.StringSliceP(flag.Name, flag.Shorthand, defaultValue, FlagDesc(flag.Description, flag.Env))
		} else {
			f.StringSlice(flag.Name, defaultValue, FlagDesc(flag.Description, flag.Env))
		}
		_ = v.BindPFlag(flag.ViperKey, f.Lookup(flag.Name))
	}
	if flag.Env != "" {
		_ = v.BindEnv(flag.ViperKey, flag.Env)
	}

	v.SetDefault(flag.ViperKey, flag.DefaultValue)
}

type BoolFlag struct {
	ViperKey     string
	Name         string
	Shorthand    string
	Env          string
	Description  string
	DefaultValue any
}

func (flag *BoolFlag) Add(v *viper.Viper, f *pflag.FlagSet) {
	defaultValue, ok := flag.DefaultValue.(bool)
	if !ok {
		ptr, ok := flag.DefaultValue.(*bool)
		if !ok {
			panic("default value for flag is not a bool")
		}
		defaultValue = common.Val(ptr)
	}

	if flag.Name != "" {
		if flag.Shorthand != "" {
			f.BoolP(flag.Name, flag.Shorthand, defaultValue, FlagDesc(flag.Description, flag.Env))
		} else {
			f.Bool(flag.Name, defaultValue, FlagDesc(flag.Description, flag.Env))
		}
		_ = v.BindPFlag(flag.ViperKey, f.Lookup(flag.Name))
	}
	if flag.Env != "" {
		_ = v.BindEnv(flag.ViperKey, flag.Env)
	}

	v.SetDefault(flag.ViperKey, flag.DefaultValue)
}

func FlagDesc(desc, envVar string) string {
	if envVar != "" {
		desc = fmt.Sprintf("%v [env: %v]", desc, envVar)
	}

	return desc
}

type IntFlag struct {
	ViperKey     string
	Name         string
	Shorthand    string
	Env          string
	Description  string
	DefaultValue any
}

func (flag *IntFlag) Add(v *viper.Viper, f *pflag.FlagSet) {
	defaultValue, ok := flag.DefaultValue.(int)
	if !ok {
		ptr, ok := flag.DefaultValue.(*int)
		if !ok {
			panic("default value for flag is not an int")
		}
		defaultValue = common.Val(ptr)
	}

	if flag.Name != "" {
		if flag.Shorthand != "" {
			f.IntP(flag.Name, flag.Shorthand, defaultValue, FlagDesc(flag.Description, flag.Env))
		} else {
			f.Int(flag.Name, defaultValue, FlagDesc(flag.Description, flag.Env))
		}
		_ = v.BindPFlag(flag.ViperKey, f.Lookup(flag.Name))
	}
	if flag.Env != "" {
		_ = v.BindEnv(flag.ViperKey, flag.Env)
	}

	v.SetDefault(flag.ViperKey, flag.DefaultValue)
}

type Flag struct {
	ViperKey     string
	Env          string
	Flag         *pflag.Flag
	DefaultValue any
}

func (flag *Flag) Add(v *viper.Viper, f *pflag.FlagSet) {
	f.AddFlag(flag.Flag)
	_ = v.BindPFlag(flag.ViperKey, flag.Flag)
	if flag.Env != "" {
		_ = v.BindEnv(flag.ViperKey, flag.Env)
	}
	v.SetDefault(flag.ViperKey, flag.DefaultValue)
}
