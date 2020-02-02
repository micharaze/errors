package errors

import (
	"fmt"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
)

// ErrorType is the type of an error
type ErrorType uint32

const (
	// InvalidArgument error
	InvalidArgument ErrorType = ErrorType(codes.InvalidArgument)
	// FailedPrecondition error
	FailedPrecondition ErrorType = ErrorType(codes.FailedPrecondition)
	// OutOfRange error
	OutOfRange ErrorType = ErrorType(codes.OutOfRange)
	// Unauthenticated error
	Unauthenticated ErrorType = ErrorType(codes.Unauthenticated)
	// PermissionDenied error
	PermissionDenied ErrorType = ErrorType(codes.PermissionDenied)
	// NotFound error
	NotFound ErrorType = ErrorType(codes.NotFound)
	// Aborted error
	Aborted ErrorType = ErrorType(codes.Aborted)
	// AlreadyExists error
	AlreadyExists ErrorType = ErrorType(codes.AlreadyExists)
	// ResourceExhausted error
	ResourceExhausted ErrorType = ErrorType(codes.ResourceExhausted)
	// Canceled error
	Canceled ErrorType = ErrorType(codes.Canceled)
	// DataLoss error
	DataLoss ErrorType = ErrorType(codes.DataLoss)
	// Unknown error
	Unknown ErrorType = ErrorType(codes.Unknown)
	// Internal error
	Internal ErrorType = ErrorType(codes.Internal)
	// Unimplemented error
	Unimplemented ErrorType = ErrorType(codes.Unimplemented)
	// Unavailable error
	Unavailable ErrorType = ErrorType(codes.Unavailable)
	// DeadlineExceeded error
	DeadlineExceeded ErrorType = ErrorType(codes.DeadlineExceeded)
)

// Map for converting ErrorType to HTTP code
var httpMap = map[ErrorType]uint32{
	InvalidArgument:    400,
	FailedPrecondition: 400,
	OutOfRange:         400,
	Unauthenticated:    401,
	PermissionDenied:   403,
	NotFound:           404,
	Aborted:            409,
	AlreadyExists:      409,
	ResourceExhausted:  429,
	Canceled:           499,
	DataLoss:           500,
	Unknown:            500,
	Internal:           500,
	Unimplemented:      501,
	Unavailable:        503,
	DeadlineExceeded:   504,
}

type customError struct {
	errorType     ErrorType
	originalError error
	context       errorContext
}

type errorContext struct {
	Field   string
	Message string
}

// New creates a new customError
func (errorType ErrorType) New(msg string) error {
	return customError{errorType: errorType, originalError: errors.New(msg)}
}

// Newf creates a new customError with formatted message
func (errorType ErrorType) Newf(msg string, args ...interface{}) error {
	return customError{errorType: errorType, originalError: fmt.Errorf(msg, args...)}
}

// Wrap creates a new wrapped error
func (errorType ErrorType) Wrap(err error, msg string) error {
	return errorType.Wrapf(err, msg)
}

// Wrapf creates a new wrapped error with formatted message
func (errorType ErrorType) Wrapf(err error, msg string, args ...interface{}) error {
	return customError{errorType: errorType, originalError: errors.Wrapf(err, msg, args...)}
}

// Code converts ErrorType to gRPC Code
func (errorType ErrorType) Code() codes.Code {
	return codes.Code(errorType)
}

// CodeString converts ErrorType to gRPC Code string
func (errorType ErrorType) CodeString() string {
	return codes.Code(errorType).String()
}

// HTTP converts ErrorType to HTTP error code
func (errorType ErrorType) HTTP() uint32 {
	return httpMap[errorType]
}

// Error returns the mssage of a customError
func (error customError) Error() string {
	return error.originalError.Error()
}

// New creates a no type error
func New(msg string) error {
	return customError{errorType: Unknown, originalError: errors.New(msg)}
}

// Newf creates a no type error with formatted message
func Newf(msg string, args ...interface{}) error {
	return customError{errorType: Unknown, originalError: errors.New(fmt.Sprintf(msg, args...))}
}

// NewCode creates an error by given gRPC Code and message.
func NewCode(code codes.Code, msg string) error {
	return customError{errorType: ErrorType(code), originalError: errors.New(msg)}
}

// NewCodef creates an error by given gRPC Code and formatted message.
func NewCodef(code codes.Code, msg string, args ...interface{}) error {
	return customError{errorType: ErrorType(code), originalError: errors.New(fmt.Sprintf(msg, args...))}
}

// Wrap an error with a string
func Wrap(err error, msg string) error {
	return Wrapf(err, msg)
}

// Cause gives the original error
func Cause(err error) error {
	return errors.Cause(err)
}

// Wrapf an error with format string
func Wrapf(err error, msg string, args ...interface{}) error {
	wrappedError := errors.Wrapf(err, msg, args...)
	if customErr, ok := err.(customError); ok {
		return customError{
			errorType:     customErr.errorType,
			originalError: wrappedError,
			context:       customErr.context,
		}
	}

	return customError{errorType: Unknown, originalError: wrappedError}
}

// AddErrorContext adds a context to an error
func AddErrorContext(err error, field, message string) error {
	context := errorContext{Field: field, Message: message}
	if customErr, ok := err.(customError); ok {
		return customError{errorType: customErr.errorType, originalError: customErr.originalError, context: context}
	}

	return customError{errorType: Unknown, originalError: err, context: context}
}

// GetErrorContext returns the error context
func GetErrorContext(err error) map[string]string {
	emptyContext := errorContext{}
	if customErr, ok := err.(customError); ok || customErr.context != emptyContext {

		return map[string]string{"field": customErr.context.Field, "message": customErr.context.Message}
	}

	return nil
}

// GetType returns the error type
func GetType(err error) ErrorType {
	if customErr, ok := err.(customError); ok {
		return customErr.errorType
	}

	return Unknown
}
