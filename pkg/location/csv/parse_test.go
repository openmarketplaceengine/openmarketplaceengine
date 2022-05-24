package csv

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	line := "4026b7e2b949a25312a6ba7a038e852b1a927ff6,40.751251220703125,-73.99410206673768,2022-04-01T03:35:51.312816Z"

	location, err := Parse(line)
	require.NoError(t, err)
	assert.Equal(t, "4026b7e2b949a25312a6ba7a038e852b1a927ff6", location.DriverID)
	assert.Equal(t, 40.751251, location.Latitude)
	assert.Equal(t, -73.994103, location.Longitude)
	assert.Equal(t, time.Unix(1648784151, 0).Unix(), location.Timestamp.Unix())
}
