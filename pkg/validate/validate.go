package validate

import (
	"fmt"
	"strings"
	"time"
)

type Validator struct {
	Errors []error
}

func (v *Validator) Error() error {
	if len(v.Errors) == 0 {
		return nil
	}
	var errs = make([]string, len(v.Errors))
	for i, err := range v.Errors {
		errs[i] = err.Error()
	}
	return fmt.Errorf("ValidationErrors:\n%s", strings.Join(errs, "\n"))
}

func wrapError(name string, value interface{}, e error) error {
	return fmt.Errorf("%s=%v, %w", name, value, e)
}

func (v *Validator) ValidateString(name string, value string) *StrWrap {
	return &StrWrap{
		validator: v,
		name:      name,
		value:     value,
	}
}

func (v *Validator) ValidateFloat64(name string, value float64) *Float64Wrap {
	return &Float64Wrap{
		validator: v,
		name:      name,
		value:     value,
	}
}

func (v *Validator) ValidateInt32(name string, value int32) *Int32Wrap {
	return &Int32Wrap{
		validator: v,
		name:      name,
		value:     value,
	}
}

func (v *Validator) ValidateTime(name string, value time.Time) *TimeWrap {
	return &TimeWrap{
		validator: v,
		name:      name,
		value:     value,
	}
}

func (v *Validator) Validate(name string, value interface{}, rule func(value interface{}) error) {
	err := rule(value)
	if err != nil {
		v.Errors = append(v.Errors, wrapError(name, value, err))
	}
}
