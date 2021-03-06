// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: api/type/v1beta1/crossing.proto

package v1beta1

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on Crossing with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Crossing) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Crossing with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in CrossingMultiError, or nil
// if none found.
func (m *Crossing) ValidateAll() error {
	return m.validate(true)
}

func (m *Crossing) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for TollgateId

	// no validation rules for WorkerId

	// no validation rules for Direction

	// no validation rules for Alg

	if all {
		switch v := interface{}(m.GetMovement()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CrossingValidationError{
					field:  "Movement",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CrossingValidationError{
					field:  "Movement",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetMovement()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CrossingValidationError{
				field:  "Movement",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetCreateTime()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CrossingValidationError{
					field:  "CreateTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CrossingValidationError{
					field:  "CreateTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCreateTime()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CrossingValidationError{
				field:  "CreateTime",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return CrossingMultiError(errors)
	}

	return nil
}

// CrossingMultiError is an error wrapping multiple validation errors returned
// by Crossing.ValidateAll() if the designated constraints aren't met.
type CrossingMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CrossingMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CrossingMultiError) AllErrors() []error { return m }

// CrossingValidationError is the validation error returned by
// Crossing.Validate if the designated constraints aren't met.
type CrossingValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CrossingValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CrossingValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CrossingValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CrossingValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CrossingValidationError) ErrorName() string { return "CrossingValidationError" }

// Error satisfies the builtin error interface
func (e CrossingValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCrossing.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CrossingValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CrossingValidationError{}
