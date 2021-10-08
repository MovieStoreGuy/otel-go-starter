package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"

	launcher "github.com/MovieStoreGuy/otel-go-starter"
	"github.com/MovieStoreGuy/otel-go-starter/config"
)

func main() {
	// Configure Signal handling to shutdown graceful
	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt,
	)
	defer cancel()

	defer launcher.Start(ctx,
		config.WithTracesPipeline(
			config.WithPipelineEnabled(),
			config.WithPipelineExporter("stdout"),
		),
	).Shutdown()

	// Configuration of the router handler
	r := mux.NewRouter()
	r.HandleFunc("/echo", func(rw http.ResponseWriter, r *http.Request) {
		io.Copy(rw, r.Body)
	})
	r.Use(otelmux.Middleware("echo-server"))

	s := &http.Server{
		Addr:    ":4096",
		Handler: r,
	}

	fmt.Println("Starting echo server")

	// Running the http server in a background
	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("Issue with server: ", err)
		}
	}()

	// Gracefully shutdown the running http server
	<-ctx.Done()
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Shutting down the application")
}
