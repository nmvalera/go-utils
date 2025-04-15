package store

import (
	"context"
	"io"
	"time"

	"github.com/kkrt-labs/go-utils/app/svc"
	"github.com/kkrt-labs/go-utils/log"
	"github.com/kkrt-labs/go-utils/tag"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type taggable struct {
	store Store

	svc.Tagged
}

func WithTags(store Store) Store {
	return &taggable{store: store}
}

func (s *taggable) Store(ctx context.Context, key string, reader io.Reader, headers *Headers) error {
	return s.store.Store(s.context(ctx, key, headers), key, reader, headers)
}

func (s *taggable) Load(ctx context.Context, key string) (io.ReadCloser, *Headers, error) {
	return s.store.Load(s.context(ctx, key, nil), key)
}

func (s *taggable) Delete(ctx context.Context, key string) error {
	return s.store.Delete(s.context(ctx, key, nil), key)
}

func (s *taggable) Copy(ctx context.Context, srcKey, dstKey string) error {
	tags := []*tag.Tag{
		tag.Key("store.src_key").String(srcKey),
		tag.Key("store.dst_key").String(dstKey),
	}
	return s.store.Copy(s.Context(ctx, tags...), srcKey, dstKey)
}

func (s *taggable) context(ctx context.Context, key string, headers *Headers) context.Context {
	tags := []*tag.Tag{
		tag.Key("store.key").String(key),
	}

	if headers != nil {
		ct := headers.ContentType
		ce := headers.ContentEncoding
		tags = append(
			tags,
			tag.Key("store.content-type").String(ct.String()),
			tag.Key("store.content-encoding").String(ce.String()),
		)
	}

	return s.Context(ctx, tags...)
}

type metrics struct {
	store Store

	loadCount   prometheus.Counter
	storeCount  prometheus.Counter
	copyCount   prometheus.Counter
	deleteCount prometheus.Counter

	loadErrCount   prometheus.Counter
	storeErrCount  prometheus.Counter
	copyErrCount   prometheus.Counter
	deleteErrCount prometheus.Counter

	loadDuration  prometheus.Histogram
	storeDuration prometheus.Histogram
	copyDuration  prometheus.Histogram
}

func WithMetrics(store Store) Store {
	return &metrics{store: store}
}

func (m *metrics) SetMetrics(system, subsystem string, _ ...*tag.Tag) {
	m.loadCount = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: system,
		Subsystem: subsystem,
		Name:      "load_count",
		Help:      "The number of objects successfully loaded from the store",
	})
	m.storeCount = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: system,
		Subsystem: subsystem,
		Name:      "store_count",
		Help:      "The number of objects successfully stored in the store",
	})
	m.deleteCount = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: system,
		Subsystem: subsystem,
		Name:      "delete_count",
		Help:      "The number of objects successfully deleted from the store",
	})
	m.copyCount = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: system,
		Subsystem: subsystem,
		Name:      "copy_count",
		Help:      "The number of objects successfully copied from the store",
	})
	m.loadErrCount = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: system,
		Subsystem: subsystem,
		Name:      "load_err_count",
		Help:      "The number of objects that failed to load from the store",
	})

	m.storeErrCount = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: system,
		Subsystem: subsystem,
		Name:      "store_err_count",
		Help:      "The number of objects that failed to store in the store",
	})
	m.deleteErrCount = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: system,
		Subsystem: subsystem,
		Name:      "delete_err_count",
		Help:      "The number of objects that failed to delete from the store",
	})
	m.copyErrCount = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: system,
		Subsystem: subsystem,
		Name:      "copy_err_count",
		Help:      "The number of objects that failed to copy from the store",
	})
	m.loadDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: system,
		Subsystem: subsystem,
		Name:      "load_duration_seconds",
		Help:      "The duration of the load method (in seconds)",
	})
	m.storeDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: system,
		Subsystem: subsystem,
		Name:      "store_duration_seconds",
		Help:      "The duration of the store method (in seconds)",
	})
	m.copyDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: system,
		Subsystem: subsystem,
		Name:      "copy_duration_seconds",
		Help:      "The duration of the copy method (in seconds)",
	})
}

