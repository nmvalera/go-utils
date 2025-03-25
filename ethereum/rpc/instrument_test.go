package rpc_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/kkrt-labs/go-utils/app/svc"
	"github.com/kkrt-labs/go-utils/ethereum/rpc"
	"github.com/kkrt-labs/go-utils/ethereum/rpc/mock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestWithMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockClient(ctrl)
	metricsClient := rpc.WithMetrics(mockClient, rpc.WithFetchInterval(time.Hour)) // Set long interval so fetch happens only once during the test

	require.Implements(t, (*rpc.Client)(nil), metricsClient)
	require.Implements(t, (*svc.Runnable)(nil), metricsClient)
	require.Implements(t, (*svc.Metricable)(nil), metricsClient)
	require.Implements(t, (*svc.MetricsCollector)(nil), metricsClient)
	require.Implements(t, (*svc.Taggable)(nil), metricsClient)

	metricsClient.(svc.Metricable).SetMetrics("test-system", "test-subsystem")

	mockClient.EXPECT().ChainID(gomock.Any()).Return(big.NewInt(1), nil)
	mockClient.EXPECT().HeaderByNumber(gomock.Any(), nil).Return(&gethtypes.Header{Number: big.NewInt(1)}, nil)

	err := metricsClient.(svc.Runnable).Start(context.Background())
	require.NoError(t, err)

	time.Sleep(100 * time.Millisecond) // Sleeps for the fetch to happen

	err = metricsClient.(svc.Runnable).Stop(context.Background())
	require.NoError(t, err)

	ch := make(chan *prometheus.Desc)
	go func() {
		metricsClient.(svc.MetricsCollector).Describe(ch)
		close(ch)
	}()

	descs := make([]*prometheus.Desc, 0)
	for desc := range ch {
		descs = append(descs, desc)
	}

	require.Len(t, descs, 2)
	assert.Equal(t, "Desc{fqName: \"test-system_test-subsystem_chain_id\", help: \"The chain ID of the network\", constLabels: {}, variableLabels: {}}", descs[0].String())
	assert.Equal(t, "Desc{fqName: \"test-system_test-subsystem_chain_head\", help: \"The block number of the chain head\", constLabels: {}, variableLabels: {}}", descs[1].String())

	chMetrics := make(chan prometheus.Metric)
	go func() {
		metricsClient.(svc.MetricsCollector).Collect(chMetrics)
		close(chMetrics)
	}()

	metrics := make([]prometheus.Metric, 0)
	for metric := range chMetrics {
		metrics = append(metrics, metric)
	}

	require.Len(t, metrics, 2)
}

func TestWithCheck(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockClient(ctrl)
	metricsClient := rpc.WithCheck(mockClient)

	require.Implements(t, (*rpc.Client)(nil), metricsClient)
	require.Implements(t, (*svc.Checkable)(nil), metricsClient)

	mockClient.EXPECT().ChainID(gomock.Any()).Return(big.NewInt(1), nil)

	err := metricsClient.(svc.Checkable).Ready(context.Background())
	require.NoError(t, err)
}
