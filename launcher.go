package otelstarter

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	metricglobal "go.opentelemetry.io/otel/metric/global"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/multierr"

	"github.com/MovieStoreGuy/otel-go-starter/config"
	"github.com/MovieStoreGuy/otel-go-starter/internal/pipeline/metric"
	"github.com/MovieStoreGuy/otel-go-starter/internal/pipeline/trace"
)

// Launcher stores all the information used from configuring the global
// Open Telemetry properties and allows for graceful shutdowns
type Launcher interface {
	Shutdown()
}

type launch struct {
	shutdownCallbacks []func() error
}

// Start configures the global context of the open telemetry functionality
// Any issues trying to configure any of the instrumentation will cause the method to panic
func Start(ctx context.Context, opts ...config.OptionFunc) Launcher {
	c := config.NewDefault()

	if err := c.Apply(opts...); err != nil {
		panic(err)
	}

	l := &launch{}

	otel.SetErrorHandler(c.GetErrorHandler())

	if c.Metrics.Enable {
		exporter, err := metric.NewExporterFactory().NewExporter(ctx, &c.Metrics.Export)
		if err != nil {
			panic(err)
		}
		if sh, ok := exporter.(metric.ShutdownExporter); ok {
			l.shutdownCallbacks = append(l.shutdownCallbacks, gracefulShutdown(sh.Shutdown))
		}

		pusher := controller.New(
			processor.NewFactory(
				selector.NewWithInexpensiveDistribution(),
				exporter,
			),
			controller.WithExporter(exporter),
			controller.WithResource(c.GetResource()),
			controller.WithCollectPeriod(time.Second),
		)

		if err := pusher.Start(ctx); err != nil {
			panic(err)
		}

		l.shutdownCallbacks = append(l.shutdownCallbacks, gracefulShutdown(pusher.Stop))
		metricglobal.SetMeterProvider(pusher)
	}

	if c.Tracing.Enable {
		exporter, err := trace.NewExporterFactory().NewExporter(ctx, &c.Tracing.Export)
		if err != nil {
			panic(err)
		}
		l.shutdownCallbacks = append(l.shutdownCallbacks, gracefulShutdown(exporter.Shutdown))

		var sampler sdktrace.Sampler
		if c.Tracing.Sample {
			sampler = sdktrace.AlwaysSample()
		}

		tp := sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sampler),
			sdktrace.WithSpanProcessor(
				sdktrace.NewBatchSpanProcessor(exporter),
			),
			sdktrace.WithResource(c.GetResource()),
		)

		l.shutdownCallbacks = append(l.shutdownCallbacks, gracefulShutdown(tp.Shutdown))

		prop, err := trace.NewPropagators(c.Tracing.Propagators)

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

func gracefulShutdown(f func(ctx context.Context) error) func() error {
	return func() error {
		ctx, done := context.WithTimeout(context.Background(), time.Second)
		defer done()

		return f(ctx)
	}
}
