package app

import (
	"github.com/hellofresh/health-go/v5"
	"github.com/nmvalera/go-utils/tag"
	"go.uber.org/zap"
)

type Option func(*App) error

// WithAppName sets the name of the application.
func WithName(name string) Option {
	return func(a *App) error {
		a.name = name
		return nil
	}
}

// WithVersion sets the version of the application.
func WithVersion(version string) Option {
	return func(a *App) error {
		a.version = version
		return nil
	}
}

// WithLogger sets the logger of the application.
func WithLogger(logger *zap.Logger) Option {
	return func(a *App) error {
		a.logger = logger
		return nil
	}
}

func WithReplaceGlobalLoggers(replace bool) Option {
	return func(a *App) error {
		a.replaceGlobalLoggers = replace
		return nil
	}
}

type ServiceOption func(*service) error

// WithHealthConfig sets the health config of the service.
func WithHealthConfig(cfg *health.Config) ServiceOption {
	return func(s *service) error {
		if cfg.Name != "" {
			s.healthConfig.Name = cfg.Name
		}

		if cfg.Check != nil {
			s.healthConfig.Check = s.wrapCheck(cfg.Check)
		}

		if cfg.Timeout != 0 {
			s.healthConfig.Timeout = cfg.Timeout
		}

		if cfg.SkipOnErr {
			s.healthConfig.SkipOnErr = true
		}

		return nil
	}
}

// WithTags sets the tags of the service.
func WithTags(tags ...*tag.Tag) ServiceOption {
	return func(s *service) error {
		s.tags = s.tags.WithTags(tags...)
		return nil
	}
}

// WithComponentName sets the name of the component.
// Multiple services can have the same component name.
// By default, the component name is the name of the service identifier which is unique.
// This enables to override the default name (without unicity constraints)
func WithComponentName(name string) ServiceOption {
	return func(s *service) error {
		s.name = name
		return nil
	}
}

// WithComponentNameChained sets the chained name of the component.
// If chained is true, the component name will be chained to the existing component name.
// If chained is false, the component name will replace the existing component name.
//
// Example:
// WithComponentNameChained(true) will return a new context with the component tag "component1.component2"
// WithComponentNameChained(false) will return a new context with the component tag "component2"
func WithComponentNameChained(chained bool) ServiceOption {
	return func(s *service) error {
		s.chainedName = chained
		return nil
	}
}
