package util

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
)

func LatInRange(r *rand.Rand, min, max float64) float64 {
	if min > max || min < -90 || min > 90 || max < -90 || max > 90 {
		panic(fmt.Errorf("invalid min=%v or max=%v range, must be valid floats and between -90 and 90", min, max))
	}
	return ToFixed(float64Range(r, min, max), 6)
}

func LonInRange(r *rand.Rand, min, max float64) float64 {
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

func IsLat(lat float64) error {
	if lat < -90 || lat > 90 {
		return fmt.Errorf("invalid latitude %v, must be between -90 and 90", lat)
	}
	return nil
}

func IsLon(lon float64) error {
	if lon < -180 || lon > 180 {
		return fmt.Errorf("invalid longitude %v, must be between -180 and 180", lon)
	}
	return nil
}

func ParseLat(s string) (float64, error) {
	lat, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	err = IsLat(lat)
	if err != nil {
		return 0, err
	}
	return Round6(lat), nil
}

func ParseLon(s string) (float64, error) {
	lon, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	err = IsLon(lon)
	if err != nil {
		return 0, err
	}
	return Round6(lon), nil
}
