package tollgate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTollgate(t *testing.T) {
	t.Run("testCrossed", func(t *testing.T) {
		testCrossed(t)
	})

	t.Run("testNotCrossed", func(t *testing.T) {
		testNotCrossed(t)
	})
}

func testCrossed(t *testing.T) {
	tol := Tollgate{
		ID: "",
		Point1: LocationXY{
			longitudeX: 2,
			latitudeY:  5,
		},
		Point2: LocationXY{
			longitudeX: 5,
			latitudeY:  2,
		},
	}

	m := Movement{
		PrevLocation: LocationXY{
			longitudeX: 2,
			latitudeY:  2,
		},
		CurrLocation: LocationXY{
			longitudeX: 5,
			latitudeY:  5,
		},
	}

	crossing := detectCrossing(&tol, &m, 0.001)
	assert.Equal(t, 3.5, crossing.Location.longitudeX)
	assert.Equal(t, 3.5, crossing.Location.latitudeY)
}

func testNotCrossed(t *testing.T) {
	tol := Tollgate{
		ID: "",
		Point1: LocationXY{
			longitudeX: 5,
			latitudeY:  5,
		},
		Point2: LocationXY{
			longitudeX: 5,
			latitudeY:  5,
		},
	}

	m := Movement{
		PrevLocation: LocationXY{
			longitudeX: 5,
			latitudeY:  5,
		},
		CurrLocation: LocationXY{
			longitudeX: 5,
			latitudeY:  5,
		},
	}

	crossing := detectCrossing(&tol, &m, 0.001)
	assert.Nil(t, crossing)
}
