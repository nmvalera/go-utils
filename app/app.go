package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/MadAppGang/httplog"
	httplogzap "github.com/MadAppGang/httplog/zap"
	"github.com/hellofresh/health-go/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/nmvalera/go-utils/app/svc"
	"github.com/nmvalera/go-utils/log"
	kkrthttp "github.com/nmvalera/go-utils/net/http"
	"github.com/nmvalera/go-utils/tag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type App struct {
	cfg *Config

	name    string
	version string

	services map[string]*service

	top     *service
	current *service

	done chan os.Signal

	logger               *zap.Logger
	replaceGlobalLoggers bool
	resetGlobalLoggers   func()

	mainMiddleware alice.Chain
	main           *kkrthttp.Entrypoint
	mainRouter     *httprouter.Router

	healthz       *kkrthttp.Entrypoint
	healthzRouter *httprouter.Router

	liveHealth  *health.Health
	readyHealth *health.Health

	prometheus *prometheus.Registry
}

func NewApp(cfg *Config, opts ...Option) (*App, error) {
	app := &App{
		cfg:            cfg,
		services:       make(map[string]*service),
		done:           make(chan os.Signal),
		logger:         zap.NewNop(),
		mainMiddleware: alice.New(),
		mainRouter:     httprouter.New(),
		healthzRouter:  httprouter.New(),
		prometheus:     prometheus.NewRegistry(),
	}

	logger, err := cfg.Log.ZapConfig().Build()
	if err != nil {
		return nil, err
	}
	app.logger = logger

	for _, opt := range opts {
		if err := opt(app); err != nil {
			return nil, err
		}
	}

	if app.replaceGlobalLoggers {
		app.replaceLoggers()
	}

	app.liveHealth = newHealth(app)
	app.readyHealth = newHealth(app)

	app.registerBaseMetrics()

	return app, nil
}

func (app *App) replaceLoggers() {
	if app.resetGlobalLoggers != nil {
		app.resetGlobalLoggers = zap.ReplaceGlobals(app.logger)
	}
}

func (app *App) resetLoggers() {
	if app.resetGlobalLoggers != nil {
		app.resetGlobalLoggers()
	}
}

func (app *App) Provide(id string, constructor func() (any, error), opts ...ServiceOption) any {
	if strings.HasPrefix(id, "system.") {
		panic(fmt.Sprintf("invalid service id: %q (system.* is reserved for internal use)", id))
	}

	return app.provide(id, constructor, opts...)
}

func (app *App) provide(id string, constructor func() (any, error), opts ...ServiceOption) any {
	if id == "" {
		id = reflect.TypeOf(constructor).Out(0).String()
	}

	if srvc, ok := app.services[id]; ok {
		app.current.addDep(srvc) // current can not be nil here
		return srvc.value
	}

	srvc := app.createService(id, constructor, opts...)
	app.services[id] = srvc

	return srvc.value
}

func (app *App) createService(id string, constructor func() (any, error), opts ...ServiceOption) *service {
	previous := app.current
	srvc := newService(id, constructor, opts...)
	srvc.app = app

	app.current = srvc // set the current service pointer
	srvc.construct()   // construct can perform calls to Provide moving the current service pointer
	if previous != nil {
		previous.addDep(srvc)
	} else {
		app.top = srvc
	}

	app.current = previous // restore the current service pointer

	return srvc
}

func Provide[T any](app *App, id string, constructor func() (T, error), opts ...ServiceOption) T {
	if strings.HasPrefix(id, "system.") {
		panic(fmt.Sprintf("invalid service id: %q (system.* is reserved for internal use)", id))
	}

	return provide(app, id, constructor, opts...)
}

func provide[T any](app *App, id string, constructor func() (T, error), opts ...ServiceOption) T {
	val := app.provide(id, func() (any, error) {
		return constructor()
	}, opts...)

	v := reflect.ValueOf(val)
	if v.Kind() == reflect.Invalid {
		// return zero value for the type T
		var zero T
		return zero
	}

	return val.(T)
}

func (app *App) Error() error {
	if app.top == nil || app.top.err == nil {
		return nil
	}
	return app.top.err
}

func (app *App) getLogger(component string) *zap.Logger {
	if component == "" {
		return app.logger.With(zap.String("component", "system"))
	}
	return app.logger.With(zap.String("component", component))
}

func (app *App) Context(ctx context.Context) context.Context {
	return app.context(ctx)
}

func (app *App) context(ctx context.Context) context.Context {
	return log.WithLogger(ctx, app.logger)
}

