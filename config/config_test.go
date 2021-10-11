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
			config.WithMetricsExporterOptions(
				config.WithExporterInsecureConnection(),
				config.WithExporterUseCompression(),
				config.WithExporterEndpoint("http://localhost:9094"),
				config.WithExporterNamed("otlpgrpc"),
			),
		),
		config.WithTracesPipeline(
			config.WithTracingExporterOptions(
				config.WithExporterInsecureConnection(),
				config.WithExporterUseCompression(),
				config.WithExporterEndpoint("http://localhost:9094"),
				config.WithExporterNamed("otlpgrpc"),
				config.WithExporterHeaders(map[string]string{
					"Service-Domain": "pineapples",
				}),
			),
			config.WithTracingPropagators(
				"b3",
				"ot",
			),
			config.WithTracingSampled(),
		),
	), "Must not error when applying valid configuration")

	assert.True(t, conf.Metrics.Enable)
	assert.True(t, conf.Metrics.Export.AllowInsecure)
	assert.True(t, conf.Metrics.Export.UseCompression)
	assert.Equal(t, "http://localhost:9094", conf.Metrics.Export.Endpoint)
	assert.Equal(t, "otlpgrpc", conf.Metrics.Export.Named)

	assert.True(t, conf.Tracing.Enable)
	assert.True(t, conf.Tracing.Export.AllowInsecure)
	assert.True(t, conf.Tracing.Export.UseCompression)
	assert.True(t, conf.Tracing.Sample)
	assert.Equal(t, "http://localhost:9094", conf.Tracing.Export.Endpoint)
	assert.Equal(t, "otlpgrpc", conf.Tracing.Export.Named)
	assert.Equal(t, map[string]string{"Service-Domain": "pineapples"}, conf.Tracing.Export.Headers)
	assert.Equal(t, []string{"b3", "ot"}, conf.Tracing.Propagators)
}
