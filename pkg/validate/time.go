package validate

import (
	"fmt"
	"time"
)

type TimeWrap struct {
	validator *Validator
	name      string
	value     time.Time
}

func (w *TimeWrap) NotBefore(before time.Time) {
	if w.value.Before(before) {
		err := wrapError(w.name, w.value, fmt.Errorf("must not be in the past"))
		w.validator.Errors = append(w.validator.Errors, err)
	}
}
