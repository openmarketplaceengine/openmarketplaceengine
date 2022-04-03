package util

import (
	"fmt"
	"math"
	"math/rand"
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
