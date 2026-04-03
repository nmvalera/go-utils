package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/nmvalera/go-utils/app/svc"
	"github.com/nmvalera/go-utils/common"
	"github.com/nmvalera/go-utils/log"
	kkrthttp "github.com/nmvalera/go-utils/net/http"
	"github.com/nmvalera/go-utils/tag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func newTestApp(t *testing.T) *App {
	cfg := &Config{
		MainEntrypoint: &kkrthttp.EntrypointConfig{
			HTTP: &kkrthttp.ServerConfig{
				ReadTimeout:       common.Ptr(time.Second),
				ReadHeaderTimeout: common.Ptr(time.Second),
				WriteTimeout:      common.Ptr(time.Second),
				IdleTimeout:       common.Ptr(time.Second),
			},
			Net: &kkrthttp.ListenConfig{
				KeepAlive: common.Ptr(time.Second),
			},
		},
		HealthzEntrypoint: &kkrthttp.EntrypointConfig{
			HTTP: &kkrthttp.ServerConfig{
				ReadTimeout:       common.Ptr(time.Second),
				ReadHeaderTimeout: common.Ptr(time.Second),
				WriteTimeout:      common.Ptr(time.Second),
				IdleTimeout:       common.Ptr(time.Second),
			},
			Net: &kkrthttp.ListenConfig{
				KeepAlive: common.Ptr(time.Second),
			},
		},
		HealthzServer: &HealthzServerConfig{
			LivenessPath:  common.Ptr("/live"),
			ReadinessPath: common.Ptr("/ready"),
			MetricsPath:   common.Ptr("/metrics"),
		},
		Log:          log.DefaultConfig(),
		StartTimeout: common.Ptr(15 * time.Second),
		StopTimeout:  common.Ptr(15 * time.Second),
	}
	app, err := NewApp(
		cfg,
		WithLogger(zap.NewNop()),
		WithName("test"),
		WithVersion("1.0.0"),
	)
	require.NoError(t, err)
	return app
}

func TestAppProvide(t *testing.T) {
	var testCase = []struct {
		desc        string
		constructor func() (any, error)
		expected    any
		expectErr   bool
	}{
		{
			desc: "string",
			constructor: func() (any, error) {
				return "test", nil
			},
			expected: "test",
		},
		{
			desc: "int",
			constructor: func() (any, error) {
				return 1, nil
			},
			expected: 1,
		},
		{
			desc: "nil",
			constructor: func() (any, error) {
				return nil, nil
			},
			expected: nil,
		},
		{
			desc: "error",
			constructor: func() (any, error) {
				return nil, errors.New("error")
			},
			expectErr: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.desc, func(t *testing.T) {
			app := newTestApp(t)
			res := app.Provide("test", tc.constructor)
			err := app.Error()
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, res)
			}
		})
	}
}

func TestProvide(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		app := newTestApp(t)
		res := Provide(app, "test", func() (string, error) {
			return "test", nil
		})
		assert.Equal(t, res, "test")
	})

	t.Run("int", func(t *testing.T) {
		app := newTestApp(t)
		res := Provide(app, "test", func() (int, error) {
			return 1, nil
		})
		assert.Equal(t, res, 1)
	})

	t.Run("*string", func(t *testing.T) {
		app := newTestApp(t)
		res := Provide(app, "test", func() (*string, error) {
			return nil, nil
		})
		assert.Nil(t, res)
	})

	t.Run("*string#nil", func(t *testing.T) {
		app := newTestApp(t)
		res := Provide(app, "test", func() (*string, error) {
			return nil, nil
		})
		assert.Nil(t, res)
	})

	t.Run("error", func(t *testing.T) {
		app := newTestApp(t)
		res := Provide(app, "test", func() (error, error) {
			return errors.New("error"), nil
		})
		assert.Equal(t, errors.New("error"), res)
	})

	t.Run("interface", func(t *testing.T) {
		app := newTestApp(t)
		res := Provide(app, "test", func() (interface{}, error) {
			return "test", nil
		})
		assert.Equal(t, res, "test")
	})

	t.Run("interface#nil", func(t *testing.T) {
		app := newTestApp(t)
		res := Provide(app, "test", func() (interface{}, error) {
			return nil, nil
		})
		assert.Nil(t, res)
	})
}

