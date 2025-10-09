package main

import (
	"context"
	"fmt"

	"github.com/nmvalera/go-utils/app"
	"github.com/nmvalera/go-utils/config"
	"github.com/spf13/pflag"
)

type MyService struct {
}

func main() {
	v := config.NewViper()
	app.AddFlags(v, pflag.NewFlagSet("test", pflag.ContinueOnError))

	cfg := new(app.Config)
	err := cfg.Unmarshal(v)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal config: %v", err))
	}

	a, err := app.NewApp(
		cfg,
		app.WithName("my-app"),
		app.WithVersion("1.0.0"),
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
