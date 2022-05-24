package validate

import (
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func IsLat(lat float64) error {
	if lat < -90 || lat > 90 {
		return fmt.Errorf("must be valid floats between -90 and 90")
	}
	return nil
}

func IsLon(lon float64) error {
	if lon < -180 || lon > 180 {
		return fmt.Errorf("must be valid floats between -180 and 180")
	}
	return nil
}

func IsTimestamp(timestamp *timestamppb.Timestamp) error {
	if timestamp == nil {
		return fmt.Errorf("must not be null")
	}
	if timestamp.GetSeconds() < 0 {
		return fmt.Errorf("seconds must not be negative")
	}
	if timestamp.GetNanos() < 0 {
		return fmt.Errorf("nanos must not be negative")
	}
	return nil
}

func IsNull(str string) error {
	if len(str) > 0 {
		return fmt.Errorf("must be null")
	}
	return nil
}

func IsNotNull(str string) error {
	if len(str) == 0 {
		return fmt.Errorf("must not be null")
	}
	return nil
}

type ValidationError struct {
	Name  string
	Value interface{}
	Err   error
}

func (e ValidationError) Error() string {
	return e.Err.Error()
}

type Rule func(string) error

func String(name string, value string, rule func(value string) error) error {
	err := rule(value)
	if err != nil {
		return ValidationError{
			Name:  name,
			Value: value,
			Err:   err,
		}
	}
	return nil
}

func Float64(name string, value float64, rule func(value float64) error) error {
	err := rule(value)
	if err != nil {
		return ValidationError{
			Name:  name,
			Value: value,
			Err:   err,
		}
	}
	return nil
}

type Validator struct {
	Errors []error
}

func (v *Validator) ValidateString(name string, value string, rule func(value string) error) {
	err := String(name, value, rule)
	if err != nil {
		v.Errors = append(v.Errors, err)
	}
}

func (v *Validator) ValidateFloat64(name string, value float64, rule func(value float64) error) {
	err := Float64(name, value, rule)
	if err != nil {
		v.Errors = append(v.Errors, err)
	}
}

func (v *Validator) ValidateLat(name string, value float64) {
	err := Float64(name, value, IsLat)
	if err != nil {
		v.Errors = append(v.Errors, err)
	}
}

func (v *Validator) ValidateLon(name string, value float64) {
	err := Float64(name, value, IsLon)
	if err != nil {
		v.Errors = append(v.Errors, err)
	}
}

func (v *Validator) ValidateTimestamp(name string, timestamp *timestamppb.Timestamp) {
	err := IsTimestamp(timestamp)
	if err != nil {
		v.Errors = append(v.Errors, ValidationError{
			Name:  name,
			Value: timestamp,
			Err:   err,
		},
		)
	}
}

func (v *Validator) Validate(name string, value interface{}, rule func(value interface{}) error) {
	err := rule(value)
	if err != nil {
		v.Errors = append(v.Errors, ValidationError{
			Name:  name,
			Value: value,
			Err:   err,
		},
		)
	}
}

func (v *Validator) ErrorInfo() *errdetails.ErrorInfo {
	if len(v.Errors) == 0 {
		return nil
	}
	data := make(map[string]string, 0)
	for _, err := range v.Errors {
		vErr, ok := err.(ValidationError)
		if !ok {
			data[vErr.Name] = fmt.Errorf("ValidationError type check failed on error: %w", err).Error()
		} else {
			data[vErr.Name] = vErr.Error()
		}
	}
	return &errdetails.ErrorInfo{
		Reason:   "bad request",
		Domain:   "",
		Metadata: data,
	}
}

func (v *Validator) BadRequest(errors []error) *errdetails.BadRequest {
	if len(v.Errors) == 0 {
		return nil
	}
	violations := make([]*errdetails.BadRequest_FieldViolation, 0)
	for _, err := range errors {
		vErr, ok := err.(ValidationError)
		if !ok {
			violations = append(violations, &errdetails.BadRequest_FieldViolation{
				Field:       vErr.Name,
				Description: fmt.Errorf("ValidationError type check failed on error: %w", err).Error(),
			})
		} else {
			violations = append(violations, &errdetails.BadRequest_FieldViolation{
				Field:       vErr.Name,
				Description: vErr.Error(),
			})
		}
	}
	return &errdetails.BadRequest{
		FieldViolations: violations,
	}
}
