package validate

import (
	"github.com/stretchr/testify/assert"
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

func TestValidateInt32F(t *testing.T) {
	v := Validator{}
	v.ValidateInt32("a", 13).GreaterThan(14)
	v.ValidateInt32("a", 13).LessThan(10)

	require.Len(t, v.Errors, 2)
	assert.Equal(t, v.Errors[0].Error(), "13 must be greater than 14")
	assert.Equal(t, v.Errors[1].Error(), "13 must be less than 10")
}

func TestValidateInt32(t *testing.T) {
	v := Validator{}
	v.ValidateInt32("a", 13).GreaterThan(10)
	v.ValidateInt32("a", 13).LessThan(14)

	require.Len(t, v.Errors, 0)
}

func TestValidateLon(t *testing.T) {
	err := Float64("test", 0, IsLongitude)
	require.NoError(t, err)

	err = Float64("test", 200, IsLongitude)
	require.EqualError(t, err, "must be valid floats between -180 and 180")
}

func TestValidateLat(t *testing.T) {
	err := Float64("test", 0, IsLatitude)
	require.NoError(t, err)

	err = Float64("test", 93, IsLatitude)
	require.EqualError(t, err, "must be valid floats between -90 and 90")
}

func TestValidate(t *testing.T) {
	v := Validator{}
	v.ValidateString("a", "").NotEmpty()
	v.ValidateString("b", "12345").LenLessThan(3)
	v.ValidateTimestamp("timestamp", &timestamppb.Timestamp{
		Seconds: -1,
		Nanos:   0,
	})
	v.ValidateLongitude("longitude", 1000.00)
	v.ValidateLatitude("latitude", 1000.00)

	require.Len(t, v.Errors, 5)
}
