package trace_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/MovieStoreGuy/otel-go-starter/config"
	"github.com/MovieStoreGuy/otel-go-starter/internal/pipeline/trace"
)

func TestLoadingPropagators(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		use      []string
	}{
		{scenario: "Using all propagators", use: []string{"b3", "baggage", "tracecontext", "ottrace"}},
		{scenario: "Using one propagator", use: []string{"b3"}},
	}

	for _, tc := range testCases {
		t.Run(tc.scenario, func(t *testing.T) {
			propagator, err := trace.NewPropagators(tc.use)

			assert.NoError(t, err, "Must not error when a valid propagator set is provided")
			assert.NotNil(t, propagator, "Must have a valid propagator returned")
		})
	}

	_, err := trace.NewPropagators(nil)
	assert.ErrorIs(t, err, config.ErrInvalidParam, "Must error when no values are provided")
}
