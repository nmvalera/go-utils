package main

import (
	"context"
	"fmt"

	"github.com/kkrt-labs/go-utils/app"
	"github.com/kkrt-labs/go-utils/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type MyAppConfig struct {
	App app.Config `mapstructure:"app"`
	Log log.Config `mapstructure:"log"`
}

type MyService struct {
}

func main() {
	v := viper.New()
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	app.AddFlags(v, f)
	log.AddFlags(v, f)

	cfg := &MyAppConfig{}
	err := v.Unmarshal(cfg)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal config: %v", err))
	}

	logCfg, err := log.ParseConfig(&cfg.Log)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse log config: %v", err))
	}

	logger, err := logCfg.Build()
	if err != nil {
		panic(fmt.Sprintf("Failed to build logger: %v", err))
	}

	a, err := app.NewApp(
		&cfg.App,
		app.WithName("my-app"),
		app.WithVersion("1.0.0"),
		app.WithLogger(logger),
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to create app: %v", err))
	}

	app.Provide(a, "my-service", func() (*MyService, error) {
		a.EnableMainEntrypoint()
		a.EnableHealthzEntrypoint()
		return &MyService{}, nil
	})

	err = a.Run(context.Background())
	if err != nil {
		panic(fmt.Sprintf("Failed to run app: %v", err))
	}
}
