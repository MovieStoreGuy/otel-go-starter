package main

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	instrumentationName = "collatz-conjecture"
	spanName            = "compute"
	cancelledEvent      = "context.cancelled"
)

type Conjecture struct {
	tracer trace.Tracer
}

func NewConjecture(tracerOpts ...trace.TracerOption) *Conjecture {
	return &Conjecture{
		tracer: otel.GetTracerProvider().Tracer(instrumentationName, tracerOpts...),
	}
}

func (c *Conjecture) ComputeRecursive(ctx context.Context, n int) {
	ctx, span := c.tracer.Start(ctx, spanName,
		trace.WithAttributes(attribute.Int("starting.value", n)),
	)
	defer span.End()

	select {
	case <-ctx.Done():
		span.AddEvent(cancelledEvent,
			trace.WithTimestamp(time.Now()),
		)
	default:
		// Allow to passthrough
	}

	switch {
	case n == 1:
		return
	case n%2 == 0:
		c.ComputeRecursive(ctx, n/2)
	case n%2 == 1:
		c.ComputeRecursive(ctx, (3*n+1)/2)
	}
}
