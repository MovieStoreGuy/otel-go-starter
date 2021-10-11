package otelstarter_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"

	launcher "github.com/MovieStoreGuy/otel-go-starter"
	"github.com/MovieStoreGuy/otel-go-starter/config"
)

type OtelTestHandler struct {
	T testing.TB
}

var _ otel.ErrorHandler = (*OtelTestHandler)(nil)

func (ot *OtelTestHandler) Handle(err error) {
	assert.NoError(ot.T, err, "Must not error throughout testing")
}

func TestLauncherUsingDefault(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	otel.SetErrorHandler(&OtelTestHandler{t})

	launcher.Start(ctx).Shutdown()
}

func TestLauncherWithConfiguredPipelines(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	assert.NotPanics(t, func() {
		launcher.Start(ctx,
			config.WithOtelErrorHandler(&OtelTestHandler{t}),
			config.WithTracesPipeline(
				config.WithTracingExporterOptions(
					config.WithExporterNamed("stdout"),
				),
			),
			config.WithMetricsPipeline(
				config.WithMetricsCollectionPeriod(time.Minute),
				config.WithMetricsExporterOptions(
					config.WithExporterNamed("stdout"),
				),
			),
		).Shutdown()
	})
}

func TestLauncherPanicsWithInvalidTracingConfig(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	assert.Panics(t, func() {
		launcher.Start(ctx, config.WithResource(ctx, nil))
	})

	assert.Panics(t, func() {
		launcher.Start(ctx, config.WithTracesPipeline(
			config.WithTracingExporterOptions(
				config.WithExporterNamed("unsupported-exporter"),
			),
		))
	})

	assert.Panics(t, func() {
		launcher.Start(ctx, config.WithTracesPipeline(
			config.WithTracingExporterOptions(
				config.WithExporterNamed("stdout"),
			),
			config.WithTracingPropagators("excellent-propagator"),
		))
	})
}
