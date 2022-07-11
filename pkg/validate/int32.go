package validate

import (
	"fmt"
)

type Int32Wrap struct {
	validator *Validator
	name      string
	value     int32
}

func (w *Int32Wrap) GreaterThan(v int32) {
	if w.value < v {
		err := wrapError(w.name, w.value, fmt.Errorf("%v must be greater than %v", w.value, v))
		w.validator.Errors = append(w.validator.Errors, err)
	}
}

func (w *Int32Wrap) LessThan(v int32) {
	if w.value > v {
		err := wrapError(w.name, w.value, fmt.Errorf("%v must be less than %v", w.value, v))
		w.validator.Errors = append(w.validator.Errors, err)
	}
}