type testService struct {
	start chan error
	stop  chan error
}

func (s *testService) Start(_ context.Context) error {
	return <-s.start
}

func (s *testService) Stop(_ context.Context) error {
	return <-s.stop
}

func TestAppNoDeps(t *testing.T) {
	start, stop := make(chan error), make(chan error)
	defer close(start)
	defer close(stop)

	testApp := func() *App {
		app := newTestApp(t)
		_ = Provide(app, "test", func() (*testService, error) {
			return &testService{
				start: start,
				stop:  stop,
			}, nil
		})
		return app
	}

	recStart, recStop := make(chan error), make(chan error)
	defer close(recStart)
	defer close(recStop)

	t.Run("no errors", func(t *testing.T) {
		app := testApp()
		require.Equal(t, app.services["test"].Status(), Constructed)

		go func() {
			recStart <- app.Start(context.Background())
		}()
		time.Sleep(100 * time.Millisecond) // wait for the service to start
		assert.Equal(t, app.services["test"].Status(), Starting)

		// Trigger start
		start <- nil
		assert.Nil(t, <-recStart)
		assert.Equal(t, app.services["test"].Status(), Running)

		go func() {
			recStop <- app.Stop(context.Background())
		}()
		time.Sleep(100 * time.Millisecond) // wait for the service to start
		assert.Equal(t, app.services["test"].Status(), Stopping)

		// Trigger stop
		stop <- nil
		assert.Nil(t, <-recStop)
		assert.Equal(t, app.services["test"].Status(), Stopped)
	})

	t.Run("error on start", func(t *testing.T) {
		app := testApp()
		go func() {
			recStart <- app.Start(context.Background())
		}()

		start <- errors.New("error on start")
		startErr := <-recStart
		require.IsType(t, startErr, &ServiceError{})
		assert.Equal(t, startErr.(*ServiceError).directErr, errors.New("error on start"))
		assert.Equal(t, app.services["test"].Status(), Error)
	})

	t.Run("error on stop", func(t *testing.T) {
		app := testApp()
		go func() {
			recStart <- app.Start(context.Background())
		}()
		start <- nil
		<-recStart

		go func() {
			recStop <- app.Stop(context.Background())
		}()
		stop <- errors.New("error on stop")
		stopErr := <-recStop
		require.IsType(t, stopErr, &ServiceError{})
		assert.Equal(t, stopErr.(*ServiceError).directErr, errors.New("error on stop"))
		assert.Equal(t, app.services["test"].Status(), Error)
	})
}

func TestAppWithDeps(t *testing.T) {
	app := newTestApp(t)
	startMain, stopMain, startDep, stopDep := make(chan error), make(chan error), make(chan error), make(chan error)
	defer close(startMain)
	defer close(stopMain)
	defer close(startDep)
	defer close(stopDep)

	_ = Provide(app, "main", func() (*testService, error) {
		_ = Provide(app, "dep", func() (*testService, error) {
			return &testService{
				start: startDep,
				stop:  stopDep,
			}, nil
		})
		return &testService{
			start: startMain,
			stop:  stopMain,
		}, nil
	})

	// Test dependency tree
	assert.Equal(t, app.services["main"].deps["dep"], app.services["dep"])
	assert.Equal(t, app.services["dep"].depsOf["main"], app.services["main"])

	recStart, recStop := make(chan error), make(chan error)
	defer close(recStart)
	defer close(recStop)

	go func() {
		recStart <- app.Start(context.Background())
	}()
	startDep <- nil
	startMain <- nil
	assert.Nil(t, <-recStart)

	go func() {
		recStop <- app.Stop(context.Background())
	}()
	stopMain <- nil
	stopDep <- nil
	assert.Nil(t, <-recStop)
}

func TestServiceCanBeRetrievedBeforeConstruction(t *testing.T) {
	app := newTestApp(t)
	myDep := func() string {
		return Provide(app, "dep", func() (string, error) { return "test-dep", nil })
	}
	myMain := func() string {
		return Provide(app, "main", func() (string, error) { return fmt.Sprintf("test-main with dep: %s", myDep()), nil })
	}

	assert.Equal(t, myMain(), "test-main with dep: test-dep")
	assert.Equal(t, myDep(), "test-dep")
}

