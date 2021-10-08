package metric

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	sdkmetric "go.opentelemetry.io/otel/sdk/export/metric"
	"google.golang.org/grpc/encoding/gzip"

	"github.com/MovieStoreGuy/otel-go-starter/config"
)

var ErrNotDefinedExporter = errors.New("invalid exporter provided")

type Factory map[string]generatorFunc

type generatorFunc func(ctx context.Context, pipe *config.Pipeline) (sdkmetric.Exporter, error)

func (ef Factory) NewExporter(ctx context.Context, pipe *config.Pipeline) (sdkmetric.Exporter, error) {
	factory, exist := ef[pipe.Exporter]
	if !exist {
		return nil, fmt.Errorf("unknown exporter %s: %w", pipe.Exporter, ErrNotDefinedExporter)
	}
	return factory(ctx, pipe)
}

func NewExporterFactory() Factory {
	return map[string]generatorFunc{
		"stdout": func(ctx context.Context, pipe *config.Pipeline) (sdkmetric.Exporter, error) {
			return stdoutmetric.New(stdoutmetric.WithPrettyPrint())
		},
		"otlpgrpc": func(ctx context.Context, pipe *config.Pipeline) (sdkmetric.Exporter, error) {
			var grpcOpts []otlpmetricgrpc.Option

			if endpoint := pipe.Endpoint; endpoint != "" {
				grpcOpts = append(grpcOpts, otlpmetricgrpc.WithEndpoint(endpoint))
			}
			if headers := pipe.Headers; headers != nil {
				grpcOpts = append(grpcOpts, otlpmetricgrpc.WithHeaders(headers))
			}
			if pipe.AllowInsecure {
				grpcOpts = append(grpcOpts, otlpmetricgrpc.WithInsecure())
			}
			if pipe.UseCompression {
				grpcOpts = append(grpcOpts, otlpmetricgrpc.WithCompressor(gzip.Name))
			}

			return otlpmetricgrpc.New(ctx, grpcOpts...)
		},
		"otlphttp": func(ctx context.Context, pipe *config.Pipeline) (sdkmetric.Exporter, error) {
			var httpOpts []otlpmetrichttp.Option

			if endpoint := pipe.Endpoint; endpoint != "" {
				httpOpts = append(httpOpts, otlpmetrichttp.WithEndpoint(endpoint))
			}
			if headers := pipe.Headers; len(headers) != 0 {
				httpOpts = append(httpOpts, otlpmetrichttp.WithHeaders(headers))
			}
			if pipe.AllowInsecure {
				httpOpts = append(httpOpts, otlpmetrichttp.WithInsecure())
			}
			if pipe.UseCompression {
				httpOpts = append(httpOpts, otlpmetrichttp.WithCompression(otlpmetrichttp.GzipCompression))
			}

			return otlpmetrichttp.New(ctx, httpOpts...)
		},
	}
}
