# App Package

## Overview

The `app` package provides a dependency injection framework for building production-ready Go applications with built-in support for service lifecycle management, health checks, metrics, and observability. It follows the application container pattern, allowing you to declare services and their dependencies in a type-safe manner, while the framework handles initialization ordering, startup, shutdown, and error propagation.

### Rationale

Building production-ready applications requires handling cross-cutting concerns like:
- **Dependency Management**: Ensuring services are initialized in the correct order
- **Lifecycle Management**: Properly starting and stopping services with their dependencies
- **Health Checks**: Exposing readiness and liveness probes for orchestration platforms
- **Metrics**: Automatic instrumentation with Prometheus
- **Observability**: Structured logging with contextual information
- **HTTP Servers**: Configurable main API and health check endpoints
- **Graceful Shutdown**: Handling OS signals and cascading shutdown

The `app` package solves these problems by providing a container that automatically manages these concerns, allowing developers to focus on business logic.

## Key Features

1. **Dependency Injection**: Type-safe service registration with automatic dependency resolution
2. **Lifecycle Management**: Automatic service startup/shutdown with dependency ordering
3. **Health Checks**: Built-in liveness and readiness probes compatible with Kubernetes
4. **Prometheus Metrics**: Automatic metric collection and exposition
5. **HTTP Entrypoints**: Configurable main API and health check servers
6. **Structured Logging**: Integrated zap logger with contextual tags
7. **Graceful Shutdown**: OS signal handling (SIGINT, SIGTERM) with coordinated service shutdown
8. **Service Interfaces**: Optional interfaces for common service patterns (Runnable, Checkable, API, etc.)
9. **Error Aggregation**: Hierarchical error reporting across service dependencies
10. **Middleware Support**: Composable HTTP middleware chains

## Architecture

### The App Object

The `App` is the central container that manages your application's services. It:

- **Maintains a service registry**: Each service is identified by a unique ID and stored in a dependency graph
- **Handles construction**: Services are lazily constructed when first requested via `Provide()`
- **Manages lifecycle**: The `Start()` method walks the dependency graph and starts services in the correct order (dependencies first), while `Stop()` shuts them down in reverse order
- **Provides HTTP endpoints**: Optional main and healthz HTTP servers for APIs and observability
- **Exposes metrics**: Integrated Prometheus registry for application metrics
- **Propagates context**: Enriches contexts with logging and tagging information

### Service Lifecycle States

Services transition through the following states:

```
Constructing -> Constructed -> Starting -> Running -> Stopping -> Stopped
                                      \                    /
                                       \--> Error <-------/
```

- **Constructing**: Service constructor is being called
- **Constructed**: Service successfully created and registered
- **Starting**: Service and its dependencies are starting
- **Running**: Service is fully operational
- **Stopping**: Service is shutting down
- **Stopped**: Service successfully stopped
- **Error**: Service encountered an error at any stage

## Service Interfaces

Services can implement optional interfaces in the `app/svc` package to hook into application lifecycle:

### svc.Runnable
Services with long-running tasks (servers, workers, etc.)
```go
type Runnable interface {
    Start(context.Context) error
    Stop(context.Context) error
}
```

### svc.Checkable
Services that can report their health status
```go
type Checkable interface {
    Ready(ctx context.Context) error
}
```

### svc.API
Services that expose HTTP routes on the main server
```go
type API interface {
    RegisterHandler(mux *mux.Router)
}
```

### svc.Healthz
Services that expose routes on the health check server
```go
type Healthz interface {
    RegisterHealthzHandler(mux *mux.Router)
}
```

### svc.Middleware
Services that provide HTTP middleware
```go
type Middleware interface {
    RegisterMiddleware(chain alice.Chain) alice.Chain
}
```

### svc.Metricable
Services that can configure their metrics namespace
```go
type Metricable interface {
    SetMetrics(system, subsystem string, tags ...*tag.Tag)
}
```

