package metric_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MovieStoreGuy/otel-go-starter/config"
	"github.com/MovieStoreGuy/otel-go-starter/internal/pipeline/metric"
	"github.com/stretchr/testify/assert"
)

func TestConfiguringExporters(t *testing.T) {
	t.Parallel()

	s := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		_, err := io.Copy(rw, r.Body)
		assert.NoError(t, err, "Must not error when copying data")
	}))
	t.Cleanup(s.Close)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testCases := []struct {
		scenario string
		conf     *config.Export
	}{
		{scenario: "Stdout Exporter", conf: &config.Export{Named: "stdout"}},
		{scenario: "Basic grpc otlp exporter", conf: &config.Export{Named: "otlpgrpc", Endpoint: s.URL}},
		{scenario: "Basic http otlp exporter", conf: &config.Export{Named: "otlphttp", Endpoint: s.URL}},
		{
			scenario: "Configured grpc otlp exporter",
			conf: &config.Export{
				Named:    "otlpgrpc",
				Endpoint: s.URL,
				Headers: map[string]string{
					"Service-Domain": "icecream",
				},
				AllowInsecure:  true,
				UseCompression: true,
			},
		},
		{
			scenario: "Configured http otlp exporter",
			conf: &config.Export{
				Named:    "otlphttp",
				Endpoint: s.URL,
				Headers: map[string]string{
					"Service-Domain": "icecream",
				},
				AllowInsecure:  true,
				UseCompression: true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.scenario, func(t *testing.T) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			exporter, err := metric.NewExporterFactory().NewExporter(ctx, tc.conf)

			assert.NoError(t, err, "Must not error when ")
			assert.NotNil(t, exporter, "Must not return a nil exporter")
			if sh, ok := exporter.(metric.ShutdownExporter); ok {
				assert.NoError(t, sh.Shutdown(ctx), "Must not error when shutting down")
			}
		})
	}

	_, err := metric.NewExporterFactory().NewExporter(ctx, &config.Export{Named: "undefined-exporter"})
	assert.ErrorIs(t, err, metric.ErrNotDefinedExporter)
}
