package jsonrpc

import (
	"context"
	"time"

	"github.com/kkrt-labs/go-utils/app/svc"
	"github.com/kkrt-labs/go-utils/log"
	"github.com/kkrt-labs/go-utils/tag"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type tagged struct {
	client Client
	tagged *svc.Tagged
}

func (t *tagged) WithTags(tags ...*tag.Tag) {
	t.tagged.WithTags(tags...)
}

func (t *tagged) Call(ctx context.Context, req *Request, res interface{}) error {
	tags := []*tag.Tag{
		tag.Key("req.method").String(req.Method),
		tag.Key("req.version").String(req.Version),
		tag.Key("req.params").Object(req.Params),
		tag.Key("req.id").Object(req.ID),
	}

	ctx = t.tagged.Context(ctx, tags...)

	return t.client.Call(ctx, req, res)
}

// WithTags is a decorator that attaches following JSON-RPC specific tags to the Call context .
// - req.method: JSON-RPC method
// - req.version: JSON-RPC version
// - req.params: JSON-RPC params
// - req.id: JSON-RPC id
//
// It also makes the client Taggable
func WithTags(client Client) Client {
	return &tagged{
		client: client,
		tagged: svc.NewTagged(),
	}
}

func WithLog(namespaces ...string) ClientDecorator {
	return func(c Client) Client {
		return ClientFunc(func(ctx context.Context, req *Request, res interface{}) error {
			logger := log.LoggerWithFieldsFromNamespaceContext(ctx, namespaces...)

			logger.Debug("Call JSON-RPC")
			err := c.Call(ctx, req, res)
			if err != nil {
				logger.Error("JSON-RPC call failed", zap.Error(err))
			}

			return err
		})
	}
}

type metricable struct {
	client Client

	duration      *prometheus.HistogramVec
	counterTotal  *prometheus.CounterVec
	counterErrors *prometheus.CounterVec
}

func WithMetrics(client Client) Client {
	return &metricable{
		client: client,
	}
}

func (m *metricable) SetMetrics(system, subsystem string, tags ...*tag.Tag) {
	m.duration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: system,
		Subsystem: subsystem,
		Name:      "request_duration_seconds",
		Help:      "The duration of requests in seconds (per method)",
	}, []string{"method"})

	m.counterTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: system,
		Subsystem: subsystem,
		Name:      "requests_total",
		Help:      "The total number of requests (per method)",
	}, []string{"method"})

	m.counterErrors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: system,
		Subsystem: subsystem,
		Name:      "requests_errors",
		Help:      "The number of requests that failed (per method)",
	}, []string{"method"})
}

func (m *metricable) Call(ctx context.Context, req *Request, res interface{}) error {
	start := time.Now()
	err := m.client.Call(ctx, req, res)
	m.duration.WithLabelValues(req.Method).Observe(time.Since(start).Seconds())
	m.counterTotal.WithLabelValues(req.Method).Inc()
	if err != nil {
		m.counterErrors.WithLabelValues(req.Method).Inc()
	}

	return err
}

func (m *metricable) Describe(ch chan<- *prometheus.Desc) {
	m.duration.Describe(ch)
	m.counterTotal.Describe(ch)
	m.counterErrors.Describe(ch)
}

func (m *metricable) Collect(ch chan<- prometheus.Metric) {
	m.duration.Collect(ch)
	m.counterTotal.Collect(ch)
	m.counterErrors.Collect(ch)
}
