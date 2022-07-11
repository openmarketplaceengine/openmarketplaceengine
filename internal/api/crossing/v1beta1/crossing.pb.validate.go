// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: api/crossing/v1beta1/crossing.proto

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

// Validate checks the field values on ListCrossingsRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ListCrossingsRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListCrossingsRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListCrossingsRequestMultiError, or nil if none found.
func (m *ListCrossingsRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *ListCrossingsRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for TollgateId

	// no validation rules for WorkerId

	// no validation rules for PageSize

	// no validation rules for PageToken

	if len(errors) > 0 {
		return ListCrossingsRequestMultiError(errors)
	}

	return nil
}

// ListCrossingsRequestMultiError is an error wrapping multiple validation
// errors returned by ListCrossingsRequest.ValidateAll() if the designated
// constraints aren't met.
type ListCrossingsRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListCrossingsRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListCrossingsRequestMultiError) AllErrors() []error { return m }

// ListCrossingsRequestValidationError is the validation error returned by
// ListCrossingsRequest.Validate if the designated constraints aren't met.
type ListCrossingsRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListCrossingsRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListCrossingsRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListCrossingsRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListCrossingsRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListCrossingsRequestValidationError) ErrorName() string {
	return "ListCrossingsRequestValidationError"
}

// Error satisfies the builtin error interface
func (e ListCrossingsRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListCrossingsRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListCrossingsRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListCrossingsRequestValidationError{}

// Validate checks the field values on ListCrossingsResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ListCrossingsResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListCrossingsResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListCrossingsResponseMultiError, or nil if none found.
func (m *ListCrossingsResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *ListCrossingsResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetCrossings() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ListCrossingsResponseValidationError{
						field:  fmt.Sprintf("Crossings[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ListCrossingsResponseValidationError{
						field:  fmt.Sprintf("Crossings[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ListCrossingsResponseValidationError{
					field:  fmt.Sprintf("Crossings[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	// no validation rules for NextPageToken

	if len(errors) > 0 {
		return ListCrossingsResponseMultiError(errors)
	}

	return nil
}

// ListCrossingsResponseMultiError is an error wrapping multiple validation
// errors returned by ListCrossingsResponse.ValidateAll() if the designated
// constraints aren't met.
type ListCrossingsResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListCrossingsResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListCrossingsResponseMultiError) AllErrors() []error { return m }

// ListCrossingsResponseValidationError is the validation error returned by
// ListCrossingsResponse.Validate if the designated constraints aren't met.
type ListCrossingsResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListCrossingsResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListCrossingsResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListCrossingsResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListCrossingsResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListCrossingsResponseValidationError) ErrorName() string {
	return "ListCrossingsResponseValidationError"
}

// Error satisfies the builtin error interface
func (e ListCrossingsResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListCrossingsResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListCrossingsResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListCrossingsResponseValidationError{}
