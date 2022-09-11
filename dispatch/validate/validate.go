package validate

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var v = newValidator()

func newValidator() *validator.Validate {
	v := NewValidator()
	return v
}

type Error struct {
	Field   string `json:"field"`
	Value   string `json:"value"`
	Message string `json:"message"`
	Details string `json:"details"`
}

type Var struct {
	Field interface{}
	Name  string
	Tag   string
}

func Vars(variables ...Var) []Error {
	var fieldErrors = make([]validator.FieldError, 0)
	for _, variable := range variables {
		err := v.Var(variable.Field, variable.Tag)
		if err == nil {
			continue
		}
		// nolint: errorlint
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			panic(fmt.Errorf("expected validator.ValidationErrors but received %T error: %w", err, err))
		}
		if len(validationErrors) > 0 {
			fieldErrors = append(fieldErrors, validationErrors...)
		}
	}

	return toErrors(fieldErrors, "")
}

func Struct(s interface{}, trimNsPrefix string) []Error {
	var vErrs = make([]Error, 0)
	err := v.Struct(s)
	if err == nil {
		return nil
	}
	// nolint: errorlint
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		msg := fmt.Sprintf("expected validator.ValidationErrors but received %T", err)
		vErrs = append(vErrs, Error{
			Field:   "unknown",
			Value:   "unknown",
			Message: msg,
			Details: msg,
		})

		return vErrs
	}

	return toErrors(validationErrors, trimNsPrefix)
}

func toErrors(fieldErrors []validator.FieldError, trimNsPrefix string) []Error {
	var errs = make([]Error, 0)
	for _, fieldError := range fieldErrors {
		k := strings.TrimPrefix(fieldError.Namespace(), trimNsPrefix)
		if k == "" {
			k = fieldError.ActualTag()
		}
		v := fmt.Sprintf("%v", fieldError.Value())
		errs = append(errs, Error{
			Field:   k,
			Value:   v,
			Message: fieldError.Error(),
			Details: fieldError.Error(),
		})
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}
