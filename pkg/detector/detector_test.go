package detector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMovement(t *testing.T) {
	t.Run("testDirection", func(t *testing.T) {
		testDirection(t)
	})
}

func testDirection(t *testing.T) {
	m1 := &Movement{
		From: &Location{
			Longitude: -77.036400,
			Latitude:  38.895100,
		},
		To: &Location{
			Longitude: -77.036400,
			Latitude:  39.895100,
		},
	}
	assert.Equal(t, Direction("N"), m1.Direction())
	m2 := &Movement{
		From: &Location{
			Longitude: -77.036400,
			Latitude:  38.895100,
		},
		To: &Location{
			Longitude: -77.036400,
			Latitude:  37.895100,
		},
	}
	assert.Equal(t, Direction("S"), m2.Direction())
	m3 := &Movement{
		From: &Location{
			Longitude: -77.036400,
			Latitude:  38.895100,
		},
		To: &Location{
			Longitude: -76.036400,
			Latitude:  38.895100,
		},
	}
	assert.Equal(t, Direction("E"), m3.Direction())
	m4 := &Movement{
		From: &Location{
			Longitude: -77.036400,
			Latitude:  38.895100,
		},
		To: &Location{
			Longitude: -78.036400,
			Latitude:  38.895100,
		},
	}
	assert.Equal(t, Direction("W"), m4.Direction())
	m5 := &Movement{
		From: &Location{
			Longitude: -77.036400,
			Latitude:  38.895100,
		},
		To: &Location{
			Longitude: -78.036400,
			Latitude:  39.895100,
		},
	}
	assert.Equal(t, Direction("NW"), m5.Direction())
	m6 := &Movement{
		From: &Location{
			Longitude: -77.036400,
			Latitude:  38.895100,
		},
		To: &Location{
			Longitude: -78.036400,
			Latitude:  37.895100,
		},
	}
	assert.Equal(t, Direction("SW"), m6.Direction())
}
