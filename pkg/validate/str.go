package validate

import (
	"fmt"
)

type strWrap struct {
	validator *Validator
	name      string
	value     string
}

func (w *strWrap) NotEmpty() {
	if len(w.value) == 0 {
		err := wrapToValidationError(w.name, w.value, fmt.Errorf("must not be empty"))
		w.validator.Errors = append(w.validator.Errors, err)
	}
}

func (w *strWrap) Empty() {
	if len(w.value) > 0 {
		err := wrapToValidationError(w.name, w.value, fmt.Errorf("must be empty"))
		w.validator.Errors = append(w.validator.Errors, err)
	}
}

func (w *strWrap) LenLessThan(limit int) {
	if len(w.value) > limit {
		err := wrapToValidationError(w.name, w.value, fmt.Errorf("length must be less than %v", limit))
		w.validator.Errors = append(w.validator.Errors, err)
	}
}