func TestServiceError(t *testing.T) {
	rootSvc := newService("svc", nil)
	dep1 := newService("dep1", nil)
	dep2 := newService("dep2", nil)
	dep11 := newService("dep11", nil)
	dep21 := newService("dep21", nil)

	svcErr := rootSvc.fail(errors.New("error on svc"))
	dep1Err := dep1.fail(nil)
	dep2Err := dep2.fail(errors.New("error on dep2"))
	dep1Err.addDepsErr(dep11.fail(errors.New("error on dep11")))
	dep2Err.addDepsErr(dep21.fail(nil))

	svcErr.addDepsErr(dep1Err)
	svcErr.addDepsErr(dep2Err)

	expected := `service "svc": error on svc
>service "dep1"
>>service "dep11": error on dep11
>service "dep2": error on dep2
>>service "dep21"`
	assert.Equal(t, svcErr.Error(), expected)
}

func TestAppServers(t *testing.T) {
	app := newTestApp(t)
	require.NoError(t, app.Error())

	app.Provide("top", func() (any, error) {
		app.EnableMainEntrypoint()
		app.EnableHealthzEntrypoint()
		return nil, nil
	})

	err := app.Start(context.Background())
	require.NoError(t, err)

	// Check main server is running
	require.NotNil(t, app.main)
	mainAddr := app.main.Addr()
	require.NotEmpty(t, mainAddr)

	conn, err := net.Dial("tcp", mainAddr)
	require.NoError(t, err)
	_ = conn.Close()

	// Check main server is running
	require.NotNil(t, app.healthz)
	healthzAddr := app.healthz.Addr()
	require.NotEmpty(t, healthzAddr)

	conn, err = net.Dial("tcp", healthzAddr)
	require.NoError(t, err)
	_ = conn.Close()

	// Check healthz server is running
	err = app.Stop(context.Background())
	require.NoError(t, err)
}

type checkableService struct {
	err error
}

func (s *checkableService) Ready(_ context.Context) error {
	return s.err
}

func TestHealthChecks(t *testing.T) {
	app := newTestApp(t)
	require.NoError(t, app.Error())

	checkable := new(checkableService)
	Provide(app, "checkable", func() (*checkableService, error) {
		app.EnableHealthzEntrypoint()
		return checkable, nil
	})

	err := app.Start(context.Background())
	require.NoError(t, err)

	require.NotNil(t, app.healthz)
	healthAddr := app.healthz.Addr()
	require.NotEmpty(t, healthAddr)

	// Test live check
	req, err := http.NewRequest("GET", "http://"+healthAddr+"/live", http.NoBody)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	// Test ready check
	req, err = http.NewRequest("GET", "http://"+healthAddr+"/ready", http.NoBody)
	require.NoError(t, err)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	// Test ready check with error
	checkable.err = errors.New("test error")

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, resp.StatusCode, http.StatusServiceUnavailable)
}

type metricsService struct {
	name  string
	count prometheus.Counter
}

func (s *metricsService) SetMetrics(appName, subsystem string, _ ...*tag.Tag) {
	s.count = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: appName,
		Subsystem: subsystem,
		Name:      fmt.Sprintf("%s_count", s.name),
	})
}

func (s *metricsService) incr() {
	s.count.Inc()
}

func (s *metricsService) Describe(ch chan<- *prometheus.Desc) {
	ch <- s.count.Desc()
}

func (s *metricsService) Collect(ch chan<- prometheus.Metric) {
	ch <- s.count
}

func TestMetrics(t *testing.T) {
	app := newTestApp(t)
	_ = WithName("testApp")(app)

	require.NoError(t, app.Error())

	metrics := &metricsService{
		name: "A",
	}
	app.Provide("test-wo-cfg", func() (any, error) {
		app.EnableHealthzEntrypoint()
		app.Provide(
			"test-w-cfg",
			func() (any, error) {
				return &metricsService{
					name: "B",
				}, nil
			},
			WithComponentName("subsystem"),
		)
		return metrics, nil
	})

	err := app.Start(context.Background())
	require.NoError(t, err)

	require.NotNil(t, app.healthz)
	healthAddr := app.healthz.Addr()
	require.NotEmpty(t, healthAddr)

	// Test metrics endpoint
	req, err := http.NewRequest("GET", "http://"+healthAddr+"/metrics", http.NoBody)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	// Test collectors are registered with correct labels
	families, err := app.prometheus.Gather()
	require.NoError(t, err)
	familyCount := len(families)
	assert.GreaterOrEqual(t, familyCount, 2)
	assert.Equal(t, "testApp_subsystem_B_count", families[familyCount-2].GetName())
	assert.Equal(t, "testApp_test_wo_cfg_A_count", families[familyCount-1].GetName())

	// Test metrics are updated
	assert.Equal(t, float64(0), families[familyCount-1].GetMetric()[0].GetCounter().GetValue())
	metrics.incr()
	metrics.incr()
	metrics.incr()

	families, err = app.prometheus.Gather()
	require.NoError(t, err)
	assert.Equal(t, float64(3), families[familyCount-1].GetMetric()[0].GetCounter().GetValue())

	err = app.Stop(context.Background())
	require.NoError(t, err)
}

