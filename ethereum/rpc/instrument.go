package rpc

import (
	"context"
	"sync"
	"time"

	"github.com/nmvalera/go-utils/app/svc"
	"github.com/nmvalera/go-utils/log"
	"github.com/nmvalera/go-utils/tag"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type MetricsOption func(*metricsClient)

type metricsClient struct {
	Client

	chainID   prometheus.Gauge
	chainHead prometheus.Gauge

	featchInterval time.Duration

	wg   sync.WaitGroup
	stop chan struct{}

	svc.Tagged
}

func WithMetrics(client Client, opts ...MetricsOption) Client {
	mc := &metricsClient{
		Client:         client,
		featchInterval: 500 * time.Millisecond,
	}
	for _, opt := range opts {
		opt(mc)
	}
	return mc
}

func (c *metricsClient) SetMetrics(system, subsystem string, _ ...*tag.Tag) {
	c.chainID = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: system,
		Subsystem: subsystem,
		Name:      "chain_id",
		Help:      "The chain ID of the network",
	})
	c.chainHead = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: system,
		Subsystem: subsystem,
		Name:      "chain_head",
		Help:      "The block number of the chain head",
	})
}

func (c *metricsClient) Start(ctx context.Context) error {
	log.LoggerFromContext(c.wrapCtx(ctx)).Info("Start collecting Ethereum RPC metrics")
	c.stop = make(chan struct{})
	c.wg.Add(1)
	go func() {
		c.fetchLoop()
		c.wg.Done()
	}()

	return nil
}

func (c *metricsClient) Stop(ctx context.Context) error {
	log.LoggerFromContext(c.wrapCtx(ctx)).Info("Stop collecting Ethereum RPC metrics...")
	close(c.stop)
	c.wg.Wait()
	return nil
}

func (c *metricsClient) fetchLoop() {
	ctx := c.wrapCtx(context.Background())
	ticker := time.NewTicker(c.featchInterval)
	defer ticker.Stop()
	for {
		_ = c.fetch(ctx)
		select {
		case <-c.stop:
			return
		case <-ticker.C:
		}
	}
}

// fetch fetches the necessary metrics from the Ethereum RPC.
func (c *metricsClient) fetch(ctx context.Context) error {
	chainID, err := c.ChainID(ctx)
	if err != nil {
		log.LoggerFromContext(ctx).Error("Ethereum RPC metrics failed to fetch chain ID", zap.Error(err))
		return err
	}
	c.chainID.Set(float64(chainID.Int64()))

	head, err := c.HeaderByNumber(ctx, nil)
	if err != nil {
		log.LoggerFromContext(ctx).Error("Ethereum RPC metrics failed to fetch chain head", zap.Error(err))
		return err
	}
	c.chainHead.Set(float64(head.Number.Int64()))

	return nil
}

func (c *metricsClient) wrapCtx(ctx context.Context) context.Context {
	return c.Context(ctx)
}

func (c *metricsClient) Describe(ch chan<- *prometheus.Desc) {
	c.chainID.Describe(ch)
	c.chainHead.Describe(ch)
}

func (c *metricsClient) Collect(ch chan<- prometheus.Metric) {
	c.chainID.Collect(ch)
	c.chainHead.Collect(ch)
}

// WithFetchInterval sets the interval for fetching the metrics from the Ethereum RPC.
func WithFetchInterval(interval time.Duration) MetricsOption {
	return func(mc *metricsClient) {
		mc.featchInterval = interval
	}
}

type checkable struct {
	Client
}

func WithCheck(client Client) Client {
	return &checkable{Client: client}
}

func (c *checkable) Ready(ctx context.Context) error {
	_, err := c.ChainID(ctx)
	return err
}
