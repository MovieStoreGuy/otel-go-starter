package config_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"

	"github.com/MovieStoreGuy/otel-go-starter/config"
)

type TestDetector struct {
	t *testing.T
}

var _ resource.Detector = (*TestDetector)(nil)

func (td TestDetector) Detect(_ context.Context) (*resource.Resource, error) {
	return resource.NewSchemaless(attribute.String("test.name", td.t.Name())), nil
}

func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	conf := config.NewDefault()

	assert.NotNil(t, conf.GetResource(), "Must not return a nil resource")
	assert.NotNil(t, conf.GetErrorHandler(), "Must not return a nil error handler")
}

func TestApplyingConfig(t *testing.T) {
	t.Parallel()

	conf := config.NewDefault()

	assert.NoError(t, conf.Apply(), "Must not error when applying valid configuration")
	assert.NoError(t, conf.Apply(
		config.WithResource(context.Background(), TestDetector{t}),
		config.WithMetricsPipeline(
			config.WithPipelineEnabled(),
			config.WithPipelineInsecureConnection(),
			config.WithPipelineCompression(),
			config.WithPipelineEndpoint("http://localhost:9094"),
			config.WithPipelineExporter("otlpgrpc"),
		),
		config.WithTracesPipeline(
			config.WithPipelineEnabled(),
			config.WithPipelineInsecureConnection(),
			config.WithPipelineCompression(),
			config.WithPipelineEndpoint("http://localhost:9094"),
			config.WithPipelineExporter("otlpgrpc"),
			config.WithPipelineHeaders(map[string]string{
				"Service-Domain": "pineapples",
			}),
			config.WithPipelinePropagators(
				"b3",
				"ot",
			),
		),
	), "Must not error when applying valid configuration")

	assert.True(t, conf.Metrics.Enable)
	assert.True(t, conf.Metrics.AllowInsecure)
	assert.True(t, conf.Metrics.UseCompression)
	assert.Equal(t, "http://localhost:9094", conf.Metrics.Endpoint)
	assert.Equal(t, "otlpgrpc", conf.Metrics.Exporter)

	assert.True(t, conf.Tracing.Enable)
	assert.True(t, conf.Tracing.AllowInsecure)
	assert.True(t, conf.Tracing.UseCompression)
	assert.Equal(t, "http://localhost:9094", conf.Tracing.Endpoint)
	assert.Equal(t, "otlpgrpc", conf.Tracing.Exporter)
	assert.Equal(t, map[string]string{"Service-Domain": "pineapples"}, conf.Tracing.Headers)
	assert.Equal(t, []string{"b3", "ot"}, conf.Tracing.Propagators)
}