### svc.MetricsCollector
Services that expose Prometheus metrics (same as `prometheus.Collector`)
```go
type MetricsCollector interface {
    Describe(ch chan<- *prometheus.Desc)
    Collect(ch chan<- prometheus.Metric)
}
```

### svc.Taggable
Services that accept contextual tags
```go
type Taggable interface {
    WithTags(tags ...*tag.Tag)
}
```

## Usage Examples

### Basic Application Setup

```go
package main

import (
    "context"
    "github.com/nmvalera/go-utils/app"
    "go.uber.org/zap"
)

func main() {
    // Create app configuration
    cfg := app.DefaultConfig()
    
    // Create application with options
    application, err := app.NewApp(
        cfg,
        app.WithName("my-service"),
        app.WithVersion("1.0.0"),
        app.WithReplaceGlobalLoggers(true),
    )
    if err != nil {
        panic(err)
    }
    
    // Register services
    // ... (see examples below)
    
    // Run application (blocks until SIGINT/SIGTERM)
    if err := application.Run(context.Background()); err != nil {
        panic(err)
    }
}
```

### Service Registration with Dependencies

```go
// Register a database service
db := app.Provide(application, "database", func() (*sql.DB, error) {
    return sql.Open("postgres", "postgresql://localhost/mydb")
})

// Register a repository that depends on the database
repo := app.Provide(application, "user-repo", func() (*UserRepository, error) {
    // The database service is already available here
    return NewUserRepository(db), nil
})

// Register a service that depends on the repository
svc := app.Provide(application, "user-service", func() (*UserService, error) {
    return NewUserService(repo), nil
})
```

### Implementing a Runnable Service

```go
type Worker struct {
    cancel context.CancelFunc
    done   chan struct{}
}

func (w *Worker) Start(ctx context.Context) error {
    ctx, w.cancel = context.WithCancel(ctx)
    w.done = make(chan struct{})
    
    go func() {
        defer close(w.done)
        ticker := time.NewTicker(10 * time.Second)
        defer ticker.Stop()
        
        for {
            select {
            case <-ctx.Done():
                return
            case <-ticker.C:
                // Do work...
            }
        }
    }()
    
    return nil
}

func (w *Worker) Stop(ctx context.Context) error {
    w.cancel()
    select {
    case <-w.done:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

// Register the worker
worker := app.Provide(application, "worker", func() (*Worker, error) {
    return &Worker{}, nil
})
```

### Implementing a Checkable Service

```go
type Database struct {
    db *sql.DB
}

func (d *Database) Ready(ctx context.Context) error {
    return d.db.PingContext(ctx)
}

// Register with custom health check configuration
db := app.Provide(
    application, 
    "database", 
    func() (*Database, error) {
        sqlDB, err := sql.Open("postgres", "...")
        if err != nil {
            return nil, err
        }
        return &Database{db: sqlDB}, nil
    },
    app.WithHealthConfig(&health.Config{
        Name:    "postgres",
        Timeout: 5 * time.Second,
    }),
)
```

### Implementing an API Service

```go
type UserAPI struct {
    service *UserService
}

func (api *UserAPI) RegisterHandler(router *mux.Router) {
    router.HandleFunc("/users/{id}", api.getUser).Methods(http.MethodGet)
    router.HandleFunc("/users", api.createUser).Methods(http.MethodPost)
}

func (api *UserAPI) getUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    // Handle request...
}

func (api *UserAPI) createUser(w http.ResponseWriter, r *http.Request) {
    // Handle request...
}

// Register the API and enable main entrypoint
api := app.Provide(application, "user-api", func() (*UserAPI, error) {
    application.EnableMainEntrypoint()  // Start HTTP server
    return &UserAPI{service: userService}, nil
})
```

### Implementing a Middleware Service

```go
type AuthMiddleware struct {
    secret string
}

func (m *AuthMiddleware) RegisterMiddleware(chain alice.Chain) alice.Chain {
    return chain.Append(m.authenticate)
}

func (m *AuthMiddleware) authenticate(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if !m.validateToken(token) {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}

func (m *AuthMiddleware) validateToken(token string) bool {
    // Validate token...
    return true
}

// Register middleware
auth := app.Provide(application, "auth", func() (*AuthMiddleware, error) {
    return &AuthMiddleware{secret: "my-secret"}, nil
})
```

