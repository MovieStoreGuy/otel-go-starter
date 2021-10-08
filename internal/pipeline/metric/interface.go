package metric

import (
	"context"

	sdkmetric "go.opentelemetry.io/otel/sdk/export/metric"
)

// ShutdownExporter is used to help check if an exporter
// is able to be shutdown.
type ShutdownExporter interface {
	sdkmetric.Exporter

	Shutdown(ctx context.Context) error
}
