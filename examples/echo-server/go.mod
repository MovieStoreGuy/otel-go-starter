module github.com/MovieStoreGuy/otel-go-starter/examples/echo-server

go 1.17

replace github.com/MovieStoreGuy/otel-go-starter => ../..

require (
	github.com/MovieStoreGuy/otel-go-starter v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.8.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.25.0
)

require (
	github.com/cenkalti/backoff/v4 v4.1.1 // indirect
	github.com/felixge/httpsnoop v1.0.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/openzipkin/zipkin-go v0.2.5 // indirect
	go.opentelemetry.io/contrib/propagators/b3 v1.0.0 // indirect
	go.opentelemetry.io/contrib/propagators/ot v1.0.0 // indirect
	go.opentelemetry.io/otel v1.0.1 // indirect
	go.opentelemetry.io/otel/exporters/jaeger v1.0.1 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.0.1 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.0.1 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.0.1 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.0.1 // indirect
	go.opentelemetry.io/otel/exporters/zipkin v1.0.1 // indirect
	go.opentelemetry.io/otel/sdk v1.0.1 // indirect
	go.opentelemetry.io/otel/trace v1.0.1 // indirect
	go.opentelemetry.io/proto/otlp v0.9.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20210604141403-392c879c8b08 // indirect
	google.golang.org/grpc v1.41.0 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)
