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
		{method: "WithResource", opt: config.WithResource(context.TODO(), nil)},
		{method: "WithOtelErrorHandler", opt: config.WithOtelErrorHandler(nil)},
		{method: "WithMetricsPipeline", opt: config.WithMetricsPipeline(nil)},
		{method: "WithTracesPipeline", opt: config.WithTracesPipeline(nil)},
		{method: "WithPipelineHeaders", opt: config.WithMetricsPipeline(config.WithPipelineHeaders(nil))},
		{method: "WithPipelinePropagators", opt: config.WithMetricsPipeline(config.WithPipelinePropagators())},
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
		{method: "WithPipelineEndpoint.InvalidScheme", opt: config.WithMetricsPipeline(config.WithPipelineEndpoint("wss://localhost"))},
		{method: "WithPipelineEndpoint.InvalidHost", opt: config.WithMetricsPipeline(config.WithPipelineEndpoint("http://pineapples.pizza"))},
		{method: "WithPipelineHeaders.DuplicateEntries", opt: config.WithTracesPipeline(
			config.WithPipelineHeaders(map[string]string{"Otel-Service": "foo"}),
			config.WithPipelineHeaders(map[string]string{"Otel-Service": "bar"}),
		)},
		{method: "WithPipelineExporter", opt: config.WithMetricsPipeline(config.WithPipelineExporter(""))},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			assert.ErrorIs(t, config.NewDefault().Apply(tc.opt), config.ErrInvalidParam)
		})
	}
}
