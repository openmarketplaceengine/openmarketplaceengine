package validate

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// NewValidator creates a validator instance.
// NB! needs to be a single instance as it caches struct info.
func NewValidator() *validator.Validate {
	v := validator.New()
	// register function to get tag name from json tags.
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}

		return name
	})

	return v
}