type healthzAPIService struct{}

func (s *healthzAPIService) RegisterHealthzHandler(router *mux.Router) {
	router.Path("/main-test").Methods(http.MethodGet).Handler(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
}
func TestHealthzAPI(t *testing.T) {
	app := newTestApp(t)
	require.NoError(t, app.Error())

	healthz := &healthzAPIService{}
	Provide(app, "healthz", func() (*healthzAPIService, error) {
		app.EnableHealthzEntrypoint()
		return healthz, nil
	})

	err := app.Start(context.Background())
	require.NoError(t, err)

	require.NotNil(t, app.healthz)
	healthAddr := app.healthz.Addr()
	require.NotEmpty(t, healthAddr)

	req, err := http.NewRequest("GET", "http://"+healthAddr+"/main-test", http.NoBody)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	err = app.Stop(context.Background())
	require.NoError(t, err)
}

type middlewareService struct{}

func (s *middlewareService) RegisterMiddleware(chain alice.Chain) alice.Chain {
	return chain.Append(func(_ http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]string{"message": "ok"})
		})
	})
}

func TestMiddlewareService(t *testing.T) {
	app := newTestApp(t)
	require.NoError(t, app.Error())

	middleware := &middlewareService{}
	Provide(app, "middleware", func() (*middlewareService, error) {
		app.EnableMainEntrypoint()
		return middleware, nil
	})

	err := app.Start(context.Background())
	require.NoError(t, err)

	require.NotNil(t, app.main)
	mainAddr := app.main.Addr()
	require.NotEmpty(t, mainAddr)

	req, err := http.NewRequest("GET", "http://"+mainAddr, http.NoBody)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, resp.StatusCode, http.StatusOK)
	var body map[string]string
	err = json.NewDecoder(resp.Body).Decode(&body)
	require.NoError(t, err)
	assert.Equal(t, body, map[string]string{"message": "ok"})

	err = app.Stop(context.Background())
	require.NoError(t, err)
}

type healthzService struct{}

