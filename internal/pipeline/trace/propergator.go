package trace

import (
	"fmt"

	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/contrib/propagators/ot"
	"go.opentelemetry.io/otel/propagation"

	"github.com/MovieStoreGuy/otel-go-starter/config"
)

func NewPropagator(use []string) (propagation.TextMapPropagator, error) {
	propergatorMap := map[string]propagation.TextMapPropagator{
		"b3":           b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader)),
		"baggage":      propagation.Baggage{},
		"tracecontext": propagation.TraceContext{},
		"ottrace":      ot.OT{},
	}
	var props []propagation.TextMapPropagator
	for _, key := range use {
		if prop, exist := propergatorMap[key]; exist {
			props = append(props, prop)
		}
	}

	if len(props) == 0 {
		return nil, fmt.Errorf("missing propagator values: %w", config.ErrInvalidParam)
	}

	return propagation.NewCompositeTextMapPropagator(props...), nil
}