func (app *App) Start(ctx context.Context) error {
	logger := app.getLogger("")
	logger.Info("System starting...")
	err := app.start(app.context(ctx))
	if err != nil {
		logger.Error("System failed to start", zap.Error(err))
		return err
	}
	logger.Info("System successfully started")
	return nil
}

func (app *App) start(ctx context.Context) error {
	if app.top == nil {
		return fmt.Errorf("no service constructed yet")
	}

	if app.top.err != nil {
		return app.top.err
	}

	app.replaceLoggers()
	app.setHandlers()
	if err := app.top.start(ctx); err != nil {
		app.resetLoggers()
		return err
	}

	return nil
}

func (app *App) Stop(ctx context.Context) error {
	logger := app.getLogger("")
	logger.Info("System stopping...")
	err := app.stop(app.context(ctx))
	if err != nil {
		logger.Error("System failed to stop", zap.Error(err))
		return err
	}
	logger.Info("System successfully stopped")
	return nil
}

func (app *App) stop(ctx context.Context) error {
	defer app.resetLoggers()

	if app.top == nil {
		return fmt.Errorf("no service constructed yet")
	}

	if err := app.top.stop(ctx); err != nil {
		return err
	}

	return nil
}

func (app *App) Run(ctx context.Context) error {
	err := app.Start(ctx)
	if err != nil {
		return err
	}

	app.listenSignals()

	sig := <-app.done
	app.getLogger("").Warn("Received signal", zap.String("signal", sig.String()))

	app.stopListeningSignals()

	return app.Stop(ctx)
}

func (app *App) listenSignals() {
	signal.Notify(app.done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
}

func (app *App) stopListeningSignals() {
	signal.Stop(app.done)
}

func newHealth(app *App) *health.Health {
	h, _ := health.New(health.WithComponent(health.Component{Name: app.name, Version: app.version}))
	return h
}

func (app *App) MainEntrypoint() *kkrthttp.Entrypoint {
	return app.main
}

func (app *App) HealthzEntrypoint() *kkrthttp.Entrypoint {
	return app.healthz
}

func (app *App) EnableMainEntrypoint() {
	app.main = app.entrypoint("main", app.cfg.MainEntrypoint)
}

func (app *App) EnableHealthzEntrypoint() {
	app.healthz = app.entrypoint("healthz", app.cfg.HealthzEntrypoint)
}

func (app *App) entrypoint(name string, cfg *kkrthttp.EntrypointConfig) *kkrthttp.Entrypoint {
	return provide(app, fmt.Sprintf("system.%v.entrypoint", name), func() (*kkrthttp.Entrypoint, error) {
		return cfg.Entrypoint()
	})
}

func (app *App) setHandlers() {
	app.setMainHandler()
	app.setHealthzHandler()
}

func (app *App) setMainHandler() {
	if app.main != nil {
		h := app.instrumentMiddleware().Extend(app.mainMiddleware).Then(app.mainRouter)
		app.main.SetHandler(h)
	}
}

func (app *App) instrumentMiddleware() alice.Chain {
	return alice.New(
		// Log Requests on main router
		httplog.LoggerWithConfig(httplog.LoggerConfig{
			Formatter: httplogzap.ZapLogger(app.getLogger("system.main.entrypoint"), zapcore.InfoLevel, ""),
		}),
		// Instrument main router with prometheus metrics
		func(next http.Handler) http.Handler {
			return promhttp.InstrumentMetricHandler(app.prometheus, next)
		},
	)
}

func (app *App) setHealthzHandler() {
	app.healthzRouter.Handler(http.MethodGet, *app.cfg.HealthzServer.LivenessPath, app.liveHealth.Handler())
	app.healthzRouter.Handler(http.MethodGet, *app.cfg.HealthzServer.ReadinessPath, app.readyHealth.Handler())
	app.healthzRouter.Handler(http.MethodGet, *app.cfg.HealthzServer.MetricsPath, promhttp.HandlerFor(app.prometheus, promhttp.HandlerOpts{}))

	if app.healthz != nil {
		app.healthz.SetHandler(app.healthzRouter)
	}
}

func (app *App) registerBaseMetrics() {
	app.prometheus.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	app.prometheus.MustRegister(collectors.NewGoCollector())
}

type ServiceStatus uint32

const (
	Constructing ServiceStatus = iota
	Constructed
	Starting
	Running
	Stopping
	Stopped
	Error
)

type service struct {
	id string

	app *App

	constructor func() (any, error)
	value       any

	deps   map[string]*service
	depsOf map[string]*service

	mux    sync.RWMutex
	status atomic.Uint32
	err    *ServiceError

	startOnce sync.Once

	stopOnce sync.Once
	stopChan chan struct{}

	name          string
	tags          tag.Set
	healthConfig  *health.Config
	metricsPrefix string
}

func newService(id string, constructor func() (any, error), opts ...ServiceOption) *service {
	s := &service{
		id:           id,
		constructor:  constructor,
		deps:         make(map[string]*service),
		depsOf:       make(map[string]*service),
		stopChan:     make(chan struct{}),
		tags:         tag.EmptySet.WithTags(tag.Key("component").String(id)),
		healthConfig: new(health.Config),
	}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			_ = s.fail(err)
			return nil
		}
	}

	return s
}

