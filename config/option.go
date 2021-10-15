package config

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/multierr"
)

func WithServiceName(name string) OptionFunc {
	return func(c *Config) error {
		r, err := resource.Merge(c.GetResource(), resource.NewSchemaless(semconv.ServiceNameKey.String(name)))
		if err != nil {
			return err
		}
		c.resource = r
		return nil
	}
}

func WithResourceDetector(ctx context.Context, detector resource.Detector) OptionFunc {
	return func(c *Config) error {
		if detector == nil {
			return fmt.Errorf("resource detector is nil: %w", ErrNilParamProvided)
		}

		r, err := detector.Detect(ctx)
		if err != nil {
			return err
		}

		c.resource, err = resource.Merge(c.resource, r)
		return err
	}
}

func WithAttributes(attrs ...attribute.KeyValue) OptionFunc {
	return func(c *Config) error {
		r, err := resource.Merge(c.GetResource(), resource.NewSchemaless(attrs...))
		if err != nil {
			return err
		}
		c.resource = r
		return nil
	}
}

func WithOtelErrorHandler(handler otel.ErrorHandler) OptionFunc {
	return func(c *Config) error {
		if handler == nil {
			return fmt.Errorf("error handler is nil: %w", ErrNilParamProvided)
		}
		c.errHandler = handler
		return nil
	}
}

func WithMetricsPipeline(pipeOpts ...MetricsOption) OptionFunc {
	return func(c *Config) (err error) {
		c.Metrics.Enable = true
		for _, opt := range pipeOpts {
			if opt == nil {
				return fmt.Errorf("nil metric pipeline option provided: %w", ErrNilParamProvided)
			}
			err = multierr.Append(err, opt(&c.Metrics))
		}
		return err
	}
}

func WithTracesPipeline(pipeOpts ...TracingOption) OptionFunc {
	return func(c *Config) (err error) {
		c.Tracing.Enable = true
		for _, opt := range pipeOpts {
			if opt == nil {
				return fmt.Errorf("nil traces pipeline option provided: %w", ErrNilParamProvided)
			}
			err = multierr.Append(err, opt(&c.Tracing))
		}
		return err
	}
}

func WithTracingExporterOptions(opts ...ExportOption) TracingOption {
	return func(t *Tracing) (err error) {
		for _, opt := range opts {
			if opt == nil {
				return fmt.Errorf("nil trace exporter option provided: %w", ErrNilParamProvided)
			}
			err = multierr.Append(err, opt(&t.Export))
		}
		return err
	}
}

func WithMetricsExporterOptions(opts ...ExportOption) MetricsOption {
	return func(m *Metrics) (err error) {
		for _, opt := range opts {
			if opt == nil {
				return fmt.Errorf("nil trace exporter option provided: %w", ErrNilParamProvided)
			}
			err = multierr.Append(err, opt(&m.Export))
		}
		return err
	}
}

func WithExporterInsecureConnection() ExportOption {
	return func(p *Export) error {
		p.AllowInsecure = true
		return nil
	}
}

func WithExporterUseCompression() ExportOption {
	return func(p *Export) error {
		p.UseCompression = true
		return nil
	}
}

// WithPipelineEndpoint will validate that the provided endpoint
// has a valid schema and that the hostname can be resolved
func WithExporterEndpoint(endpoint string) ExportOption {
	return func(p *Export) error {
		u, err := url.Parse(endpoint)
		if err != nil {
			return err
		}
		switch u.Scheme {
		case "http", "https":
			// expected schemas for the endpoint
		default:
			return fmt.Errorf("unknown scheme provided; must be http(s): %w", ErrInvalidParam)
		}

		if _, err := net.LookupHost(u.Hostname()); err != nil {
			return multierr.Append(err, ErrInvalidParam)
		}

		p.Endpoint = u.String()

		return nil
	}
}

func WithExporterHeaders(headers map[string]string) ExportOption {
	return func(p *Export) error {
		if headers == nil {
			return fmt.Errorf("header is nil: %w", ErrNilParamProvided)
		}

		if p.Headers == nil {
			p.Headers = make(map[string]string)
		}

		for k, v := range headers {
			if _, exist := p.Headers[k]; exist {
				return fmt.Errorf("conflict in head key %s: %w", k, ErrInvalidParam)
			}
			p.Headers[k] = v
		}

		return nil
	}
}

func WithExporterNamed(named string) ExportOption {
	return func(p *Export) error {
		if named == "" {
			return fmt.Errorf("no exporter named: %w", ErrInvalidParam)
		}
		p.Named = named
		return nil
	}
}

func WithTracingPropagators(use ...string) TracingOption {
	return func(p *Tracing) error {
		if len(use) == 0 {
			return fmt.Errorf("no pipeline propagators defined: %w", ErrNilParamProvided)
		}
		p.Propagators = use
		return nil
	}
}

func WithTracingSampled() TracingOption {
	return func(t *Tracing) error {
		t.Sample = true
		return nil
	}
}

func WithMetricsCollectionPeriod(t time.Duration) MetricsOption {
	return func(m *Metrics) error {
		if t < 0 {
			return fmt.Errorf("collection period must be positive value: %w", ErrInvalidParam)
		}
		m.CollectPeriod = t
		return nil
	}
}
