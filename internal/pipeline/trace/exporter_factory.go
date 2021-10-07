package trace

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/exporters/zipkin"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/encoding/gzip"

	"github.com/MovieStoreGuy/otel-go-starter/config"
)

var (
	ErrNotDefinedExporter = errors.New("invalid exporter provided")
)

type ExporterFactory map[string]generatorFunc

type generatorFunc func(ctx context.Context, conf *config.Pipeline) (sdktrace.SpanExporter, error)

func (ef ExporterFactory) NewExporter(ctx context.Context, conf *config.Pipeline) (sdktrace.SpanExporter, error) {
	factory, exist := (ef)[conf.Exporter]
	if !exist {
		return nil, fmt.Errorf("unknown exporter %s: %w", conf.Exporter, ErrNotDefinedExporter)
	}
	return factory(ctx, conf)
}

func NewExporterFactory() ExporterFactory {
	return map[string]generatorFunc{
		"otlpgrpc": func(ctx context.Context, conf *config.Pipeline) (sdktrace.SpanExporter, error) {
			var grpcOpts []otlptracegrpc.Option

			if endpoint := conf.Endpoint; endpoint != "" {
				grpcOpts = append(grpcOpts, otlptracegrpc.WithEndpoint(endpoint))
			}
			if headers := conf.Headers; len(headers) != 0 {
				grpcOpts = append(grpcOpts, otlptracegrpc.WithHeaders(headers))
			}
			if conf.AllowInsecure {
				grpcOpts = append(grpcOpts, otlptracegrpc.WithInsecure())
			}
			if conf.UseCompression {
				grpcOpts = append(grpcOpts, otlptracegrpc.WithCompressor(gzip.Name))
			}

			return otlptracegrpc.New(ctx, grpcOpts...)
		},
		"otlphttp": func(ctx context.Context, conf *config.Pipeline) (sdktrace.SpanExporter, error) {
			var httpOpts []otlptracehttp.Option

			if endpoint := conf.Endpoint; endpoint != "" {
				httpOpts = append(httpOpts, otlptracehttp.WithEndpoint(endpoint))
			}
			if headers := conf.Headers; len(headers) != 0 {
				httpOpts = append(httpOpts, otlptracehttp.WithHeaders(headers))
			}
			if conf.AllowInsecure {
				httpOpts = append(httpOpts, otlptracehttp.WithInsecure())
			}
			if conf.UseCompression {
				httpOpts = append(httpOpts, otlptracehttp.WithCompression(otlptracehttp.GzipCompression))
			}

			return otlptracehttp.New(ctx, httpOpts...)
		},
		"zipkin": func(ctx context.Context, conf *config.Pipeline) (sdktrace.SpanExporter, error) {
			return zipkin.New(conf.Endpoint)
		},
		"jaeger": func(ctx context.Context, conf *config.Pipeline) (sdktrace.SpanExporter, error) {
			return jaeger.New(jaeger.WithCollectorEndpoint(
				jaeger.WithEndpoint(conf.Endpoint),
			))
		},
		"stdout": func(_ context.Context, _ *config.Pipeline) (sdktrace.SpanExporter, error) {
			return stdouttrace.New(stdouttrace.WithPrettyPrint())
		},
	}
}
