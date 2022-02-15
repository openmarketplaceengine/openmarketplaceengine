package itinerary

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestItinerary(t *testing.T) {
	itinerary := NewItinerary("i-1", "immediate")

	require.NotNil(t, itinerary)
}
