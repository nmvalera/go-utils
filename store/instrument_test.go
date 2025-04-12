package store_test

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/kkrt-labs/go-utils/app/svc"
	kkrtgomock "github.com/kkrt-labs/go-utils/gomock"
	"github.com/kkrt-labs/go-utils/store"
	"github.com/kkrt-labs/go-utils/store/mock"
	"github.com/kkrt-labs/go-utils/tag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestImplementsInterface(t *testing.T) {
	assert.Implements(t, (*store.Store)(nil), store.WithTags(nil))
	assert.Implements(t, (*store.Store)(nil), store.WithMetrics(nil))
	assert.Implements(t, (*store.Store)(nil), store.WithLog(nil))
}

func TestWithTags(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockStore(ctrl)
	taggedStore := store.WithTags(mockStore)

	require.Implements(t, (*store.Store)(nil), taggedStore)
	require.Implements(t, (*svc.Taggable)(nil), taggedStore)

	taggedStore.(svc.Taggable).WithTags(tag.Key("component").String("test-component"))

	validateCtx := func(ctx context.Context) error {
		return tag.ExpectTagsOnContext(
			ctx,
			tag.Key("component").String("test-component"),
			tag.Key("store.key").String("test-key"),
			tag.Key("store.content-type").String("application/protobuf"),
			tag.Key("store.content-encoding").String("gzip"),
		)
	}
	mockStore.EXPECT().Store(kkrtgomock.ContextMatcher(validateCtx), "test-key", gomock.Any(), gomock.Any()).Return(nil)

	err := taggedStore.Store(context.Background(), "test-key", strings.NewReader("test-value"), &store.Headers{
		ContentType:     store.ContentTypeProtobuf,
		ContentEncoding: store.ContentEncodingGzip,
	})

	require.NoError(t, err)
}

func TestWithMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockStore(ctrl)
	metricsStore := store.WithMetrics(mockStore)

	require.Implements(t, (*store.Store)(nil), metricsStore)
	require.Implements(t, (*svc.Metricable)(nil), metricsStore)
	require.Implements(t, (*svc.MetricsCollector)(nil), metricsStore)

	metricsStore.(svc.Metricable).SetMetrics("test-system", "test-subsystem")

	ctx := context.TODO()
	reader := io.NopCloser(strings.NewReader("test-value"))
	headers := new(store.Headers)
	mockStore.EXPECT().Store(ctx, "test-key", reader, headers).Return(nil)
	err := metricsStore.Store(ctx, "test-key", reader, headers)
	require.NoError(t, err)

	mockStore.EXPECT().Load(ctx, "test-key").Return(reader, headers, nil)
	resReader, resHeaders, err := metricsStore.Load(ctx, "test-key")
	require.NoError(t, err)
	require.Equal(t, reader, resReader)
	require.Equal(t, headers, resHeaders)
	ch := make(chan *prometheus.Desc)
	go func() {
		metricsStore.(svc.MetricsCollector).Describe(ch)
		close(ch)
	}()

	descs := make([]*prometheus.Desc, 0)
	for desc := range ch {
		descs = append(descs, desc)
	}

	require.Len(t, descs, 6)

	assert.Equal(t, "Desc{fqName: \"test-system_test-subsystem_load_count\", help: \"The number of objects successfully loaded from the store\", constLabels: {}, variableLabels: {}}", descs[0].String())
	assert.Equal(t, "Desc{fqName: \"test-system_test-subsystem_store_count\", help: \"The number of objects successfully stored in the store\", constLabels: {}, variableLabels: {}}", descs[1].String())
	assert.Equal(t, "Desc{fqName: \"test-system_test-subsystem_load_err_count\", help: \"The number of objects that failed to load from the store\", constLabels: {}, variableLabels: {}}", descs[2].String())
	assert.Equal(t, "Desc{fqName: \"test-system_test-subsystem_store_err_count\", help: \"The number of objects that failed to store in the store\", constLabels: {}, variableLabels: {}}", descs[3].String())
	assert.Equal(t, "Desc{fqName: \"test-system_test-subsystem_load_duration_seconds\", help: \"The duration of the load method (in seconds)\", constLabels: {}, variableLabels: {}}", descs[4].String())
	assert.Equal(t, "Desc{fqName: \"test-system_test-subsystem_store_duration_seconds\", help: \"The duration of the store method (in seconds)\", constLabels: {}, variableLabels: {}}", descs[5].String())

	chMetrics := make(chan prometheus.Metric)
	go func() {
		metricsStore.(svc.MetricsCollector).Collect(chMetrics)
		close(chMetrics)
	}()

	metrics := make([]prometheus.Metric, 0)
	for metric := range chMetrics {
		metrics = append(metrics, metric)
	}

	require.Len(t, metrics, 6)
}

func TestWithLog(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockStore(ctrl)
	logStore := store.WithLog(mockStore)

	require.Implements(t, (*store.Store)(nil), logStore)

	ctx := context.TODO()
	reader := io.NopCloser(strings.NewReader("test-value"))
	headers := new(store.Headers)

	mockStore.EXPECT().Store(ctx, "test-key", reader, headers).Return(nil)
	err := logStore.Store(ctx, "test-key", reader, headers)
	require.NoError(t, err)

	mockStore.EXPECT().Load(ctx, "test-key").Return(reader, headers, nil)
	resReader, resHeaders, err := logStore.Load(ctx, "test-key")
	require.NoError(t, err)
	require.Equal(t, reader, resReader)
	require.Equal(t, headers, resHeaders)
}
