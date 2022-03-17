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

	t.Run("testDirection", func(t *testing.T) {
		testDirection(t)
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
		From: LocationXY{
			longitudeX: 2,
			latitudeY:  2,
		},
		To: LocationXY{
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
		From: LocationXY{
			longitudeX: 5,
			latitudeY:  5,
		},
		To: LocationXY{
			longitudeX: 5,
			latitudeY:  5,
		},
	}

	crossing := detectCrossing(&tol, &m, 0.001)
	assert.Nil(t, crossing)
}

func testDirection(t *testing.T) {
	assert.Equal(t, Direction("N"), detectDirection(&Movement{
		"",
		LocationXY{
			-77.036400,
			38.895100,
		},
		LocationXY{
			-77.036400,
			39.895100,
		},
	}))
	assert.Equal(t, Direction("S"), detectDirection(&Movement{
		"",
		LocationXY{
			-77.036400,
			38.895100,
		},
		LocationXY{
			-77.036400,
			37.895100,
		},
	}))
	assert.Equal(t, Direction("E"), detectDirection(&Movement{
		"",
		LocationXY{
			-77.036400,
			38.895100,
		},
		LocationXY{
			-76.036400,
			38.895100,
		},
	}))
	assert.Equal(t, Direction("W"), detectDirection(&Movement{
		"",
		LocationXY{
			-77.036400,
			38.895100,
		},
		LocationXY{
			-78.036400,
			38.895100,
		},
	}))
	assert.Equal(t, Direction("NW"), detectDirection(&Movement{
		"",
		LocationXY{
			-77.036400,
			38.895100,
		},
		LocationXY{
			-78.036400,
			39.895100,
		},
	}))
	assert.Equal(t, Direction("SW"), detectDirection(&Movement{
		"",
		LocationXY{
			-77.036400,
			38.895100,
		},
		LocationXY{
			-78.036400,
			37.895100,
		},
	}))
}
