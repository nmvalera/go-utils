package net

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/kkrt-labs/go-utils/log"
	"github.com/kkrt-labs/go-utils/tag"
	"go.uber.org/zap"
)

type EntrypointConfig struct {
	Network   string
	Address   string
	KeepAlive time.Duration
	TLS       *tls.Config
}

func (cfg *EntrypointConfig) SetDefault() *EntrypointConfig {
	if cfg.Network == "" {
		cfg.Network = "tcp"
	}

	return cfg
}

type Entrypoint struct {
	cfg *EntrypointConfig

	lCfg net.ListenConfig
}

func NewEntrypoint(cfg *EntrypointConfig) *Entrypoint {
	return &Entrypoint{
		cfg: cfg,
		lCfg: net.ListenConfig{
			KeepAlive: cfg.KeepAlive,
		},
	}
}

func (ep *Entrypoint) logger(ctx context.Context) *zap.Logger {
	ctx = tag.WithComponent(ctx, "entrypoint")
	ctx = tag.WithTags(ctx, tag.Key("network").String(ep.cfg.Network), tag.Key("address").String(ep.cfg.Address))

	return log.LoggerFromContext(ctx)
}

func (ep *Entrypoint) Listen(ctx context.Context) (l net.Listener, err error) {
	logger := ep.logger(ctx)

	logger.Info("Announces on local network...")
	l, err = ep.lCfg.Listen(ctx, ep.cfg.Network, ep.cfg.Address)
	if err != nil {
		logger.Error("Failed to announce on local network", zap.Error(err))
		return
	}
	logger.Info("Announced on local network", zap.String("address", l.Addr().String()))

	if ep.cfg.TLS != nil {
		logger.Info("Configure TLS")
		l = tls.NewListener(l, ep.cfg.TLS)
	}

	return l, nil
}
