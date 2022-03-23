package bbox

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	redisClient "github.com/openmarketplaceengine/openmarketplaceengine/redis/client"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestDetector(t *testing.T) {
	err := cfg.Load()
	require.NoError(t, err)

	client := redisClient.NewStoreClient()
	require.NotNil(t, client)

	storage := NewStorage(client)
	ctx := context.Background()

	t.Run("testPassThrough", func(t *testing.T) {
		testPassThrough(ctx, t, storage)
	})
}

func testPassThrough(ctx context.Context, t *testing.T, s Storage) {
	d, err := NewDetector(
		"toll-123",
		[]*BBox{{
			Left:   1,
			Bottom: 1,
			Right:  5,
			Top:    5,
		}, {
			Left:   5,
			Bottom: 5,
			Right:  10,
			Top:    10,
		}, {
			Left:   10,
			Bottom: 10,
			Right:  15,
			Top:    15,
		}, {
			Left:   15,
			Bottom: 15,
			Right:  20,
			Top:    20,
		}, {
			Left:   20,
			Bottom: 20,
			Right:  35,
			Top:    35,
		}},
		4,
		s,
	)

	require.NoError(t, err)

	loc1 := LocationXY{
		LongitudeX: 2,
		LatitudeY:  2,
	}
	loc2 := LocationXY{
		LongitudeX: 7,
		LatitudeY:  7,
	}
	loc3 := LocationXY{
		LongitudeX: 13,
		LatitudeY:  13,
	}
	//noisy GPS at 4, skipping
	loc5 := LocationXY{
		LongitudeX: 33,
		LatitudeY:  33,
	}
	subjectID := uuid.NewString()
	passedThrough1, _ := d.Trace(ctx, subjectID, &loc1)
	require.Nil(t, passedThrough1)

	passedThrough2, _ := d.Trace(ctx, subjectID, &loc2)
	require.Nil(t, passedThrough2)

	passedThrough3, _ := d.Trace(ctx, subjectID, &loc3)
	require.Nil(t, passedThrough3)

	passedThrough4, _ := d.Trace(ctx, subjectID, &LocationXY{})
	require.Nil(t, passedThrough4)

	passedThrough5, _ := d.Trace(ctx, subjectID, &loc5)
	require.NotNil(t, passedThrough5)

	assert.Equal(t, []int64{1, 1, 1, 0, 1}, passedThrough5.Visits)
}
