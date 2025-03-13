package http

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/kkrt-labs/go-utils/app/svc"
	"github.com/kkrt-labs/go-utils/log"
	"github.com/kkrt-labs/go-utils/tag"
	"go.uber.org/zap"
)

// Entrypoint listens on a local network address and serves incoming HTTP requests.
type Entrypoint struct {
	addr string

	lCfg   *net.ListenConfig
	server *http.Server

	tlsCfg *tls.Config

	mux sync.RWMutex
	l   net.Listener

	done   chan struct{}
	srvErr error

	tagged *svc.Tagged
}

type EntrypointOption func(*Entrypoint) error

// WithServer sets the http.Server to use for the entrypoint.
func WithServer(srv *http.Server) EntrypointOption {
	return func(ep *Entrypoint) error {
		ep.server = srv
		return nil
	}
}

// WithListenConfig sets the net.ListenConfig to use for the entrypoint.
func WithListenConfig(lCfg *net.ListenConfig) EntrypointOption {
	return func(ep *Entrypoint) error {
		ep.lCfg = lCfg
		return nil
	}
}

// WithTags sets the tags to use for the entrypoint.
func WithTags(tags ...*tag.Tag) EntrypointOption {
	return func(ep *Entrypoint) error {
		ep.WithTags(tags...)
		return nil
	}
}

// WithTLSConfig sets the tls.Config to use for the entrypoint.
func WithTLSConfig(tlsCfg *tls.Config) EntrypointOption {
	return func(ep *Entrypoint) error {
		ep.tlsCfg = tlsCfg
		return nil
	}
}

// NewEntrypoint creates a new Entrypoint.
func NewEntrypoint(addr string, opts ...EntrypointOption) (*Entrypoint, error) {
	ep := &Entrypoint{
		addr:   addr,
		lCfg:   &net.ListenConfig{},
		server: &http.Server{},
		tagged: svc.NewTagged(),
	}

	for _, opt := range opts {
		if err := opt(ep); err != nil {
			return nil, err
		}
	}

	baseCtxFunc := ep.server.BaseContext
	if baseCtxFunc == nil {
		ep.server.BaseContext = func(_ net.Listener) context.Context {
			return ep.context(context.Background())
		}
	} else {
		ep.server.BaseContext = func(l net.Listener) context.Context {
			return ep.context(baseCtxFunc(l))
		}
	}

	return ep, nil
}

func (ep *Entrypoint) context(ctx context.Context) context.Context {
	return ep.tagged.Context(ctx)
}

// Addr returns the address the entrypoint is exposed to after Start() is called.
func (ep *Entrypoint) Addr() string {
	ep.mux.RLock()
	defer ep.mux.RUnlock()

	if ep.l == nil {
		return ""
	}
	return ep.l.Addr().String()
}

// SetHandler sets the handler for the entrypoint.
func (ep *Entrypoint) SetHandler(handler http.Handler) {
	ep.server.Handler = handler
}

// Start starts the entrypoint.
func (ep *Entrypoint) Start(ctx context.Context) error {
	ctx = ep.context(ctx)

	// Open connection and return possibly error
	l, err := ep.listen(ctx)
	if err != nil {
		return err
	}

	ep.mux.Lock()
	ep.l = l
	ep.mux.Unlock()

	return ep.serve(ctx, l)
}

// Stop stops the entrypoint.
func (ep *Entrypoint) Stop(ctx context.Context) error {
	ctx = ep.context(ctx)

	logger := log.LoggerFromContext(ctx)
	logger.Info("Entrypoint gracefully stopping...")

	// Gracefully shutdown server
	err := ep.server.Shutdown(ctx)
	if err != nil {
		logger.Error("Error while stopping entrypoint", zap.Error(err))
		_ = ep.server.Close()
		return err
	}

	// Wait for Serve(...) to be done
	<-ep.done

	// Return possible error from Serve(...)
	if ep.srvErr != nil && ep.srvErr != http.ErrServerClosed {
		return ep.srvErr
	}

	logger.Info("Entrypoint successfully stopped")

	return nil
}

func (ep *Entrypoint) listen(ctx context.Context) (net.Listener, error) {
	logger := ep.logger(ctx)

	logger.Info(
		"Open entrypoint on local network",
		zap.String("network", "tcp"),
		zap.String("address", ep.addr),
	)

	l, err := ep.lCfg.Listen(ctx, "tcp", ep.addr)
	if err != nil {
		ep.srvErr = err
		logger.Error("Failed to open entrypoint on local network", zap.Error(err))
		return nil, err
	}

	if ep.tlsCfg != nil {
		logger.Info("Entrypoint upgrades to TLS")
		l = tls.NewListener(l, ep.tlsCfg)
	}

	return l, nil
}

// serve serves incoming HTTP requests.
func (ep *Entrypoint) serve(ctx context.Context, l net.Listener) error {
	logger := ep.logger(ctx)

	logger.Info("Entrypoint is accepting and serving incoming HTTP requests...")
	ep.done = make(chan struct{})

	go func() {
		ep.srvErr = ep.server.Serve(l)
		if ep.srvErr != nil && ep.srvErr != http.ErrServerClosed {
			logger.Error("Entrypoint failed while serving incoming HTTP requests", zap.Error(ep.srvErr))
		}
		close(ep.done)
	}()

	return nil
}

func (ep *Entrypoint) logger(ctx context.Context) *zap.Logger {
	return log.LoggerFromContext(ep.context(ctx))
}

// Ready returns the error from Serve(...) if it's not nil.
func (ep *Entrypoint) Ready(_ context.Context) error {
	return ep.srvErr
}

// WithTags sets the tags for the entrypoint.
func (ep *Entrypoint) WithTags(tags ...*tag.Tag) {
	ep.tagged.WithTags(tags...)
}

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
