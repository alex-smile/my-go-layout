package errorx

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
)

const (
	systemSource = "mygo"
)

// Code references to the kind of error that was returned
type Code string

// StatusMessage in order to make it simple
// to keep a list of the used codes and the default message
type StatusMessage struct {
	Message    string
	StatusCode int
}

// NewErrorStatusMessage build a message with a status code
func NewErrorStatusMessage(message string, status int) StatusMessage {
	return StatusMessage{
		Message:    message,
		StatusCode: status,
	}
}

const (
	// ErrorCodeResourceNotFound error when the requested resource was not found
	ErrorCodeResourceNotFound Code = "resource_not_found"
	// ErrorCodeInvalidArgs error when the request data is incorrect or incomplete
	ErrorCodeInvalidArgs Code = "invalid_args"
	// ErrorCodeBadRequest error when the needed input is not provided
	ErrorCodeBadRequest Code = "bad_request"
	//ErrorCodeResourceAlreadyExists when a resource already exists
	ErrorCodeResourceAlreadyExists Code = "resource_already_exist"
	// ErrorCodeUnknownIssue when the issue is unknown
	ErrorCodeUnknownIssue Code = "unknown_issue"
	// ErrorCodePermissionDenied when user don't have necessary permissions
	ErrorCodePermissionDenied Code = "permission_denied"
	// ErrorCodeResourceStateConflict when the resource is in another state and generates a conflict
	ErrorCodeResourceStateConflict Code = "resource_state_conflict"
	// ErrorCodeConflict conflict
	ErrorCodeConflict Code = "conflict"
	// ErrorCodeNotImplemented when some method is not implemented
	ErrorCodeNotImplemented Code = "not_implemented"
	// ErrorCodeUnauthorized when user is not authorized
	ErrorCodeUnauthorized Code = "unauthorized"
	// ErrorCodeNotFound when is not found
	ErrorCodeNotFound Code = "not_found"
	// ErrorCodeDatabaseError when using database commands and it returned an error
	ErrorCodeDatabaseError Code = "database_error"
	// ErrorCodeInternalServerError internal server error
	ErrorCodeInternalServerError Code = "internal_server_error"
	// ErrorCodeRequestBackendError request backend service error
	ErrorCodeRequestBackendError Code = "request_backend_error"
)

var (
	// ErrorMessageList general default messages for all error codes
	ErrorMessageList = map[Code]StatusMessage{
		ErrorCodeResourceNotFound:      NewErrorStatusMessage("resource not found", http.StatusNotFound),
		ErrorCodeInvalidArgs:           NewErrorStatusMessage("parameters were invalid", http.StatusBadRequest),
		ErrorCodeBadRequest:            NewErrorStatusMessage("bad request", http.StatusBadRequest),
		ErrorCodeResourceAlreadyExists: NewErrorStatusMessage("the posted resource already existed.", http.StatusBadRequest),
		ErrorCodeUnknownIssue:          NewErrorStatusMessage("unknown issue was caught and message was not specified.", http.StatusInternalServerError),
		ErrorCodePermissionDenied:      NewErrorStatusMessage("current user has no permission to perform the action.", http.StatusForbidden),
		ErrorCodeResourceStateConflict: NewErrorStatusMessage("the posted resource already existed.", http.StatusConflict),
		ErrorCodeNotImplemented:        NewErrorStatusMessage("method not implemented", http.StatusNotImplemented),
		ErrorCodeUnauthorized:          NewErrorStatusMessage("unauthorized", http.StatusUnauthorized),
		ErrorCodeNotFound:              NewErrorStatusMessage("not found", http.StatusNotFound),
		ErrorCodeDatabaseError:         NewErrorStatusMessage("database error.", http.StatusInternalServerError),
		ErrorCodeInternalServerError:   NewErrorStatusMessage("internal server error", http.StatusInternalServerError),
		ErrorCodeRequestBackendError:   NewErrorStatusMessage("request backend service error", http.StatusInternalServerError),
	}
)

// AddError add error
func AddError(code Code, message string, status int) {
	ErrorMessageList[code] = NewErrorStatusMessage(message, status)
}

// MygoError common error struct used in the whole application
type MygoError struct {
	Source     string `json:"source"`
	Message    string `json:"message"`
	Code       Code   `json:"code"`
	StatusCode int    `json:"-"`
}

// New Constructor function for the error structure
func New(code Code) *MygoError {
	var (
		message = "Error not described"
		status  = http.StatusInternalServerError
	)

	if val, ok := ErrorMessageList[code]; ok {
		message = val.Message
		status = val.StatusCode
	}

	return &MygoError{
		Source:     systemSource,
		Message:    message,
		Code:       code,
		StatusCode: status,
	}
}

func NewWithMessage(code Code, message string, args ...interface{}) *MygoError {
	return New(code).SetMessage(message, args...)
}

func NewWithError(err error) *MygoError {
	if x, ok := err.(*MygoError); ok {
		return x
	}
	return NewWithMessage(ErrorCodeInternalServerError, err.Error())
}

// Error satisfies the error interface
func (e *MygoError) Error() string {
	return fmt.Sprintf("%v: %s", e.Code, e.Message)
}

// SetMessage sets a message using the given format and parameters returns itself
func (e *MygoError) SetMessage(format string, args ...interface{}) *MygoError {
	e.Message = fmt.Sprintf(format, args...)
	return e
}

// SetCode sets a code and returns itself for chaining calls
func (e *MygoError) SetCode(code Code) *MygoError {
	e.Code = code
	return e
}

func (e *MygoError) SetSource(source string) *MygoError {
	e.Source = source
	return e
}

type Fields map[string]interface{}

func (e *MygoError) WithField(key string, value interface{}) *MygoError {
	e.WithFields(Fields{key: value})
	return e
}

func (e *MygoError) WithFields(fields Fields) *MygoError {
	if len(fields) == 0 {
		return e
	}

	fieldsMessage := e.formatFields(fields)
	if e.Message == "" {
		e.Message = fieldsMessage
	} else {
		e.Message = fmt.Sprintf("%s [%s]", e.Message, fieldsMessage)
	}
	return e
}

func (e *MygoError) formatFields(fields Fields) string {
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s=\"%v\"", key, fields[key]))
	}
	return strings.Join(parts, " ")
}
