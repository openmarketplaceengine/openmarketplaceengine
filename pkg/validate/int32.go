package validate

import (
	"fmt"
)

type int32Wrap struct {
	validator *Validator
	name      string
	value     int32
}

func (w *int32Wrap) GreaterThan(v int32) {
	if w.value < v {
		err := wrapToValidationError(w.name, w.value, fmt.Errorf("%v must be greater than %v", w.value, v))
		w.validator.Errors = append(w.validator.Errors, err)
	}
}

func (w *int32Wrap) LessThan(v int32) {
	if w.value > v {
		err := wrapToValidationError(w.name, w.value, fmt.Errorf("%v must be less than %v", w.value, v))
		w.validator.Errors = append(w.validator.Errors, err)
	}
}