func (m *metrics) Store(ctx context.Context, key string, reader io.Reader, headers *Headers) error {
	start := time.Now()
	err := m.store.Store(ctx, key, reader, headers)
	duration := time.Since(start)
	m.storeDuration.Observe(duration.Seconds())
	if err != nil {
		m.storeErrCount.Inc()
	} else {
		m.storeCount.Inc()
	}
	return err
}

func (m *metrics) Load(ctx context.Context, key string) (io.ReadCloser, *Headers, error) {
	start := time.Now()
	reader, headers, err := m.store.Load(ctx, key)
	duration := time.Since(start)
	m.loadDuration.Observe(duration.Seconds())
	if err != nil {
		m.loadErrCount.Inc()
	} else {
		m.loadCount.Inc()
	}
	return reader, headers, err
}

func (m *metrics) Delete(ctx context.Context, key string) error {
	err := m.store.Delete(ctx, key)
	if err != nil {
		m.deleteErrCount.Inc()
	} else {
		m.deleteCount.Inc()
	}
	return err
}

func (m *metrics) Copy(ctx context.Context, srcKey, dstKey string) error {
	start := time.Now()
	err := m.store.Copy(ctx, srcKey, dstKey)
	duration := time.Since(start)
	m.copyDuration.Observe(duration.Seconds())
	if err != nil {
		m.copyErrCount.Inc()
	} else {
		m.copyCount.Inc()
	}
	return err
}

func (m *metrics) Describe(ch chan<- *prometheus.Desc) {
	m.loadCount.Describe(ch)
	m.storeCount.Describe(ch)
	m.deleteCount.Describe(ch)
	m.copyCount.Describe(ch)
	m.loadErrCount.Describe(ch)
	m.storeErrCount.Describe(ch)
	m.deleteErrCount.Describe(ch)
	m.copyErrCount.Describe(ch)
	m.loadDuration.Describe(ch)
	m.storeDuration.Describe(ch)
	m.copyDuration.Describe(ch)
}

func (m *metrics) Collect(ch chan<- prometheus.Metric) {
	m.loadCount.Collect(ch)
	m.storeCount.Collect(ch)
	m.deleteCount.Collect(ch)
	m.copyCount.Collect(ch)
	m.loadErrCount.Collect(ch)
	m.storeErrCount.Collect(ch)
	m.deleteErrCount.Collect(ch)
	m.copyErrCount.Collect(ch)
	m.loadDuration.Collect(ch)
	m.storeDuration.Collect(ch)
	m.copyDuration.Collect(ch)
}

type loggable struct {
	store Store
}

func WithLog(store Store) Store {
	return &loggable{store: store}
}

func (l *loggable) Store(ctx context.Context, key string, reader io.Reader, headers *Headers) error {
	logger := log.LoggerFromContext(ctx)
	logger.Debug("Store store object")
	err := l.store.Store(ctx, key, reader, headers)
	if err != nil {
		logger.Error("Failed to store store object", zap.Error(err))
	}
	return err
}

func (l *loggable) Load(ctx context.Context, key string) (io.ReadCloser, *Headers, error) {
	logger := log.LoggerFromContext(ctx)
	logger.Debug("Load store object")
	reader, headers, err := l.store.Load(ctx, key)
	if err != nil {
		logger.Error("Failed to load store object", zap.Error(err))
	}
	return reader, headers, err
}

func (l *loggable) Delete(ctx context.Context, key string) error {
	logger := log.LoggerFromContext(ctx)
	logger.Debug("Delete store object")
	err := l.store.Delete(ctx, key)
	if err != nil {
		logger.Error("Failed to delete store object", zap.Error(err))
	}
	return err
}

func (l *loggable) Copy(ctx context.Context, srcKey, dstKey string) error {
	logger := log.LoggerFromContext(ctx)
	logger.Debug("Copy store object")
	err := l.store.Copy(ctx, srcKey, dstKey)
	if err != nil {
		logger.Error("Failed to copy store object", zap.Error(err))
	}
	return err
}
