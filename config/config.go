package config

import (
	"errors"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.uber.org/multierr"
)

var (
	ErrNilParamProvided = errors.New("nil value provided")
	ErrInvalidParam     = errors.New("invalid value provided")
)

type Config struct {
	Metrics Metrics
	Tracing Tracing

	errHandler otel.ErrorHandler
	resource   *resource.Resource
}

type Export struct {
	AllowInsecure  bool
	UseCompression bool
	Named          string
	Endpoint       string
	Headers        map[string]string
}

type Tracing struct {
	Enable bool

	Export Export

	Sample      bool
	Propagators []string
}

type Metrics struct {
	Enable bool

	Export Export

	CollectPeriod time.Duration
}

// Method types to programatically validate additions
// to the existing config
type (
	OptionFunc    func(*Config) error
	ExportOption  func(*Export) error
	TracingOption func(*Tracing) error
	MetricsOption func(*Metrics) error
)

func NewDefault() *Config {
	return &Config{
		Metrics: Metrics{
			Enable: false,
			Export: Export{
				Named:   "otlpgrpc",
				Headers: map[string]string{},
			},
			CollectPeriod: time.Second,
		},
		Tracing: Tracing{
			Enable: false,
			Export: Export{
				Named:   "otlpgrpc",
				Headers: map[string]string{},
			},
			Propagators: []string{
				"baggage",
				"tracecontext",
			},
			Sample: false,
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
