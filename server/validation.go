package server

import (
	"errors"
	"fmt"
	"net/mail"
	"net/url"
	"strings"
)

// ValidationErrors represents a set of errors detected during validation.
// It satisfies the error interface.
type ValidationErrors []error

// Append adds a new validation error to the set.
func (ves *ValidationErrors) Append(fieldName, errorMsg string) {
	err := errors.New(fmt.Sprintf("%s: %s", fieldName, errorMsg))
	*ves = append(*ves, err)
}

// Error concatenates all errors in the set into a single error message.
func (ves ValidationErrors) Error() string {
	var sb strings.Builder
	for _, err := range ves {
		sb.WriteString(err.Error())
		sb.WriteRune('\n')
	}
	return sb.String()
}

// ValidateStringNonEmpty checks whether a string is non-empty.
func ValidateStringNonEmpty(errs *ValidationErrors, fieldName string, value string) {
	if len(value) == 0 {
		errs.Append(fieldName, "Missing required field")
	}
}

// ValidateURL checks whether a string is a valid URL.
func ValidateURL(errs *ValidationErrors, fieldName string, value string) {
	if value == "" {
		errs.Append(fieldName, "URL cannot be empty")
		return
	}

	_, err := url.Parse(value)
	if err != nil {
		errs.Append(fieldName, "Invalid URL")
	}
}

// ValidateEmailAddress checks whether a string is a valid email address (RFC 5322).
func ValidateEmailAddress(errs *ValidationErrors, fieldName string, value string) {
	if value == "" {
		errs.Append(fieldName, "Email address cannot be empty")
		return
	}

	_, err := mail.ParseAddress(value)
	if err != nil {
		errs.Append(fieldName, "Invalid email address")
	}
}
