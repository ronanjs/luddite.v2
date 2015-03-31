package luddite

import (
	"encoding/xml"
	"fmt"
)

const (
	// Common error codes
	EcodeUnknown               = 0
	EcodeInternal              = 1
	EcodeUnsupportedMediaType  = 2
	EcodeSerializationFailed   = 3
	EcodeDeserializationFailed = 4
	EcodeIdentifierMismatch    = 5

	// Service-specific error codes
	EcodeServiceBase = 1024
)

var commonErrorMessages = map[int]string{
	EcodeUnknown:               "Unknown error",
	EcodeInternal:              "Internal error",
	EcodeUnsupportedMediaType:  "Unsupported media type: %s",
	EcodeSerializationFailed:   "Serialization failed: %s",
	EcodeDeserializationFailed: "Deserialization failed: %s",
	EcodeIdentifierMismatch:    "Resource identifier in URL doesn't match value in body",
}

// Error is a structured error that is returned as the body in all 4xx and 5xx responses.
type Error struct {
	XMLName xml.Name `json:"-" xml:"error"`
	Code    int      `json:"code" xml:"code"`
	Message string   `json:"message" xml:"message"`
	Stack   string   `json:"stack,omitempty" xml:"stack,omitempty"`
}

func (e *Error) Error() string {
	return e.Message
}

// NewError allocates and initializes an Error. If a non-nil
// errorMessages map is passed, the error message string is resolved
// using this map. Otherwise a map containing common error message
// strings is used. Services or resources should generally use error
// code values greater than or equal to EcodeServiceBase. Error codes
// below this value are reserved for common use.
func NewError(errorMessages map[int]string, code int, args ...interface{}) *Error {
	var (
		format string
		ok     bool
	)

	// Lookup an error message string by error code: first try the
	// caller provided error message map with fallback to the
	// common error message map.
	if errorMessages != nil {
		format, ok = errorMessages[code]
	}

	if !ok {
		format, ok = commonErrorMessages[code]
	}

	// If no error message could be found, failsafe by using a
	// known-good common error message along with the caller's
	// error code.
	if !ok {
		format = commonErrorMessages[EcodeUnknown]
		args = nil
	}

	// Optionally format the error message
	var message string
	if len(args) != 0 {
		message = fmt.Sprintf(format, args...)
	} else {
		message = format
	}

	return &Error{
		Code:    code,
		Message: message,
	}
}