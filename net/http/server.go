package http

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/kkrt-labs/go-utils/log"
	kkrtnet "github.com/kkrt-labs/go-utils/net"
	"github.com/kkrt-labs/go-utils/tag"
	"go.uber.org/zap"
)

type ServerConfig struct {
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

func (cfg *ServerConfig) SetDefault() *ServerConfig {
	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = 30 * time.Second
	}

	if cfg.ReadHeaderTimeout == 0 {
		cfg.ReadHeaderTimeout = 30 * time.Second
	}

	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = 90 * time.Second
	}

	if cfg.IdleTimeout == 0 {
		cfg.IdleTimeout = 90 * time.Second
	}

	return cfg
}

func NewServer(cfg *ServerConfig) *http.Server {
	return &http.Server{
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}
}

type Server struct {
	Entrypoint *kkrtnet.Entrypoint
	Server     *http.Server

	mux sync.Mutex
	l   net.Listener

	done   chan struct{}
	srvErr error
}

func (s *Server) logger(ctx context.Context) *zap.Logger {
	ctx = tag.WithComponent(ctx, "server")
	return log.LoggerFromContext(ctx)
}

func (s *Server) Addr() string {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.l == nil {
		return ""
	}
	return s.l.Addr().String()
}

func (s *Server) Start(ctx context.Context) error {
	// Open connection and return possibly error
	l, err := s.Entrypoint.Listen(ctx)
	if err != nil {
		return err
	}

	s.mux.Lock()
	s.l = l
	s.mux.Unlock()

	s.logger(ctx).Info("Start serving incoming HTTP requests")
	s.done = make(chan struct{})

	// Start serving in a separate go-routine
	go func() {
		s.srvErr = s.Server.Serve(l)
		if s.srvErr != nil && s.srvErr != http.ErrServerClosed {
			s.logger(ctx).Error("Error while serving incoming HTTP requests", zap.Error(s.srvErr))
		} else {
			s.logger(ctx).Info("stopped serving HTTP request")
		}
		close(s.done)
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger(ctx).Info("Stop server...")

	// Gracefully shutdown server
	err := s.Server.Shutdown(ctx)
	if err != nil {
		s.logger(ctx).Error("Error while shutting down server", zap.Error(err))
		_ = s.Server.Close()
		return err
	}

	// Wait for Serve(...) to be done
	<-s.done

	// Return possible error from Serve(...)
	if err == nil && s.srvErr != nil && s.srvErr != http.ErrServerClosed {
		s.logger(ctx).Error("Error while serving incoming HTTP requests", zap.Error(s.srvErr))
		return s.srvErr
	}

	s.logger(ctx).Info("Server successfully stopped")

	return nil
}