func (s *service) context(ctx context.Context) context.Context {
	return tag.WithTags(ctx, s.tags...)
}

func (s *service) Name() string {
	return s.id
}

func (s *service) Status() ServiceStatus {
	return ServiceStatus(s.status.Load())
}

func (s *service) setStatus(status ServiceStatus) {
	s.status.Store(uint32(status))
}

func (s *service) setStatusWithLock(status ServiceStatus) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.setStatus(status)
}

func (s *service) fail(err error) *ServiceError {
	if svcErr, ok := err.(*ServiceError); ok {
		s.err = svcErr
	} else {
		s.err = &ServiceError{
			svc:       s,
			directErr: err,
		}
	}
	s.setStatus(Error)

	return s.err
}

func (s *service) failWithLock(err error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	_ = s.fail(err)
}

func (s *service) construct() {
	s.setStatus(Constructing)
	val, constructorErr := s.constructor()
	if constructorErr != nil {
		_ = s.fail(constructorErr)
		return
	}

	if t, ok := val.(svc.Taggable); ok {
		t.WithTags(s.tags...)
	}

	if t, ok := val.(svc.Checkable); ok {
		s.healthConfig.Check = s.wrapCheck(t.Ready)
	}

	if t, ok := val.(svc.API); ok {
		t.RegisterHandler(s.app.mainRouter)
	}

	if t, ok := val.(svc.Middleware); ok {
		s.app.mainMiddleware = t.RegisterMiddleware(s.app.mainMiddleware)
	}

	if t, ok := val.(svc.Healthz); ok {
		t.RegisterHealthzHandler(s.app.healthzRouter)
	}

	s.value = val
	if err := s.registerReadyCheck(); err != nil {
		_ = s.fail(err)
		return
	}

	s.setMetrics()

	s.setStatus(Constructed)
}

// sanitizeMetricName sanitizes a name by replacing all non-alphanumeric characters with underscores
// except "_" and ":"
func sanitizeMetricName(name string) string {
	re := regexp.MustCompile("[^a-zA-Z0-9_:]")
	return re.ReplaceAllString(name, "_")
}

func (s *service) addDep(dep *service) {
	if s.isCircularDependency(dep) {
		_ = s.fail(fmt.Errorf("circular dependency detected: %v -> %v", s.id, dep.id))
		return
	}

	// detect circular dependencies
	s.deps[dep.id] = dep
	dep.depsOf[s.id] = s
	if dep.err != nil {
		if s.err == nil {
			_ = s.fail(nil)
		}
		s.err.addDepsErr(dep.err)
	}
}

func (s *service) isCircularDependency(dep *service) bool {
	if dep.id == s.id {
		return true
	}
	for _, d := range dep.deps {
		if s.isCircularDependency(d) {
			return true
		}
	}
	return false
}

func (s *service) getLogger() *zap.Logger {
	logger := s.app.getLogger("").With(zap.String("service.id", s.id))
	if s.name != "" {
		logger = logger.With(zap.String("service.name", s.name))
	}
	return logger
}

func (s *service) start(ctx context.Context) *ServiceError {
	s.startOnce.Do(func() {
		if s.err != nil {
			return
		}

		s.setStatusWithLock(Starting)

		// Start dependencies
		startErr := &ServiceError{
			svc: s,
		}
		wg := sync.WaitGroup{}
		wg.Add(len(s.deps))
		for _, dep := range s.deps {
			go func(dep *service) {
				defer wg.Done()
				if err := dep.start(ctx); err != nil {
					startErr.addDepsErr(err)
				}
			}(dep)
		}
		wg.Wait()

		if len(startErr.depsErrs) > 0 {
			s.failWithLock(startErr)
			return
		}

		// If all dependencies started successfully then start the service
		if s.err == nil {
			if start, ok := s.value.(svc.Runnable); ok {
				logger := s.getLogger()
				logger.Info("Service starting...")
				err := start.Start(s.context(ctx))
				if err != nil {
					s.failWithLock(err)
					logger.Error("Service failed to start", zap.Error(err))
					return
				}
				logger.Info("Service started successfully")
			}
		}

		s.registerMetric()
		s.setStatusWithLock(Running)
	})

	return s.err
}

