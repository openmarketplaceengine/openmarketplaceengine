package itinerary

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestItinerary(t *testing.T) {
	itinerary := NewItinerary()

	require.NotNil(t, itinerary)
}