### Implementing a Metrics Service

```go
type MetricsService struct {
    requestCount prometheus.Counter
    requestDuration prometheus.Histogram
}

func (m *MetricsService) SetMetrics(system, subsystem string, tags ...*tag.Tag) {
    m.requestCount = prometheus.NewCounter(prometheus.CounterOpts{
        Namespace: system,
        Subsystem: subsystem,
        Name:      "requests_total",
        Help:      "Total number of requests",
    })
    
    m.requestDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
        Namespace: system,
        Subsystem: subsystem,
        Name:      "request_duration_seconds",
        Help:      "Request duration in seconds",
    })
}

func (m *MetricsService) Describe(ch chan<- *prometheus.Desc) {
    m.requestCount.Describe(ch)
    m.requestDuration.Describe(ch)
}

func (m *MetricsService) Collect(ch chan<- prometheus.Metric) {
    m.requestCount.Collect(ch)
    m.requestDuration.Collect(ch)
}

// Register service with custom component name
metrics := app.Provide(
    application, 
    "metrics", 
    func() (*MetricsService, error) {
        return &MetricsService{}, nil
    },
    app.WithComponentName("api"),  // Metrics will be under "myapp_api_*"
)
```

### Service with Tags

```go
type Service struct {
    tags tag.Set
}

func (s *Service) WithTags(tags ...*tag.Tag) {
    s.tags = s.tags.WithTags(tags...)
}

// Register with custom tags
svc := app.Provide(
    application,
    "my-service",
    func() (*Service, error) {
        return &Service{}, nil
    },
    app.WithTags(
        tag.Key("env").String("production"),
        tag.Key("region").String("us-east-1"),
    ),
)
```

### Enabling HTTP Entrypoints

```go
// Enable main API server (for business logic endpoints)
app.Provide(application, "top", func() (any, error) {
    application.EnableMainEntrypoint()
    application.EnableHealthzEntrypoint()
    return nil, nil
})

// Now the application will have:
// - Main server on :8080 (configurable)
// - Health server on :8081 (configurable) with:
//   - GET /live  (liveness probe)
//   - GET /ready (readiness probe, checks all Checkable services)
//   - GET /metrics (Prometheus metrics)
```

### Complete Example

```go
package main

import (
    "context"
    "database/sql"
    "net/http"
    "time"
    
    "github.com/gorilla/mux"
    "github.com/nmvalera/go-utils/app"
    "github.com/nmvalera/go-utils/app/svc"
    "go.uber.org/zap"
)

type Database struct {
    db *sql.DB
}

func (d *Database) Start(ctx context.Context) error {
    return d.db.PingContext(ctx)
}

func (d *Database) Stop(ctx context.Context) error {
    return d.db.Close()
}

func (d *Database) Ready(ctx context.Context) error {
    return d.db.PingContext(ctx)
}

type UserAPI struct {
    db     *Database
}

func (api *UserAPI) RegisterHandler(router *mux.Router) {
    router.HandleFunc("/users/{id}", api.getUser).Methods(http.MethodGet)
}

func (api *UserAPI) getUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"id": "123", "name": "John Doe"}`))
}

