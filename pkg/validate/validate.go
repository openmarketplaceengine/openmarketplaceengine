package validate

import (
	"fmt"
	"time"
)

type Validator struct {
	Errors []error
}

func wrapToValidationError(name string, value interface{}, e error) error {
	return fmt.Errorf("ValidationError: %s=%v, %w", name, value, e)
}

func (v *Validator) ValidateString(name string, value string) *strWrap {
	return &strWrap{
		validator: v,
		name:      name,
		value:     value,
	}
}

func (v *Validator) ValidateFloat64(name string, value float64) *float64Wrap {
	return &float64Wrap{
		validator: v,
		name:      name,
		value:     value,
	}
}

func (v *Validator) ValidateInt32(name string, value int32) *int32Wrap {
	return &int32Wrap{
		validator: v,
		name:      name,
		value:     value,
	}
}

func (v *Validator) ValidateTime(name string, value time.Time) *timeWrap {
	return &timeWrap{
		validator: v,
		name:      name,
		value:     value,
	}
}

func (v *Validator) Validate(name string, value interface{}, rule func(value interface{}) error) {
	err := rule(value)
	if err != nil {
		v.Errors = append(v.Errors, wrapToValidationError(name, value, err))
	}
}
