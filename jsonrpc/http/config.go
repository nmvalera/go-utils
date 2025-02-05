package jsonrpchttp

import (
	comhttp "github.com/kkrt-labs/go-utils/net/http"
)

type Config struct {
	HTTP *comhttp.ClientConfig
}

func (cfg *Config) SetDefault() *Config {
	if cfg.HTTP == nil {
		cfg.HTTP = new(comhttp.ClientConfig)
	}

	cfg.HTTP.SetDefault()

	return cfg
}
