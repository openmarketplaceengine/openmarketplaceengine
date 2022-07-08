package validate

import (
	"fmt"
)

type Float64Wrap struct {
	validator *Validator
	name      string
	value     float64
}

func (w *Float64Wrap) Latitude() {
	if w.value < -90 || w.value > 90 {
		err := wrapError(w.name, w.value, fmt.Errorf("must be valid floats between -90 and 90"))
		w.validator.Errors = append(w.validator.Errors, err)
	}
}

func (w *Float64Wrap) Longitude() {
	if w.value < -180 || w.value > 180 {
		err := wrapError(w.name, w.value, fmt.Errorf("must be valid floats between -180 and 180"))
		w.validator.Errors = append(w.validator.Errors, err)
	}
}
