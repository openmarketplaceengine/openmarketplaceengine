package validate

import (
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
	"testing"
)

type upper struct {
	Email string `json:"email" validate:"required,email"`
	Inner inner  `json:"inner" validate:"required,alpha"`
}

type inner struct {
	Number string `json:"number" validate:"required"`
	Code   string `json:"code" validate:"required"`
}

func TestValidateStruct(t *testing.T) {
	v := NewValidator()
	s := upper{
		Email: "",
		Inner: inner{
			Number: "",
			Code:   "",
		},
	}
	err := v.Struct(s)
	require.Error(t, err)
	require.Len(t, err.(validator.ValidationErrors), 3)
}
