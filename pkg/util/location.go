package util

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
)

func LatitudeInRange(r *rand.Rand, min, max float64) float64 {
	if min > max || min < -90 || min > 90 || max < -90 || max > 90 {
		panic(fmt.Errorf("invalid min=%v or max=%v range, must be valid floats and between -90 and 90", min, max))
	}
	return ToFixed(float64Range(r, min, max), 6)
}

func LongitudeInRange(r *rand.Rand, min, max float64) float64 {
	if min > max || min < -180 || min > 180 || max < -180 || max > 180 {
		panic(fmt.Errorf("invalid min=%v or max=%v range, must be valid floats and between -180 and 180", min, max))
	}
	return ToFixed(float64Range(r, min, max), 6)
}

func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return math.Floor(num*output) / output
}

func Round6(num float64) float64 {
	output := math.Pow(10, float64(6))
	return math.Floor(num*output) / output
}

func float64Range(r *rand.Rand, min, max float64) float64 {
	if min == max {
		return min
	}
	return r.Float64()*(max-min) + min
}

func IsLatitude(latitude float64) error {
	if latitude < -90 || latitude > 90 {
		return fmt.Errorf("invalid latitude %v, must be between -90 and 90", latitude)
	}
	return nil
}

func IsLongitude(longitude float64) error {
	if longitude < -180 || longitude > 180 {
		return fmt.Errorf("invalid longitude %v, must be between -180 and 180", longitude)
	}
	return nil
}

func ParseLatitude(s string) (float64, error) {
	latitude, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	err = IsLatitude(latitude)
	if err != nil {
		return 0, err
	}
	return Round6(latitude), nil
}

func ParseLongitude(s string) (float64, error) {
	longitude, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	err = IsLongitude(longitude)
	if err != nil {
		return 0, err
	}
	return Round6(longitude), nil
}
