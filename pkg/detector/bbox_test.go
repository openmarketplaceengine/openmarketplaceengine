package detector

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestDetectCrossingBBox(t *testing.T) {
	storage := NewMapStorage()
	ctx := context.Background()

	t.Run("testCrossing", func(t *testing.T) {
		testCrossing(ctx, t, storage)
	})
}

func testCrossing(ctx context.Context, t *testing.T, s Storage) {
	required := int32(4)
	bBoxes := []*BBox{{
		LonMin: 1,
		LatMin: 1,
		LonMax: 5,
		LatMax: 5,
	}, {
		LonMin: 5,
		LatMin: 5,
		LonMax: 10,
		LatMax: 10,
	}, {
		LonMin: 10,
		LatMin: 10,
		LonMax: 15,
		LatMax: 15,
	}, {
		LonMin: 15,
		LatMin: 15,
		LonMax: 20,
		LatMax: 20,
	}, {
		LonMin: 20,
		LatMin: 20,
		LonMax: 35,
		LatMax: 35,
	}}

	loc0 := Location{}

	loc1 := Location{
		Lon: 2,
		Lat: 2,
	}
	loc2 := Location{
		Lon: 7,
		Lat: 7,
	}
	loc3 := Location{
		Lon: 13,
		Lat: 13,
	}
	//noisy GPS at 4, skipping
	loc5 := Location{
		Lon: 33,
		Lat: 33,
	}
	workerID := uuid.NewString()
	tollgateID := "toll-1"
	c1, _ := detectCrossingBBox(ctx, s, tollgateID, workerID, bBoxes, required, &Movement{From: &loc0, To: &loc1})
	require.Nil(t, c1)

	c2, _ := detectCrossingBBox(ctx, s, tollgateID, workerID, bBoxes, required, &Movement{From: &loc0, To: &loc2})
	require.Nil(t, c2)

	c3, _ := detectCrossingBBox(ctx, s, tollgateID, workerID, bBoxes, required, &Movement{From: &loc0, To: &loc3})
	require.Nil(t, c3)

	c4, _ := detectCrossingBBox(ctx, s, tollgateID, workerID, bBoxes, required, &Movement{From: &loc0, To: &loc0})
	require.Nil(t, c4)

	key := storageKey(tollgateID, workerID)

	visits, err := s.Visits(ctx, key, len(bBoxes))
	require.NoError(t, err)

	assert.Equal(t, []int{1, 1, 1, 0, 0}, visits)

	c5, _ := detectCrossingBBox(ctx, s, tollgateID, workerID, bBoxes, required, &Movement{From: &loc0, To: &loc5})
	require.NotNil(t, c5)

	visits, err = s.Visits(ctx, key, len(bBoxes))
	require.NoError(t, err)

	assert.Nil(t, visits)
}
