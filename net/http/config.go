package http

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

type EntrypointConfig struct {
	Addr string       `mapstructure:"addr"`
	HTTP ServerConfig `mapstructure:"http"`
	Net  ListenConfig `mapstructure:"net"`
}

func (cfg *EntrypointConfig) Entrypoint() (*Entrypoint, error) {
	srv, err := cfg.HTTP.Server()
	if err != nil {
		return nil, fmt.Errorf("failed to create server: %w", err)
	}

	netCfg, err := cfg.Net.ListenConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create listen config: %w", err)
	}

	return NewEntrypoint(cfg.Addr, WithServer(srv), WithListenConfig(netCfg))
}

type ServerConfig struct {
	ReadTimeout       string `mapstructure:"read-timeout"`
	ReadHeaderTimeout string `mapstructure:"read-header-timeout"`
	WriteTimeout      string `mapstructure:"write-timeout"`
	IdleTimeout       string `mapstructure:"idle-timeout"`
	MaxHeaderBytes    int    `mapstructure:"max-header-bytes"`
}

func (cfg *ServerConfig) Server() (*http.Server, error) {
	srv := &http.Server{
		MaxHeaderBytes: cfg.MaxHeaderBytes,
	}

	if cfg.ReadTimeout != "" {
		timeout, err := time.ParseDuration(cfg.ReadTimeout)
		if err != nil {
			return nil, fmt.Errorf("failed to parse server read timeout: %w", err)
		}
		srv.ReadTimeout = timeout
	}

	if cfg.ReadHeaderTimeout != "" {
		timeout, err := time.ParseDuration(cfg.ReadHeaderTimeout)
		if err != nil {
			return nil, fmt.Errorf("failed to parse server read header timeout: %w", err)
		}
		srv.ReadHeaderTimeout = timeout
	}

	if cfg.WriteTimeout != "" {
		timeout, err := time.ParseDuration(cfg.WriteTimeout)
		if err != nil {
			return nil, fmt.Errorf("failed to parse server write timeout: %w", err)
		}
		srv.WriteTimeout = timeout
	}

	if cfg.IdleTimeout != "" {
		timeout, err := time.ParseDuration(cfg.IdleTimeout)
		if err != nil {
			return nil, fmt.Errorf("failed to parse server idle timeout: %w", err)
		}
		srv.IdleTimeout = timeout
	}

	return srv, nil
}

type ListenConfig struct {
	KeepAlive      string               `mapstructure:"keep-alive"`
	KeepAliveProbe KeepAliveProbeConfig `mapstructure:"keep-alive-probe"`
}

func (cfg *ListenConfig) ListenConfig() (*net.ListenConfig, error) {
	netCfg := &net.ListenConfig{}

	if cfg.KeepAlive != "" {
		idle, err := time.ParseDuration(cfg.KeepAlive)
		if err != nil {
			return nil, err
		}
		netCfg.KeepAlive = idle
	}

	keepAliveProbe, err := cfg.KeepAliveProbe.KeepAliveProbe()
	if err != nil {
		return nil, err
	}
	netCfg.KeepAliveConfig = *keepAliveProbe

	return netCfg, nil
}

type KeepAliveProbeConfig struct {
	Enable   bool   `mapstructure:"enable"`
	Idle     string `mapstructure:"idle"`
	Interval string `mapstructure:"interval"`
	Count    int    `mapstructure:"count"`
}

func (cfg *KeepAliveProbeConfig) KeepAliveProbe() (*net.KeepAliveConfig, error) {
	netCfg := &net.KeepAliveConfig{
		Enable: cfg.Enable,
		Count:  cfg.Count,
	}

	if cfg.Idle != "" {
		idle, err := time.ParseDuration(cfg.Idle)
		if err != nil {
			return nil, err
		}
		netCfg.Idle = idle
	}

	if cfg.Interval != "" {
		interval, err := time.ParseDuration(cfg.Interval)
		if err != nil {
			return nil, err
		}
		netCfg.Interval = interval
	}

	return netCfg, nil
}
