package trace_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/MovieStoreGuy/otel-go-starter/config"
	"github.com/MovieStoreGuy/otel-go-starter/internal/pipeline/trace"
)

func TestBuildingExporters(t *testing.T) {
	t.Parallel()

	factory := trace.NewExporterFactory()
	require.NotNil(t, factory, "Must have a valid exporter factory")

	s := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		_, err := io.Copy(rw, r.Body)
		assert.NoError(t, err, "Must not error when copying data")
	}))
	t.Cleanup(s.Close)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testCases := []struct {
		name string
		conf *config.Export
	}{
		{name: "stdout exporter", conf: &config.Export{Named: "stdout"}},
		{name: "jaeger exporter", conf: &config.Export{Named: "jaeger", Endpoint: s.URL}},
		{name: "zipkin exporter", conf: &config.Export{Named: "zipkin", Endpoint: s.URL}},
		{name: "otel http exporter", conf: &config.Export{Named: "otlphttp", Endpoint: s.URL}},
		{name: "otel grpc exporter", conf: &config.Export{Named: "otlpgrpc", Endpoint: s.URL}},
		{
			name: "otel http exporter with all options",
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
		{
			name: "otel grpc exporter with all options",
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, done := context.WithCancel(ctx)
			exporter, err := factory.NewExporter(ctx, tc.conf)

			require.NoError(t, err, "Must not error when configuring exporter")
			require.NotNil(t, exporter, "Must have a valid exporter")

			assert.NoError(t, exporter.Shutdown(ctx), "Must not error when shutting down exporter")
			done()
		})
	}
	_, err := factory.NewExporter(context.Background(), &config.Export{})
	assert.ErrorIs(t, err, trace.ErrNotDefinedExporter, "Must error when invalid exporter name is provided")
}
