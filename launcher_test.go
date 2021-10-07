package launcher_test

import (
	"context"
	"testing"

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

func TestLauncherWithConfiguredTracePipeline(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	assert.NotPanics(t, func() {
		launcher.Start(ctx,
			config.WithOtelErrorHandler(&OtelTestHandler{t}),
			config.WithTracesPipeline(
				config.WithPipelineEnabled(),
				config.WithPipelineExporter("stdout"),
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
			config.WithPipelineEnabled(),
			config.WithPipelineExporter("unsupported-exporter"),
		))
	})

	assert.Panics(t, func() {
		launcher.Start(ctx, config.WithTracesPipeline(
			config.WithPipelineEnabled(),
			config.WithPipelineExporter("stdout"),
			config.WithPipelinePropagators("excellent-propagator"),
		))
	})
}
