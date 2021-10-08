package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"

	otelstarter "github.com/MovieStoreGuy/otel-go-starter"
	"github.com/MovieStoreGuy/otel-go-starter/config"
)

var (
	workers  int = runtime.NumCPU()
	maxValue int = 1_000_000
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt,
	)
	defer cancel()

	defer otelstarter.Start(ctx,
		config.WithTracesPipeline(
			config.WithPipelineEnabled(),
			config.WithPipelineExporter("stdout"),
		),
	).Shutdown()

	con := NewConjecture()

	var wg sync.WaitGroup
	jobs := make(chan int, workers)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				con.ComputeRecursive(ctx, job)
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
