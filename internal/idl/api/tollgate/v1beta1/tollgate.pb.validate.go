// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: api/tollgate/v1beta1/tollgate.proto

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

// Validate checks the field values on GetTollgateRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetTollgateRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetTollgateRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetTollgateRequestMultiError, or nil if none found.
func (m *GetTollgateRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetTollgateRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for TollgateId

	if len(errors) > 0 {
		return GetTollgateRequestMultiError(errors)
	}

	return nil
}

// GetTollgateRequestMultiError is an error wrapping multiple validation errors
// returned by GetTollgateRequest.ValidateAll() if the designated constraints
// aren't met.
type GetTollgateRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetTollgateRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetTollgateRequestMultiError) AllErrors() []error { return m }

// GetTollgateRequestValidationError is the validation error returned by
// GetTollgateRequest.Validate if the designated constraints aren't met.
type GetTollgateRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetTollgateRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetTollgateRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetTollgateRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetTollgateRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetTollgateRequestValidationError) ErrorName() string {
	return "GetTollgateRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetTollgateRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetTollgateRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetTollgateRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetTollgateRequestValidationError{}

// Validate checks the field values on GetTollgateResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetTollgateResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetTollgateResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetTollgateResponseMultiError, or nil if none found.
func (m *GetTollgateResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetTollgateResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetTollgate()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, GetTollgateResponseValidationError{
					field:  "Tollgate",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, GetTollgateResponseValidationError{
					field:  "Tollgate",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetTollgate()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return GetTollgateResponseValidationError{
				field:  "Tollgate",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return GetTollgateResponseMultiError(errors)
	}

	return nil
}

// GetTollgateResponseMultiError is an error wrapping multiple validation
// errors returned by GetTollgateResponse.ValidateAll() if the designated
// constraints aren't met.
type GetTollgateResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetTollgateResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetTollgateResponseMultiError) AllErrors() []error { return m }

// GetTollgateResponseValidationError is the validation error returned by
// GetTollgateResponse.Validate if the designated constraints aren't met.
type GetTollgateResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetTollgateResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetTollgateResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetTollgateResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetTollgateResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetTollgateResponseValidationError) ErrorName() string {
	return "GetTollgateResponseValidationError"
}

// Error satisfies the builtin error interface
func (e GetTollgateResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetTollgateResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetTollgateResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetTollgateResponseValidationError{}

// Validate checks the field values on ListTollgatesRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ListTollgatesRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListTollgatesRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListTollgatesRequestMultiError, or nil if none found.
func (m *ListTollgatesRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *ListTollgatesRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for PageSize

	// no validation rules for PageToken

	if len(errors) > 0 {
		return ListTollgatesRequestMultiError(errors)
	}

	return nil
}

// ListTollgatesRequestMultiError is an error wrapping multiple validation
// errors returned by ListTollgatesRequest.ValidateAll() if the designated
// constraints aren't met.
type ListTollgatesRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListTollgatesRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListTollgatesRequestMultiError) AllErrors() []error { return m }

// ListTollgatesRequestValidationError is the validation error returned by
// ListTollgatesRequest.Validate if the designated constraints aren't met.
type ListTollgatesRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListTollgatesRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListTollgatesRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListTollgatesRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListTollgatesRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListTollgatesRequestValidationError) ErrorName() string {
	return "ListTollgatesRequestValidationError"
}

