package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	otelstarter "github.com/MovieStoreGuy/otel-go-starter"
	"github.com/MovieStoreGuy/otel-go-starter/config"
)

const (
	instrumentationName = "collatz-conjecture"

	spanNameRecursive = "compute.recursive"
	spanNameMemoized  = "compute.memoized"
	spanNameIterative = "compute.iterative"

	cancelledEvent   = "context.cancelled"
	precomputedEvent = "result.cached"
)

var (
	workers  int = runtime.NumCPU()
	maxValue int = 50_000
)

type Conjecture struct {
	tracer trace.Tracer

	rw      sync.RWMutex
	memoize map[int]struct{}
}

func NewConjecture(tracerOpts ...trace.TracerOption) *Conjecture {
	return &Conjecture{
		tracer: otel.GetTracerProvider().Tracer(instrumentationName, tracerOpts...),
		memoize: map[int]struct{}{
			1: {},
			2: {},
			4: {},
		},
	}
}

func (c *Conjecture) ComputeRecursive(ctx context.Context, n int) {
	ctx, span := c.tracer.Start(ctx, spanNameRecursive,
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

func (c *Conjecture) ComputeMemoized(ctx context.Context, n int) {
	ctx, span := c.tracer.Start(ctx, spanNameMemoized,
		trace.WithAttributes(attribute.Int("starting.value", n)),
	)

	defer span.End()

	select {
	case <-ctx.Done():
		span.AddEvent(cancelledEvent,
			trace.WithTimestamp(time.Now()),
		)
		return
	default:
		// Allow to passthrough
	}
	c.rw.RLock()
	_, computed := c.memoize[n]
	c.rw.RUnlock()

	if computed {
		span.AddEvent(precomputedEvent,
			trace.WithTimestamp(time.Now()),
		)
		return
	}

	// Since the result isn't needed but the fact that we have
	// previously computed it then as an optimisation
	// we'll register that we have done it to stop other branches from doing so
	c.rw.Lock()
	c.memoize[n] = struct{}{}
	c.rw.Unlock()

	switch n%2 == 0 { // IsEven?
	case true:
		c.ComputeMemoized(ctx, n/2)
	case false:
		c.ComputeMemoized(ctx, (3*n+1)/2)
	}
}

func (c *Conjecture) ComputeIterative(ctx context.Context, n int) {
	ctx, span := c.tracer.Start(ctx, spanNameIterative,
		trace.WithAttributes(attribute.Int("starting.value", n)),
	)
	defer span.End()

	for n != 1 {
		select {
		case <-ctx.Done():
			span.AddEvent(cancelledEvent,
				trace.WithTimestamp(time.Now()),
			)
			return
		default:
			// All to passthrough
		}

		switch n%2 == 0 { // IsEven?
		case true:
			n = n / 2
		case false:
			n = (3*n + 1) / 2
		}
	}
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt,
	)
	defer cancel()

	defer otelstarter.Start(ctx,
		config.WithServiceName("collatz-conjecture"),
		config.WithTracesPipeline(
			config.WithPipelineEnabled(),
			config.WithPipelineExporter("zipkin"),
			config.WithPipelineEndpoint("http://localhost:9411/api/v2/spans"),
		),
	).Shutdown()

	con := NewConjecture()

	var wg sync.WaitGroup
	jobs := make(chan int, maxValue)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				con.ComputeRecursive(ctx, job)
				con.ComputeMemoized(ctx, job)
				con.ComputeIterative(ctx, job)
			}
		}()
	}

	for n := 1; n < maxValue; n++ {
		jobs <- n
	}
	close(jobs)
	wg.Wait()

	fmt.Println("Finished calculating values for the Collatz Conjecture")
}
