package http

import (
	"context"
	"errors"
	"net"
	"net/http"
	"sync"

	"github.com/nmvalera/go-utils/app/svc"
	"github.com/nmvalera/go-utils/log"
	"go.uber.org/zap"
)

// Entrypoint listens on a local network address and serves incoming HTTP requests.
type Entrypoint struct {
	addr string

	lCfg   *net.ListenConfig
	server *http.Server

	tlsCfg *TLSCertConfig

	mux sync.RWMutex
	l   net.Listener

	done   chan struct{}
	srvErr error

	*svc.RunContext
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

// WithTLSConfig sets the tls.Config to use for the entrypoint.
func WithTLSConfig(tlsCfg *TLSCertConfig) EntrypointOption {
	return func(ep *Entrypoint) error {
		ep.tlsCfg = tlsCfg
		return nil
	}
}

// NewEntrypoint creates a new Entrypoint.
func NewEntrypoint(addr string, opts ...EntrypointOption) (*Entrypoint, error) {
	ep := &Entrypoint{
		addr:       addr,
		lCfg:       &net.ListenConfig{},
		server:     &http.Server{},
		RunContext: &svc.RunContext{},
	}

	for _, opt := range opts {
		if err := opt(ep); err != nil {
			return nil, err
		}
	}

	baseCtxFunc := ep.server.BaseContext
	if baseCtxFunc == nil {
		ep.server.BaseContext = func(_ net.Listener) context.Context {
			return ep.Context()
		}
	}
	return ep, nil
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

func (ep *Entrypoint) Server() *http.Server {
	return ep.server
}

// Start starts the entrypoint.
func (ep *Entrypoint) Start(ctx context.Context) error {
	// Open connection and return possibly error
	l, err := ep.listen(ctx)
	if err != nil {
		return err
	}

	ep.mux.Lock()
	ep.l = l
	ep.mux.Unlock()

	if ep.tlsCfg != nil && ep.tlsCfg.CertFile != nil {
		return ep.serveTLS(ctx, l)
	}

	return ep.serve(ctx, l)
}

// Stop stops the entrypoint.
func (ep *Entrypoint) Stop(stopCtx context.Context) error {
	logger := log.LoggerFromContext(stopCtx)
	logger.Info("Entrypoint gracefully stopping...")

	// Gracefully shutdown server
	err := ep.server.Shutdown(stopCtx)
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

func (ep *Entrypoint) listen(startCtx context.Context) (net.Listener, error) {
	logger := log.LoggerFromContext(startCtx)

	logger.Info(
		"Open entrypoint on local network",
		zap.String("network", "tcp"),
		zap.String("address", ep.addr),
	)

	l, err := ep.lCfg.Listen(startCtx, "tcp", ep.addr)
	if err != nil {
		ep.srvErr = err
		logger.Error("Failed to open entrypoint on local network", zap.Error(err))
		return nil, err
	}

	return l, nil
}

// serve serves incoming HTTP requests.
func (ep *Entrypoint) serve(startCtx context.Context, l net.Listener) error {
	logger := log.LoggerFromContext(startCtx)

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

func (ep *Entrypoint) serveTLS(startCtx context.Context, l net.Listener) error {
	logger := log.LoggerFromContext(startCtx)
	if ep.tlsCfg.CertFile == nil {
		return errors.New("cert file is required")
	}

	if ep.tlsCfg.KeyFile == nil {
		return errors.New("key file is required")
	}

	logger.Info("Entrypoint is accepting and serving incoming HTTPS requests...")
	ep.done = make(chan struct{})

	go func() {
		ep.srvErr = ep.server.ServeTLS(l, *ep.tlsCfg.CertFile, *ep.tlsCfg.KeyFile)
		if ep.srvErr != nil && ep.srvErr != http.ErrServerClosed {
			logger.Error("Entrypoint failed while serving incoming HTTPS requests", zap.Error(ep.srvErr))
		}
		close(ep.done)
	}()

	return nil
}

// Ready returns the error from Serve(...) if it's not nil.
func (ep *Entrypoint) Ready(_ context.Context) error {
	return ep.srvErr
}
