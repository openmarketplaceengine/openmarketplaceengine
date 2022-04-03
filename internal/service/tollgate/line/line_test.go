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

	t.Run("testNotCrossedPrecision", func(t *testing.T) {
		testNotCrossedPrecision(t)
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

	crossing := detectCrossingLine("", &line, &m, 0.001)
	assert.NotNil(t, crossing)
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

	crossing := detectCrossingLine("", &line, &m, 0.0000001)
	assert.NotNil(t, crossing)
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

	crossing := DetectCrossing("", &line, &m)
	assert.NotNil(t, crossing)
}

func testNotCrossedPrecision(t *testing.T) {
	line1 := tollgate.Line{
		Lon1: -74.195995,
		Lat1: 40.636916,
		Lon2: -74.198356,
		Lat2: 40.634408,
	}
	line2 := tollgate.Line{
		Lon1: -73.951378,
		Lat1: 40.855176,
		Lon2: -73.953223,
		Lat2: 40.848359,
	}
	m := tollgate.Movement{
		From: &tollgate.Location{
			Lon: -74.172478,
			Lat: 40.663041,
		},
		To: &tollgate.Location{
			Lon: -74.154282,
			Lat: 40.669812,
		},
	}

	crossing1 := DetectCrossing("", &line1, &m)
	assert.Nil(t, crossing1)

	crossing2 := DetectCrossing("", &line2, &m)
	assert.Nil(t, crossing2)
}