func (s *healthzService) RegisterHealthzHandler(router *mux.Router) {
	router.Path("/debug-test").Methods(http.MethodGet).Handler(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
}

type runContextService struct {
	start chan error
	stop  chan error

	runCtx    context.Context
	runCtxMux sync.Mutex
}

func (s *runContextService) SetRunContext(ctx context.Context) {
	s.runCtxMux.Lock()
	defer s.runCtxMux.Unlock()
	s.runCtx = ctx
}

func (s *runContextService) getRunCtx() context.Context {
	s.runCtxMux.Lock()
	defer s.runCtxMux.Unlock()
	return s.runCtx
}

func (s *runContextService) Start(_ context.Context) error {
	return <-s.start
}

func (s *runContextService) Stop(_ context.Context) error {
	return <-s.stop
}

func TestRunContext(t *testing.T) {
	t.Run("set before start and alive during run", func(t *testing.T) {
		app := newTestApp(t)

		start, stop := make(chan error), make(chan error)
		defer close(start)
		defer close(stop)

		rcSvc := &runContextService{start: start, stop: stop}
		Provide(app, "test", func() (*runContextService, error) {
			return rcSvc, nil
		})

		recStart := make(chan error)
		go func() {
			recStart <- app.Start(context.Background())
		}()

		// Run context should be set before Start completes (it's set before Start is called)
		// We need to let Start proceed
		start <- nil
		require.NoError(t, <-recStart)

		// Run context should be set and not canceled
		runCtx := rcSvc.getRunCtx()
		require.NotNil(t, runCtx)
		assert.NoError(t, runCtx.Err())

		// Verify tags are present
		tags := tag.FromContext(runCtx)
		assert.NotEmpty(t, tags)

		recStop := make(chan error)
		go func() {
			recStop <- app.Stop(context.Background())
		}()

		// Run context should still be alive during Stop
		time.Sleep(50 * time.Millisecond)
		assert.NoError(t, runCtx.Err())

		stop <- nil
		require.NoError(t, <-recStop)

		// Run context should be canceled after Stop completes
		assert.Error(t, runCtx.Err())
		assert.Equal(t, context.Canceled, runCtx.Err())
	})

	t.Run("canceled on start failure", func(t *testing.T) {
		app := newTestApp(t)

		start, stop := make(chan error), make(chan error)
		defer close(start)
		defer close(stop)

		rcSvc := &runContextService{start: start, stop: stop}
		Provide(app, "test", func() (*runContextService, error) {
			return rcSvc, nil
		})

		recStart := make(chan error)
		go func() {
			recStart <- app.Start(context.Background())
		}()

		start <- errors.New("start failed")
		require.Error(t, <-recStart)

		// Run context should be canceled after start failure
		runCtx := rcSvc.getRunCtx()
		require.NotNil(t, runCtx)
		assert.Error(t, runCtx.Err())
	})

	t.Run("non RunContextAware service unaffected", func(t *testing.T) {
		app := newTestApp(t)

		start, stop := make(chan error), make(chan error)
		defer close(start)
		defer close(stop)

		testSvc := &testService{start: start, stop: stop}
		Provide(app, "test", func() (*testService, error) {
			return testSvc, nil
		})

		recStart := make(chan error)
		go func() {
			recStart <- app.Start(context.Background())
		}()
		start <- nil
		require.NoError(t, <-recStart)

		recStop := make(chan error)
		go func() {
			recStop <- app.Stop(context.Background())
		}()
		stop <- nil
		require.NoError(t, <-recStop)
	})

	t.Run("config tags in run context", func(t *testing.T) {
		cfg := &Config{
			MainEntrypoint: &kkrthttp.EntrypointConfig{
				HTTP: &kkrthttp.ServerConfig{
					ReadTimeout:       common.Ptr(time.Second),
					ReadHeaderTimeout: common.Ptr(time.Second),
					WriteTimeout:      common.Ptr(time.Second),
					IdleTimeout:       common.Ptr(time.Second),
				},
				Net: &kkrthttp.ListenConfig{
					KeepAlive: common.Ptr(time.Second),
				},
			},
			HealthzEntrypoint: &kkrthttp.EntrypointConfig{
				HTTP: &kkrthttp.ServerConfig{
					ReadTimeout:       common.Ptr(time.Second),
					ReadHeaderTimeout: common.Ptr(time.Second),
					WriteTimeout:      common.Ptr(time.Second),
					IdleTimeout:       common.Ptr(time.Second),
				},
				Net: &kkrthttp.ListenConfig{
					KeepAlive: common.Ptr(time.Second),
				},
			},
			HealthzServer: &HealthzServerConfig{
				LivenessPath:  common.Ptr("/live"),
				ReadinessPath: common.Ptr("/ready"),
				MetricsPath:   common.Ptr("/metrics"),
			},
			Log:          log.DefaultConfig(),
			StartTimeout: common.Ptr(15 * time.Second),
			StopTimeout:  common.Ptr(15 * time.Second),
			Tags: map[string]any{
				"env":     "production",
				"cluster": "us-east-1",
				"nested": map[string]any{
					"a": "b",
					"c": 1,
					"d": 1.0,
					"e": true,
				},
			},
		}
		app, err := NewApp(cfg, WithLogger(zap.NewNop()), WithName("myapp"), WithVersion("1.0.0"))
		require.NoError(t, err)

		start, stop := make(chan error), make(chan error)
		defer close(start)
		defer close(stop)

		rcSvc := &runContextService{start: start, stop: stop}
		Provide(app, "test", func() (*runContextService, error) {
			return rcSvc, nil
		})

		recStart := make(chan error)
		go func() {
			recStart <- app.Start(context.Background())
		}()
		start <- nil
		require.NoError(t, <-recStart)

		runCtx := rcSvc.getRunCtx()
		require.NotNil(t, runCtx)

		// Verify all tags are present in the run context
		tags := tag.FromContext(runCtx)
		tagMap := make(map[string]string)
		for _, t := range tags {
			tagMap[string(t.Key)] = t.Value.String()
		}
		assert.Equal(t, "myapp", tagMap["app"])
		assert.Equal(t, "1.0.0", tagMap["version"])
		assert.Equal(t, "production", tagMap["env"])
		assert.Equal(t, "us-east-1", tagMap["cluster"])
		// Service component tag should also be present
		assert.Equal(t, "test", tagMap["component"])

		recStop := make(chan error)
		go func() {
			recStop <- app.Stop(context.Background())
		}()
		stop <- nil
		require.NoError(t, <-recStop)
	})

	t.Run("canceled on stop context timeout", func(t *testing.T) {
		app := newTestApp(t)

		start := make(chan error)
		defer close(start)

		rcSvc := &runContextService{start: start, stop: make(chan error)}
		Provide(app, "test", func() (*runContextService, error) {
			return rcSvc, nil
		})

		recStart := make(chan error)
		go func() {
			recStart <- app.Start(context.Background())
		}()
		start <- nil
		require.NoError(t, <-recStart)

		runCtx := rcSvc.getRunCtx()
		require.NotNil(t, runCtx)
		assert.NoError(t, runCtx.Err())

		// Stop with a context that times out (Stop will block because stop channel is never fed)
		stopCtx, stopCancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer stopCancel()

		err := app.Stop(stopCtx)
		require.Error(t, err)

		// Run context should be canceled after stop context timeout
		assert.Error(t, runCtx.Err())
		assert.Equal(t, context.Canceled, runCtx.Err())
	})
}

// childTaggedService implements svc.Taggable (tags applied at construct) and can enrich a parent run context.
type childTaggedService struct {
	svc.Tagged
}

type parentRunContextService struct {
	svc.RunContext
	child *childTaggedService

	start chan error
	stop  chan error
}

func (p *parentRunContextService) Start(_ context.Context) error {
	return <-p.start
}

func (p *parentRunContextService) Stop(_ context.Context) error {
	return <-p.stop
}

func TestRunContext_parentRunCtxChildChainedTags(t *testing.T) {
	app := newTestApp(t)

	start, stop := make(chan error), make(chan error)
	defer close(start)
	defer close(stop)

	var parent *parentRunContextService
	Provide(app, "parent", func() (*parentRunContextService, error) {
		ch := Provide(app, "child", func() (*childTaggedService, error) {
			return &childTaggedService{}, nil
		}, WithComponentNameChained(true))
		parent = &parentRunContextService{
			child: ch,
			start: start,
			stop:  stop,
		}
		return parent, nil
	})

	recStart := make(chan error)
	go func() {
		recStart <- app.Start(context.Background())
	}()
	start <- nil
	require.NoError(t, <-recStart)

	require.NotNil(t, parent)
	assert.Equal(t, "parent", getComponentTag(parent.Context()))
	assert.Equal(t, "parent.child", getComponentTag(parent.child.Context(parent.Context())))
	assert.Equal(t, "child", getComponentTag(parent.child.Context(context.Background())))

	recStop := make(chan error)
	go func() {
		recStop <- app.Stop(context.Background())
	}()
	stop <- nil
	require.NoError(t, <-recStop)
}

func getComponentTag(ctx context.Context) string {
	for _, t := range tag.FromContext(ctx) {
		if t.Key == "component" {
			return t.Value.String()
		}
	}
	return ""
}

func TestHealthzService(t *testing.T) {
	app := newTestApp(t)
	require.NoError(t, app.Error())

	healthz := &healthzService{}
	Provide(app, "healthz", func() (*healthzService, error) {
		app.EnableHealthzEntrypoint()
		return healthz, nil
	})

	err := app.Start(context.Background())
	require.NoError(t, err)

	require.NotNil(t, app.healthz)
	healthAddr := app.healthz.Addr()
	require.NotEmpty(t, healthAddr)

	req, err := http.NewRequest("GET", "http://"+healthAddr+"/debug-test", http.NoBody)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	err = app.Stop(context.Background())
	require.NoError(t, err)
}
