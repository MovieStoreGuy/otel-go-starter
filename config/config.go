package config

import (
	"errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.uber.org/multierr"
)

var (
	ErrNilParamProvided = errors.New("nil value provided")
	ErrInvalidParam     = errors.New("invalid value provided")
)

type Config struct {
	Metrics Pipeline
	Tracing Pipeline

	errHandler otel.ErrorHandler
	resource   *resource.Resource
}

type Pipeline struct {
	Enable         bool
	AllowInsecure  bool
	UseCompression bool
	Exporter       string
	Endpoint       string
	Headers        map[string]string
	Propergators   []string
}

type OptionFunc func(*Config) error

type PipeOption func(*Pipeline) error

func NewDefault() Config {
	return Config{
		Metrics: Pipeline{
			Enable: false,
		},
		Tracing: Pipeline{
			Enable: false,
			Propergators: []string{
				"baggage",
				"tracecontext",
			},
		},
		errHandler: otel.GetErrorHandler(),
		resource:   resource.Default(),
	}
}

func (c *Config) Apply(opts ...OptionFunc) (err error) {
	for _, opt := range opts {
		err = multierr.Append(err, opt(c))
	}
	return err
}

func (c *Config) GetErrorHandler() otel.ErrorHandler {
	return c.errHandler
}

func (c *Config) GetResource() *resource.Resource {
	return c.resource
}
