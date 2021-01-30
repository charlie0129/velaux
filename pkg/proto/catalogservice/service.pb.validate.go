// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: pkg/proto/catalogservice/service.proto

package catalogservice

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang/protobuf/ptypes"
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
	_ = ptypes.DynamicAny{}
)

// define the regex for a UUID once up-front
var _service_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on AddCatalogRequest with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *AddCatalogRequest) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Name

	// no validation rules for Desc

	return nil
}

// AddCatalogRequestValidationError is the validation error returned by
// AddCatalogRequest.Validate if the designated constraints aren't met.
type AddCatalogRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e AddCatalogRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e AddCatalogRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e AddCatalogRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e AddCatalogRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e AddCatalogRequestValidationError) ErrorName() string {
	return "AddCatalogRequestValidationError"
}

// Error satisfies the builtin error interface
func (e AddCatalogRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAddCatalogRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = AddCatalogRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = AddCatalogRequestValidationError{}

// Validate checks the field values on AddCatalogResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *AddCatalogResponse) Validate() error {
	if m == nil {
		return nil
	}

	return nil
}

// AddCatalogResponseValidationError is the validation error returned by
// AddCatalogResponse.Validate if the designated constraints aren't met.
type AddCatalogResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e AddCatalogResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e AddCatalogResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e AddCatalogResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e AddCatalogResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e AddCatalogResponseValidationError) ErrorName() string {
	return "AddCatalogResponseValidationError"
}

// Error satisfies the builtin error interface
func (e AddCatalogResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAddCatalogResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = AddCatalogResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = AddCatalogResponseValidationError{}

// Validate checks the field values on GetCatalogRequest with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *GetCatalogRequest) Validate() error {
	if m == nil {
		return nil
	}

	return nil
}

// GetCatalogRequestValidationError is the validation error returned by
// GetCatalogRequest.Validate if the designated constraints aren't met.
type GetCatalogRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetCatalogRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetCatalogRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetCatalogRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetCatalogRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetCatalogRequestValidationError) ErrorName() string {
	return "GetCatalogRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetCatalogRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetCatalogRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetCatalogRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetCatalogRequestValidationError{}

// Validate checks the field values on GetCatalogResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *GetCatalogResponse) Validate() error {
	if m == nil {
		return nil
	}

	return nil
}

// GetCatalogResponseValidationError is the validation error returned by
// GetCatalogResponse.Validate if the designated constraints aren't met.
type GetCatalogResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetCatalogResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetCatalogResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetCatalogResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetCatalogResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetCatalogResponseValidationError) ErrorName() string {
	return "GetCatalogResponseValidationError"
}

// Error satisfies the builtin error interface
func (e GetCatalogResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetCatalogResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetCatalogResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetCatalogResponseValidationError{}

// Validate checks the field values on ListCatalogsRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *ListCatalogsRequest) Validate() error {
	if m == nil {
		return nil
	}

	return nil
}

// ListCatalogsRequestValidationError is the validation error returned by
// ListCatalogsRequest.Validate if the designated constraints aren't met.
type ListCatalogsRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListCatalogsRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListCatalogsRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListCatalogsRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListCatalogsRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListCatalogsRequestValidationError) ErrorName() string {
	return "ListCatalogsRequestValidationError"
}

// Error satisfies the builtin error interface
func (e ListCatalogsRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListCatalogsRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListCatalogsRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListCatalogsRequestValidationError{}

// Validate checks the field values on ListCatalogsResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *ListCatalogsResponse) Validate() error {
	if m == nil {
		return nil
	}

	for idx, item := range m.GetCatalogs() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ListCatalogsResponseValidationError{
					field:  fmt.Sprintf("Catalogs[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	return nil
}

// ListCatalogsResponseValidationError is the validation error returned by
// ListCatalogsResponse.Validate if the designated constraints aren't met.
type ListCatalogsResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListCatalogsResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListCatalogsResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListCatalogsResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListCatalogsResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListCatalogsResponseValidationError) ErrorName() string {
	return "ListCatalogsResponseValidationError"
}

// Error satisfies the builtin error interface
func (e ListCatalogsResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListCatalogsResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListCatalogsResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListCatalogsResponseValidationError{}

// Validate checks the field values on DelCatalogRequest with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *DelCatalogRequest) Validate() error {
	if m == nil {
		return nil
	}

	return nil
}

// DelCatalogRequestValidationError is the validation error returned by
// DelCatalogRequest.Validate if the designated constraints aren't met.
type DelCatalogRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DelCatalogRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DelCatalogRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DelCatalogRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DelCatalogRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DelCatalogRequestValidationError) ErrorName() string {
	return "DelCatalogRequestValidationError"
}

// Error satisfies the builtin error interface
func (e DelCatalogRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDelCatalogRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DelCatalogRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DelCatalogRequestValidationError{}

// Validate checks the field values on DelCatalogResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *DelCatalogResponse) Validate() error {
	if m == nil {
		return nil
	}

	return nil
}

// DelCatalogResponseValidationError is the validation error returned by
// DelCatalogResponse.Validate if the designated constraints aren't met.
type DelCatalogResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DelCatalogResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DelCatalogResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DelCatalogResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DelCatalogResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DelCatalogResponseValidationError) ErrorName() string {
	return "DelCatalogResponseValidationError"
}

// Error satisfies the builtin error interface
func (e DelCatalogResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDelCatalogResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DelCatalogResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DelCatalogResponseValidationError{}
