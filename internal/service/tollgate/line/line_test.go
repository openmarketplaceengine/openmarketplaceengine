package line

import (
	"testing"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate"

	"github.com/stretchr/testify/assert"
)

func TestTollgate(t *testing.T) {
	t.Run("testCrossed", func(t *testing.T) {
		testCrossed(t)
	})

	t.Run("testCrossedLatLong", func(t *testing.T) {
		testCrossedLatLong(t)
	})

	t.Run("testNotCrossed", func(t *testing.T) {
		testNotCrossed(t)
	})
}

func testCrossed(t *testing.T) {
	tol := NewTollgate(
		"",
		&tollgate.LocationXY{
			LongitudeX: 2,
			LatitudeY:  5,
		},
		&tollgate.LocationXY{
			LongitudeX: 5,
			LatitudeY:  2,
		},
	)

	m := tollgate.Movement{
		From: &tollgate.LocationXY{
			LongitudeX: 2,
			LatitudeY:  2,
		},
		To: &tollgate.LocationXY{
			LongitudeX: 5,
			LatitudeY:  5,
		},
	}

	crossing := detectCrossing(tol, &m, 0.001)
	assert.Equal(t, 3.5, crossing.Location.LongitudeX)
	assert.Equal(t, 3.5, crossing.Location.LatitudeY)
}

func testCrossedLatLong(t *testing.T) {
	tol := NewTollgate(
		"",
		&tollgate.LocationXY{
			LongitudeX: -79.870262,
			LatitudeY:  41.198497,
		},
		&tollgate.LocationXY{
			LongitudeX: -79.870218,
			LatitudeY:  41.200268,
		},
	)

	m := tollgate.Movement{
		From: &tollgate.LocationXY{
			LongitudeX: -79.87124651670456,
			LatitudeY:  41.199493331477335,
		},
		To: &tollgate.LocationXY{
			LongitudeX: -79.867927,
			LatitudeY:  41.199329,
		},
	}

	crossing := detectCrossing(tol, &m, 0.0000001)
	assert.InDelta(t, -79.8702, crossing.Location.LongitudeX, 0.001)
	assert.InDelta(t, 41.1994, crossing.Location.LatitudeY, 0.001)
}

func testNotCrossed(t *testing.T) {
	tol := NewTollgate(
		"",
		&tollgate.LocationXY{
			LongitudeX: 5,
			LatitudeY:  5,
		},
		&tollgate.LocationXY{
			LongitudeX: 5,
			LatitudeY:  5,
		},
	)

	m := tollgate.Movement{
		From: &tollgate.LocationXY{
			LongitudeX: 5,
			LatitudeY:  5,
		},
		To: &tollgate.LocationXY{
			LongitudeX: 5,
			LatitudeY:  5,
		},
	}

	crossing := detectCrossing(tol, &m, 0.001)
	assert.Nil(t, crossing)
}