// Error satisfies the builtin error interface
func (e ListTollgatesRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListTollgatesRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListTollgatesRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListTollgatesRequestValidationError{}

// Validate checks the field values on ListTollgatesResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ListTollgatesResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListTollgatesResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListTollgatesResponseMultiError, or nil if none found.
func (m *ListTollgatesResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *ListTollgatesResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetTollgates() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ListTollgatesResponseValidationError{
						field:  fmt.Sprintf("Tollgates[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ListTollgatesResponseValidationError{
						field:  fmt.Sprintf("Tollgates[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ListTollgatesResponseValidationError{
					field:  fmt.Sprintf("Tollgates[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	// no validation rules for NextPageToken

	if len(errors) > 0 {
		return ListTollgatesResponseMultiError(errors)
	}

	return nil
}

// ListTollgatesResponseMultiError is an error wrapping multiple validation
// errors returned by ListTollgatesResponse.ValidateAll() if the designated
// constraints aren't met.
type ListTollgatesResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListTollgatesResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListTollgatesResponseMultiError) AllErrors() []error { return m }

// ListTollgatesResponseValidationError is the validation error returned by
// ListTollgatesResponse.Validate if the designated constraints aren't met.
type ListTollgatesResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListTollgatesResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListTollgatesResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListTollgatesResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListTollgatesResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListTollgatesResponseValidationError) ErrorName() string {
	return "ListTollgatesResponseValidationError"
}

// Error satisfies the builtin error interface
func (e ListTollgatesResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListTollgatesResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListTollgatesResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListTollgatesResponseValidationError{}

// Validate checks the field values on Tollgate with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Tollgate) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Tollgate with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in TollgateMultiError, or nil
// if none found.
func (m *Tollgate) ValidateAll() error {
	return m.validate(true)
}

func (m *Tollgate) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for Name

	if all {
		switch v := interface{}(m.GetBBoxes()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, TollgateValidationError{
					field:  "BBoxes",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, TollgateValidationError{
					field:  "BBoxes",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetBBoxes()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return TollgateValidationError{
				field:  "BBoxes",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetGateLine()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, TollgateValidationError{
					field:  "GateLine",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, TollgateValidationError{
					field:  "GateLine",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetGateLine()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return TollgateValidationError{
				field:  "GateLine",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetCreated()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, TollgateValidationError{
					field:  "Created",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, TollgateValidationError{
					field:  "Created",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCreated()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return TollgateValidationError{
				field:  "Created",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetUpdated()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, TollgateValidationError{
					field:  "Updated",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, TollgateValidationError{
					field:  "Updated",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetUpdated()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return TollgateValidationError{
				field:  "Updated",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return TollgateMultiError(errors)
	}

	return nil
}

// TollgateMultiError is an error wrapping multiple validation errors returned
// by Tollgate.ValidateAll() if the designated constraints aren't met.
type TollgateMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TollgateMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TollgateMultiError) AllErrors() []error { return m }

// TollgateValidationError is the validation error returned by
// Tollgate.Validate if the designated constraints aren't met.
type TollgateValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TollgateValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TollgateValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TollgateValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TollgateValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TollgateValidationError) ErrorName() string { return "TollgateValidationError" }

// Error satisfies the builtin error interface
func (e TollgateValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTollgate.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TollgateValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TollgateValidationError{}

// Validate checks the field values on BBoxes with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *BBoxes) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on BBoxes with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in BBoxesMultiError, or nil if none found.
func (m *BBoxes) ValidateAll() error {
	return m.validate(true)
}

func (m *BBoxes) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetBBoxes() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, BBoxesValidationError{
						field:  fmt.Sprintf("BBoxes[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, BBoxesValidationError{
						field:  fmt.Sprintf("BBoxes[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return BBoxesValidationError{
					field:  fmt.Sprintf("BBoxes[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	// no validation rules for Required

	if len(errors) > 0 {
		return BBoxesMultiError(errors)
	}

	return nil
}

// BBoxesMultiError is an error wrapping multiple validation errors returned by
// BBoxes.ValidateAll() if the designated constraints aren't met.
type BBoxesMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m BBoxesMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m BBoxesMultiError) AllErrors() []error { return m }

// BBoxesValidationError is the validation error returned by BBoxes.Validate if
// the designated constraints aren't met.
type BBoxesValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e BBoxesValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e BBoxesValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e BBoxesValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e BBoxesValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e BBoxesValidationError) ErrorName() string { return "BBoxesValidationError" }

// Error satisfies the builtin error interface
func (e BBoxesValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sBBoxes.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = BBoxesValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = BBoxesValidationError{}

// Validate checks the field values on BBox with the rules defined in the proto
// definition for this message. If any rules are violated, the first error
// encountered is returned, or nil if there are no violations.
func (m *BBox) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on BBox with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in BBoxMultiError, or nil if none found.
func (m *BBox) ValidateAll() error {
	return m.validate(true)
}

func (m *BBox) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for LonMin

	// no validation rules for LatMin

	// no validation rules for LonMax

	// no validation rules for LatMax

	if len(errors) > 0 {
		return BBoxMultiError(errors)
	}

	return nil
}

// BBoxMultiError is an error wrapping multiple validation errors returned by
// BBox.ValidateAll() if the designated constraints aren't met.
type BBoxMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m BBoxMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m BBoxMultiError) AllErrors() []error { return m }

// BBoxValidationError is the validation error returned by BBox.Validate if the
// designated constraints aren't met.
type BBoxValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e BBoxValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e BBoxValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e BBoxValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e BBoxValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e BBoxValidationError) ErrorName() string { return "BBoxValidationError" }

// Error satisfies the builtin error interface
func (e BBoxValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sBBox.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = BBoxValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = BBoxValidationError{}

// Validate checks the field values on GateLine with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *GateLine) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GateLine with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in GateLineMultiError, or nil
// if none found.
func (m *GateLine) ValidateAll() error {
	return m.validate(true)
}

func (m *GateLine) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Lon1

	// no validation rules for Lat1

	// no validation rules for Lon2

	// no validation rules for Lat2

	if len(errors) > 0 {
		return GateLineMultiError(errors)
	}

	return nil
}

// GateLineMultiError is an error wrapping multiple validation errors returned
// by GateLine.ValidateAll() if the designated constraints aren't met.
type GateLineMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GateLineMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GateLineMultiError) AllErrors() []error { return m }

// GateLineValidationError is the validation error returned by
// GateLine.Validate if the designated constraints aren't met.
type GateLineValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GateLineValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GateLineValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GateLineValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GateLineValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GateLineValidationError) ErrorName() string { return "GateLineValidationError" }

// Error satisfies the builtin error interface
func (e GateLineValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGateLine.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GateLineValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GateLineValidationError{}