func (s *service) stop(ctx context.Context) *ServiceError {
	if s.err != nil {
		return s.err
	}

	// if one of the dependencies is not running then don't stop
	for _, dep := range s.depsOf {
		if dep.Status() <= Stopping {
			<-s.stopChan
			return s.err
		}
	}

	s.stopOnce.Do(func() {
		if s.err != nil {
			return
		}

		s.setStatusWithLock(Stopping)
		defer func() {
			close(s.stopChan)
		}()

		if stop, ok := s.value.(svc.Runnable); ok {
			logger := s.getLogger()
			logger.Info("Service stopping...")
			err := stop.Stop(s.context(ctx))
			if err != nil {
				s.failWithLock(err)
				logger.Error("Service failed to stop", zap.Error(err))
				return
			}
			logger.Info("Service successfully stopped")
		}
		if s.err == nil {
			s.setStatusWithLock(Stopped)
		}

		stopErr := &ServiceError{
			svc: s,
		}
		wg := sync.WaitGroup{}
		wg.Add(len(s.deps))
		for _, dep := range s.deps {
			go func(dep *service) {
				defer wg.Done()
				if err := dep.stop(ctx); err != nil {
					stopErr.addDepsErr(err)
				}
			}(dep)
		}
		wg.Wait()

		if len(stopErr.depsErrs) > 0 {
			s.failWithLock(stopErr)
		}
	})

	return s.err
}

func (s *service) registerReadyCheck() error {
	if s.healthConfig.Check == nil {
		return nil
	}

	if s.healthConfig.Name == "" {
		if s.name != "" {
			s.healthConfig.Name = s.name
		} else {
			s.healthConfig.Name = s.id
		}
	}

	return s.app.readyHealth.Register(*s.healthConfig)
}

func (s *service) wrapCheck(check health.CheckFunc) health.CheckFunc {
	return func(ctx context.Context) error {
		// we lock to make sure that the service is not
		// stopped while we are checking if it is ready
		s.mux.RLock()
		defer s.mux.RUnlock()

		switch s.Status() {
		case Constructing, Constructed:
			return fmt.Errorf("service not started")
		case Starting:
			return fmt.Errorf("service starting")
		case Running:
			return check(ctx)
		case Stopping:
			return fmt.Errorf("service stopping")
		case Stopped:
			return fmt.Errorf("service stopped")
		case Error:
			return fmt.Errorf("service in error state: %v", s.err)
		}
		return nil
	}
}

func (s *service) setMetrics() {
	if m, ok := s.value.(svc.Metricable); ok {
		subsystem := s.name
		if subsystem == "" {
			subsystem = s.id
		}
		m.SetMetrics(sanitizeMetricName(s.app.name), sanitizeMetricName(subsystem), s.tags...)
	}
}

func (s *service) registerMetric() {
	if collector, ok := s.value.(prometheus.Collector); ok {
		if s.metricsPrefix != "" {
			prometheus.WrapRegistererWithPrefix(s.metricsPrefix, s.app.prometheus).MustRegister(collector)
		} else {
			s.app.prometheus.MustRegister(collector)
		}
	}
}

type ServiceError struct {
	svc *service

	mu        sync.RWMutex
	directErr error
	depsErrs  []*ServiceError
}

func (e *ServiceError) Error() string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var s string

	if e.directErr != nil {
		s = fmt.Sprintf("service %q: %v", e.svc.id, e.directErr)
	} else {
		s = fmt.Sprintf("service %q", e.svc.id)
	}

	if len(e.depsErrs) > 0 {
		for _, dep := range e.depsErrs {
			s += "\n"
			err := dep.Error()
			lines := strings.Split(err, "\n")
			indentedLines := make([]string, len(lines))
			for i, line := range lines {
				indentedLines[i] = fmt.Sprintf(">%s", line)
			}
			s += strings.Join(indentedLines, "\n")
		}
	}

	return s
}

func (e *ServiceError) addDepsErr(err *ServiceError) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.depsErrs = append(e.depsErrs, err)
}
