package validate

import (
	"fmt"
	"time"
)

type timeWrap struct {
	validator *Validator
	name      string
	value     time.Time
}

func (w *timeWrap) NotBefore(before time.Time) {
	if w.value.Before(before) {
		err := wrapError(w.name, w.value, fmt.Errorf("must not be in the past"))
		w.validator.Errors = append(w.validator.Errors, err)
	}
}
