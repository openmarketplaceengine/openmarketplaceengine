package validate

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestValidateString(t *testing.T) {
	err := String("test", "a", IsNotNull)
	require.NoError(t, err)

	err = String("test", "", IsNull)
	require.NoError(t, err)
}

func TestValidateLon(t *testing.T) {
	err := Float64("test", 0, IsLon)
	require.NoError(t, err)

	err = Float64("test", 200, IsLon)
	require.EqualError(t, err, "must be valid floats between -180 and 180")
}

func TestValidateLat(t *testing.T) {
	err := Float64("test", 0, IsLat)
	require.NoError(t, err)

	err = Float64("test", 93, IsLat)
	require.EqualError(t, err, "must be valid floats between -90 and 90")
}

func TestValidate(t *testing.T) {
	v := Validator{}
	v.ValidateString("worker_id", "", IsNotNull)
	v.ValidateTimestamp("timestamp", &timestamppb.Timestamp{
		Seconds: -1,
		Nanos:   0,
	})
	v.ValidateLon("lon", 1000.00)
	v.ValidateLat("lat", 1000.00)

	require.Len(t, v.Errors, 4)
}
