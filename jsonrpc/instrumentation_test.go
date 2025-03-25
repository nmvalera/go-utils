package jsonrpc_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kkrt-labs/go-utils/app/svc"
	kkrtgomock "github.com/kkrt-labs/go-utils/gomock"
	"github.com/kkrt-labs/go-utils/jsonrpc"
	jsonrpcmock "github.com/kkrt-labs/go-utils/jsonrpc/mock"
	"github.com/kkrt-labs/go-utils/tag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestWithTags(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCli := jsonrpcmock.NewMockClient(ctrl)

	taggedCli := jsonrpc.WithTags(mockCli)
	require.Implementsf(t, (*svc.Taggable)(nil), taggedCli, "taggedCli should implement svc.Taggable")

	taggedCli.(svc.Taggable).WithTags(tag.Key("test-key").String("test-value"))

	req := &jsonrpc.Request{
		ID:      "test-id",
		Method:  "test-method",
		Params:  []interface{}{"test-param"},
		Version: "2.0",
	}
	res := "test"

	validateCtx := func(ctx context.Context) error {
		return tag.ExpectTagsOnContext(
			ctx,
			tag.Key("test-key").String("test-value"),
			tag.Key("req.method").String("test-method"),
			tag.Key("req.version").String("2.0"),
			tag.Key("req.params").Object([]interface{}{"test-param"}),
			tag.Key("req.id").Object("test-id"),
		)
	}

	mockCli.EXPECT().Call(kkrtgomock.ContextMatcher(validateCtx), req, res).Return(nil)
	err := taggedCli.Call(context.TODO(), req, res)
	require.NoError(t, err)
}

func TestWithLog(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCli := jsonrpcmock.NewMockClient(ctrl)
	logCli := jsonrpc.WithLog()(mockCli)

	req := new(jsonrpc.Request)
	res := new(string)
	ctx := context.TODO()
	mockCli.EXPECT().Call(ctx, req, res).Return(nil)
	err := logCli.Call(ctx, req, res)
	require.NoError(t, err)
}

func TestWithMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCli := jsonrpcmock.NewMockClient(ctrl)

	metricsCli := jsonrpc.WithMetrics(mockCli)
	require.Implementsf(t, (*svc.Metricable)(nil), metricsCli, "metricsCli should implement svc.Metricable")
	assert.Implementsf(t, (*svc.MetricsCollector)(nil), metricsCli, "metricsCli should implement svc.MetricsCollector")
	metricsCli.(svc.Metricable).SetMetrics("test-system", "test-subsystem")

	req := &jsonrpc.Request{
		Method: "test-method",
	}
	res := "test"

	// One call without error
	mockCli.EXPECT().Call(gomock.Any(), req, res).Return(nil)
	err := metricsCli.Call(context.TODO(), req, res)
	require.NoError(t, err)

	// One call with error
	mockCli.EXPECT().Call(gomock.Any(), req, res).Return(errors.New("test error"))
	err = metricsCli.Call(context.TODO(), req, res)
	require.Error(t, err)

	ch := make(chan *prometheus.Desc)
	go func() {
		metricsCli.(svc.MetricsCollector).Describe(ch)
		close(ch)
	}()

	descs := make([]*prometheus.Desc, 0)
	for desc := range ch {
		descs = append(descs, desc)
	}

	require.Len(t, descs, 3)
	assert.Equal(t, "Desc{fqName: \"test-system_test-subsystem_request_duration_seconds\", help: \"The duration of requests in seconds (per method)\", constLabels: {}, variableLabels: {method}}", descs[0].String())
	assert.Equal(t, "Desc{fqName: \"test-system_test-subsystem_requests_total\", help: \"The total number of requests (per method)\", constLabels: {}, variableLabels: {method}}", descs[1].String())
	assert.Equal(t, "Desc{fqName: \"test-system_test-subsystem_requests_errors\", help: \"The number of requests that failed (per method)\", constLabels: {}, variableLabels: {method}}", descs[2].String())

	chMetrics := make(chan prometheus.Metric)
	go func() {
		metricsCli.(svc.MetricsCollector).Collect(chMetrics)
		close(chMetrics)
	}()

	metrics := make([]prometheus.Metric, 0)
	for metric := range chMetrics {
		metrics = append(metrics, metric)
	}

	assert.Len(t, metrics, 3)
}
