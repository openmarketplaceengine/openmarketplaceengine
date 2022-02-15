package itinerary

import (
	"testing"

	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/job/flow"

	"github.com/stretchr/testify/require"
)

func TestItinerary(t *testing.T) {
	itinerary := NewItinerary("i-1", "immediate", []flow.Step{})

	require.NotNil(t, itinerary)
}