func main() {
    cfg := app.DefaultConfig()
    
    application, err := app.NewApp(
        cfg,
        app.WithName("user-service"),
        app.WithVersion("1.0.0"),
        app.WithReplaceGlobalLoggers(true),
    )
    if err != nil {
        panic(err)
    }
    
    // Register database
    db := app.Provide(
        application,
        "database",
        func() (*Database, error) {
            sqlDB, err := sql.Open("postgres", "postgresql://localhost/users")
            if err != nil {
                return nil, err
            }
            return &Database{db: sqlDB}, nil
        },
        app.WithComponentName("postgres"),
    )
    
    // Register API
    app.Provide(application, "api", func() (*UserAPI, error) {
        application.EnableMainEntrypoint()
        application.EnableHealthzEntrypoint()
        
        logger := zap.L().With(zap.String("component", "api"))
        return &UserAPI{
            db:     db,
            logger: logger,
        }, nil
    })
    
    // Run application
    if err := application.Run(context.Background()); err != nil {
        zap.L().Fatal("Application failed", zap.Error(err))
    }
}
```

## Configuration

### App Configuration

The `Config` struct provides configuration for the application:

```go
type Config struct {
    MainEntrypoint    *http.EntrypointConfig  // Main API server config
    HealthzEntrypoint *http.EntrypointConfig  // Health check server config
    HealthzServer     *HealthzServerConfig    // Health endpoint paths
    Log               *log.Config             // Logging configuration
    StartTimeout      *string                 // Startup timeout (e.g., "30s")
    StopTimeout       *string                 // Shutdown timeout (e.g., "30s")
}
```

Default configuration:
```go
cfg := app.DefaultConfig()
// Main server: :8080
// Health server: :8081
// Liveness: /live
// Readiness: /ready
// Metrics: /metrics
```

### Configuration Flags

The application supports the following configuration flags:

#### Main Server Configuration

| Flag | Environment Variable | Description | Default Value |
|------|---------------------|-------------|---------------|
| `--main-ep-addr` | `MAIN_ENTRYPOINT_ADDR` | Main entrypoint address to listen on | `:8080` |
| `--main-ep-net-keep-alive` | `MAIN_ENTRYPOINT_NET_KEEP_ALIVE` | Main entrypoint keep alive | `0` |
| `--main-ep-http-read-timeout` | `MAIN_ENTRYPOINT_HTTP_READ_TIMEOUT` | Main entrypoint maximum duration for reading an entire request including the body (zero means no timeout) | `30s` |
| `--main-ep-http-read-header-timeout` | `MAIN_READ_HEADER_TIMEOUT` | Main entrypoint maximum duration for reading request headers (zero uses the value of read timeout) | `30s` |
| `--main-ep-http-write-timeout` | `MAIN_ENTRYPOINT_HTTP_WRITE_TIMEOUT` | Main entrypoint maximum duration for writing the response (zero means no timeout) | `30s` |
| `--main-ep-http-idle-timeout` | `MAIN_ENTRYPOINT_HTTP_IDLE_TIMEOUT` | Main entrypoint maximum amount of time to wait for the next request when keep-alives are enabled (zero uses the value of read timeout) | `30s` |

#### Health Check Server Configuration

| Flag | Environment Variable | Description | Default Value |
|------|---------------------|-------------|---------------|
| `--healthz-ep-addr` | `HEALTHZ_ENTRYPOINT_ADDR` | Health check entrypoint address | `:8081` |
| `--healthz-ep-keep-alive` | `HEALTHZ_ENTRYPOINT_NET_KEEP_ALIVE` | Health check entrypoint keep alive | `0` |
| `--healthz-ep-http-read-timeout` | `HEALTHZ_NTRYPOINT_HTTP_READ_TIMEOUT` | Health entrypoint maximum duration for reading an entire request including the body (zero means no timeout) | `30s` |
| `--healthz-ep-http-read-header-timeout` | `HEALTHZ_NTRYPOINT_HTTP_READ_HEADER_TIMEOUT` | Health entrypoint maximum duration for reading request headers (zero uses the value of read timeout) | `30s` |
| `--healthz-ep-http-write-timeout` | `HEALTHZ_NTRYPOINT_HTTP_WRITE_TIMEOUT` | Health entrypoint maximum duration for writing the response (zero means no timeout) | `30s` |
| `--healthz-ep-http-idle-timeout` | `HEALTHZ_NTRYPOINT_HTTP_IDLE_TIMEOUT` | Health entrypoint maximum amount of time to wait for the next request when keep-alives are enabled (zero uses the value of read timeout) | `30s` |

All configuration flags can be set via command-line arguments, environment variables, or through a configuration file. The application uses [Viper](https://github.com/spf13/viper) for configuration management, which supports multiple configuration formats.

## Underlying Dependencies

The `app` package integrates several production-ready libraries:

### Core Dependencies

- **[uber-go/zap](https://github.com/uber-go/zap)**: High-performance structured logging
  - *Rationale*: Zero-allocation, type-safe logging with excellent performance characteristics

- **[gorilla/mux](https://github.com/gorilla/mux)**: Powerful HTTP router and URL matcher
  - *Rationale*: Feature-rich, flexible routing with path parameters, regular expressions, and middleware support

- **[justinas/alice](https://github.com/justinas/alice)**: Middleware chaining
  - *Rationale*: Clean API for composing HTTP middleware in a type-safe manner

- **[prometheus/client_golang](https://github.com/prometheus/client_golang)**: Prometheus metrics
  - *Rationale*: Industry-standard metrics collection and exposition for observability

- **[hellofresh/health-go](https://github.com/hellofresh/health-go)**: Health check library
  - *Rationale*: Kubernetes-compatible health check endpoints with multiple check support

### Supporting Dependencies

- **[spf13/viper](https://github.com/spf13/viper)**: Configuration management
  - *Rationale*: Supports multiple configuration sources (files, env vars, flags) with a unified API

- **[spf13/pflag](https://github.com/spf13/pflag)**: POSIX/GNU-style flags
  - *Rationale*: Drop-in replacement for Go's flag package with better POSIX compliance

- **[MadAppGang/httplog](https://github.com/MadAppGang/httplog)**: HTTP request logging
  - *Rationale*: Structured HTTP access logs with zap integration

## Best Practices

### 1. Service Organization
- Use meaningful service IDs that describe the service's purpose
- Group related functionality into single services
- Keep services focused on a single responsibility

### 2. Dependency Management
- Declare dependencies in the constructor by calling `app.Provide()`
- Avoid circular dependencies (the framework will detect and report them)
- Start external resources (servers, connections) in `Start()`, not in constructors

### 3. Error Handling
- Return errors from constructors for initialization failures
- Return errors from `Start()` for startup failures
- Return errors from `Stop()` for cleanup failures
- Check `app.Error()` after registration to catch construction errors early

### 4. Health Checks
- Implement `Checkable` for services
- Keep health checks lightweight (they run frequently)
- Use appropriate timeouts in health check configuration

### 5. Logging
- Always use contextual logging (e.g. `log.LoggerFromContext(ctx).Info(...)`) to enable passing fields across log messages
- Use structured logging with fields rather than string formatting
- Use base zap logging (e.g `logger.Info(..)`) and not sugared logging (e.g. `logger.Sugar().Infow(..)`) which brings a drop of performance
- Add component context via tags or logger fields
- Log at appropriate levels (Info for lifecycle, Debug for operations, Error for failures)

### 6. Metrics
- Implement `Metricable` to customize metric namespace
- Implement `prometheus.Collector` to expose custom metrics
- Use consistent naming conventions (counters end in `_total`, histograms in `_seconds`, etc.)

### 7. Testing
- Create test apps with `NewApp()` and `zap.NewNop()` for silent logging
- Use `app.Start()` and `app.Stop()` for controlled lifecycle testing
- Test service error scenarios by returning errors from constructors/Start/Stop

## Error Handling

The package provides hierarchical error reporting via `ServiceError`. When a service fails, the error includes:
- The service ID and direct error
- Errors from all failed dependencies (recursively)
- Indentation showing the dependency hierarchy

Example error output:
```
service "main-api": connection refused
>service "database": connection timeout
>>service "database-config": file not found
```

## Thread Safety

- `App.Provide()` is NOT thread-safe and should only be called during initialization
- `App.Start()` and `App.Stop()` can be called concurrently (idempotent via `sync.Once`)
- Service lifecycle methods are protected by internal locks
- Health checks are thread-safe and can be called concurrently
