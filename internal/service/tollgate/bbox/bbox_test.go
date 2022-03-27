package bbox

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate"

	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	redisClient "github.com/openmarketplaceengine/openmarketplaceengine/redis/client"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestTollgate(t *testing.T) {
	err := cfg.Load()
	require.NoError(t, err)

	client := redisClient.NewStoreClient()
	require.NotNil(t, client)

	storage := NewStorage(client)
	ctx := context.Background()

	t.Run("testCrossing", func(t *testing.T) {
		testCrossing(ctx, t, storage)
	})
}

func testCrossing(ctx context.Context, t *testing.T, s Storage) {
	bBox, err := NewTollgate(
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

	loc0 := tollgate.LocationXY{}

	loc1 := tollgate.LocationXY{
		LongitudeX: 2,
		LatitudeY:  2,
	}
	loc2 := tollgate.LocationXY{
		LongitudeX: 7,
		LatitudeY:  7,
	}
	loc3 := tollgate.LocationXY{
		LongitudeX: 13,
		LatitudeY:  13,
	}
	//noisy GPS at 4, skipping
	loc5 := tollgate.LocationXY{
		LongitudeX: 33,
		LatitudeY:  33,
	}
	subjectID := uuid.NewString()
	c1, _ := bBox.Detect(ctx, &tollgate.Movement{SubjectID: subjectID, From: &loc0, To: &loc1})
	require.Nil(t, c1)

	c2, _ := bBox.Detect(ctx, &tollgate.Movement{SubjectID: subjectID, From: &loc0, To: &loc2})
	require.Nil(t, c2)

	c3, _ := bBox.Detect(ctx, &tollgate.Movement{SubjectID: subjectID, From: &loc0, To: &loc3})
	require.Nil(t, c3)

	c4, _ := bBox.Detect(ctx, &tollgate.Movement{SubjectID: subjectID, From: &loc0, To: &loc0})
	require.Nil(t, c4)

	visits, err := s.Visits(ctx, bBox.TollgateID, subjectID, 5)
	require.NoError(t, err)

	assert.Equal(t, []int64{1, 1, 1, 0, 0}, visits)

	c5, _ := bBox.Detect(ctx, &tollgate.Movement{SubjectID: subjectID, From: &loc0, To: &loc5})
	require.NotNil(t, c5)

	visits, err = s.Visits(ctx, bBox.TollgateID, subjectID, 5)
	require.NoError(t, err)

	assert.Equal(t, []int64{0, 0, 0, 0, 0}, visits)
}
