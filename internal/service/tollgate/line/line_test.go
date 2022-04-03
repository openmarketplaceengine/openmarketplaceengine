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
	line := tollgate.Line{
		Lon1: 2,
		Lat1: 5,
		Lon2: 5,
		Lat2: 2,
	}

	m := tollgate.Movement{
		From: &tollgate.Location{
			Lon: 2,
			Lat: 2,
		},
		To: &tollgate.Location{
			Lon: 5,
			Lat: 5,
		},
	}

	crossing := DetectCrossing("", &line, &m, 0.001)
	assert.Equal(t, 3.5, crossing.Location.Lon)
	assert.Equal(t, 3.5, crossing.Location.Lat)
}

func testCrossedLatLong(t *testing.T) {
	line := tollgate.Line{
		Lon1: -79.870262,
		Lat1: 41.198497,
		Lon2: -79.870218,
		Lat2: 41.200268,
	}

	m := tollgate.Movement{
		From: &tollgate.Location{
			Lon: -79.87124651670456,
			Lat: 41.199493331477335,
		},
		To: &tollgate.Location{
			Lon: -79.867927,
			Lat: 41.199329,
		},
	}

	crossing := DetectCrossing("", &line, &m, 0.0000001)
	assert.InDelta(t, -79.8702, crossing.Location.Lon, 0.001)
	assert.InDelta(t, 41.1994, crossing.Location.Lat, 0.001)
}

func testNotCrossed(t *testing.T) {
	line := tollgate.Line{
		Lon1: 5,
		Lat1: 5,
		Lon2: 5,
		Lat2: 5,
	}
	m := tollgate.Movement{
		From: &tollgate.Location{
			Lon: 5,
			Lat: 5,
		},
		To: &tollgate.Location{
			Lon: 5,
			Lat: 5,
		},
	}

	crossing := DetectCrossing("", &line, &m, 0.001)
	assert.Nil(t, crossing)
}
