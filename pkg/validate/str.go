package validate

import (
	"fmt"
)

type StrWrap struct {
	validator *Validator
	name      string
	value     string
}

func (w *StrWrap) NotEmpty() {
	if len(w.value) == 0 {
		err := wrapError(w.name, w.value, fmt.Errorf("must not be empty"))
		w.validator.Errors = append(w.validator.Errors, err)
	}
}

func (w *StrWrap) Empty() {
	if len(w.value) > 0 {
		err := wrapError(w.name, w.value, fmt.Errorf("must be empty"))
		w.validator.Errors = append(w.validator.Errors, err)
	}
}

func (w *StrWrap) LenLessThan(limit int) {
	if len(w.value) > limit {
		err := wrapError(w.name, w.value, fmt.Errorf("length must be less than %v", limit))
		w.validator.Errors = append(w.validator.Errors, err)
	}
}
