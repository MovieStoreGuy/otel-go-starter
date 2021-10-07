package launcher

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/multierr"

	"github.com/MovieStoreGuy/otel-go-starter/config"
	"github.com/MovieStoreGuy/otel-go-starter/internal/pipeline/trace"
)

type Launcher interface {
	Shutdown()
}

type launch struct {
	conf              *config.Config
	shutdownCallbacks []func() error
}

// Start configures the global context of the open telemetry functionality
func Start(ctx context.Context, opts ...config.OptionFunc) Launcher {
	c := config.NewDefault()

	if err := c.Apply(opts...); err != nil {
		panic(err)
	}

	l := &launch{
		conf: &c,
	}

	otel.SetErrorHandler(c.GetErrorHandler())

	if c.Metrics.Enable {
		_ = 0
	}

	if c.Tracing.Enable {
		exporter, err := trace.NewExporterFactory().NewExporter(ctx, &c.Tracing)
		if err != nil {
			panic(err)
		}
		l.shutdownCallbacks = append(l.shutdownCallbacks, func() error {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			return exporter.Shutdown(ctx)
		})

		tp := sdktrace.NewTracerProvider(
			// TODO(Sean Marciniak): Add in a means of sampling here
			sdktrace.WithSpanProcessor(
				sdktrace.NewBatchSpanProcessor(exporter),
			),
			sdktrace.WithResource(c.GetResource()),
		)

		l.shutdownCallbacks = append(l.shutdownCallbacks, func() error {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			return tp.Shutdown(ctx)
		})

		prop, err := trace.NewPropagator(c.Tracing.Propergators)

		if err != nil {
			panic(err)
		}

		otel.SetTextMapPropagator(prop)
		otel.SetTracerProvider(tp)
	}

	return l
}

func (l *launch) Shutdown() {
	var err error
	for _, shutdown := range l.shutdownCallbacks {
		err = multierr.Append(err, shutdown())
	}
	if err != nil {
		otel.Handle(err)
	}
}
