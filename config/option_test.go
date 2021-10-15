package config_test

import (
	"context"
	"testing"

	"github.com/MovieStoreGuy/otel-go-starter/config"
	"github.com/stretchr/testify/assert"
)

func TestValidConfigOptions(t *testing.T) {
	t.Parallel()
}

func TestNilParamConfigOptions(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		method string
		opt    config.OptionFunc
	}{
		{method: "WithResource", opt: config.WithResourceDetector(context.TODO(), nil)},
		{method: "WithOtelErrorHandler", opt: config.WithOtelErrorHandler(nil)},
		{method: "WithMetricsPipeline", opt: config.WithMetricsPipeline(nil)},
		{method: "WithTracesPipeline", opt: config.WithTracesPipeline(nil)},
		{method: "WithPipelineHeaders", opt: config.WithMetricsPipeline(
			config.WithMetricsExporterOptions(config.WithExporterHeaders(nil)),
		)},
		{method: "WithPipelinePropagators", opt: config.WithTracesPipeline(
			config.WithTracingPropagators(),
		)},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			assert.ErrorIs(t, config.NewDefault().Apply(tc.opt), config.ErrNilParamProvided)
		})
	}
}

func TestInvalidConfigOptions(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		method string
		opt    config.OptionFunc
	}{
		{method: "WithPipelineEndpoint.InvalidScheme", opt: config.WithMetricsPipeline(
			config.WithMetricsExporterOptions(config.WithExporterEndpoint("wss://localhost")),
		)},
		{method: "WithPipelineEndpoint.InvalidHost", opt: config.WithTracesPipeline(
			config.WithTracingExporterOptions(config.WithExporterEndpoint("http://pineapples.pizza")),
		)},
		{method: "WithPipelineHeaders.DuplicateEntries", opt: config.WithTracesPipeline(
			config.WithTracingExporterOptions(
				config.WithExporterHeaders(map[string]string{"Otel-Service": "foo"}),
				config.WithExporterHeaders(map[string]string{"Otel-Service": "bar"}),
			),
		)},
		{method: "WithPipelineExporter", opt: config.WithMetricsPipeline(config.WithMetricsExporterOptions(config.WithExporterNamed("")))},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			assert.ErrorIs(t, config.NewDefault().Apply(tc.opt), config.ErrInvalidParam)
		})
	}
}
