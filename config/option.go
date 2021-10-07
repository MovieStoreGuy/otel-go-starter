package config

import (
	"context"
	"fmt"
	"net"
	"net/url"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.uber.org/multierr"
)

func WithResource(ctx context.Context, detector resource.Detector) OptionFunc {
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
func WithOtelErrorHandler(handler otel.ErrorHandler) OptionFunc {
	return func(c *Config) error {
		if handler == nil {
			return fmt.Errorf("error handler is nil: %w", ErrNilParamProvided)
		}
		c.errHandler = handler
		return nil
	}
}

func WithMetricsPipeline(pipeOpts ...PipeOption) OptionFunc {
	return func(c *Config) (err error) {
		for _, opt := range pipeOpts {
			if opt == nil {
				return fmt.Errorf("nil metric pipeline option provided: %w", ErrNilParamProvided)
			}
			err = multierr.Append(err, opt(&c.Metrics))
		}
		return err
	}
}

func WithTracesPipeline(pipeOpts ...PipeOption) OptionFunc {
	return func(c *Config) (err error) {
		for _, opt := range pipeOpts {
			if opt == nil {
				return fmt.Errorf("nil traces pipeline option provided: %w", ErrNilParamProvided)
			}
			err = multierr.Append(err, opt(&c.Tracing))
		}
		return err
	}
}

func WithPipelineEnabled() PipeOption {
	return func(p *Pipeline) error {
		p.Enable = true
		return nil
	}
}

func WithPipelineInsecureConnection() PipeOption {
	return func(p *Pipeline) error {
		p.AllowInsecure = true
		return nil
	}
}

func WithPipelineCompression() PipeOption {
	return func(p *Pipeline) error {
		p.UseCompression = true
		return nil
	}
}

func WithPipelineSampled() PipeOption {
	return func(p *Pipeline) error {
		p.Sample = true
		return nil
	}
}

// WithPipelineEndpoint will validate that the provided endpoint
// has a valid schema and that the hostname can be resolved
func WithPipelineEndpoint(endpoint string) PipeOption {
	return func(p *Pipeline) error {
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

func WithPipelineHeaders(headers map[string]string) PipeOption {
	return func(p *Pipeline) error {
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

func WithPipelineExporter(named string) PipeOption {
	return func(p *Pipeline) error {
		if named == "" {
			return fmt.Errorf("no exporter named: %w", ErrInvalidParam)
		}
		p.Exporter = named
		return nil
	}
}

func WithPipelinePropagators(use ...string) PipeOption {
	return func(p *Pipeline) error {
		if len(use) == 0 {
			return fmt.Errorf("no pipeline propagators defined: %w", ErrNilParamProvided)
		}
		p.Propagators = use
		return nil
	}
}
