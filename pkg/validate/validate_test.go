package validate

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestValidateString(t *testing.T) {
	v := Validator{}
	v.ValidateString("test", "a").Empty()
	require.Len(t, v.Errors, 1)
	//require.ErrorIs(t, v.Errors[0], errors.New("must be empty"))
	require.Equal(t, v.Errors[0].Error(), "ValidationError: test=a, must be empty")

	v.ValidateString("test", "").NotEmpty()
	require.Len(t, v.Errors, 2)
	//require.ErrorIs(t, v.Errors[1], errors.New("must be not empty"))
	require.Equal(t, v.Errors[1].Error(), "ValidationError: test=, must not be empty")
}

func TestValidateInt32F(t *testing.T) {
	v := Validator{}
	v.ValidateInt32("a", 13).GreaterThan(14)
	v.ValidateInt32("a", 13).LessThan(10)

	require.Len(t, v.Errors, 2)
	assert.Equal(t, v.Errors[0].Error(), "ValidationError: a=13, 13 must be greater than 14")
	assert.Equal(t, v.Errors[1].Error(), "ValidationError: a=13, 13 must be less than 10")
}

func TestValidateInt32(t *testing.T) {
	v := Validator{}
	v.ValidateInt32("a", 13).GreaterThan(10)
	v.ValidateInt32("a", 13).LessThan(14)

	require.Len(t, v.Errors, 0)
}

func TestValidateLon(t *testing.T) {
	v := Validator{}
	v.ValidateFloat64("lng", 0).Longitude()
	require.Len(t, v.Errors, 0)

	v.ValidateFloat64("lng", 200).Longitude()
	require.Len(t, v.Errors, 1)
	require.EqualError(t, v.Errors[0], "ValidationError: lng=200, must be valid floats between -180 and 180")
}

func TestValidateLat(t *testing.T) {
	v := Validator{}
	v.ValidateFloat64("lat", 0).Latitude()
	require.Len(t, v.Errors, 0)

	v.ValidateFloat64("lat", 93).Latitude()
	require.Len(t, v.Errors, 1)
	require.EqualError(t, v.Errors[0], "ValidationError: lat=93, must be valid floats between -90 and 90")
}

func TestValidate(t *testing.T) {
	v := Validator{}
	v.ValidateString("a", "").NotEmpty()
	v.ValidateString("b", "12345").LenLessThan(3)
	v.ValidateTime("timestamp", time.Time{}).NotBefore(time.Now())
	v.ValidateFloat64("lng", 1000.00).Longitude()
	v.ValidateFloat64("lat", 1000.00).Latitude()

	require.Len(t, v.Errors, 5)
}
