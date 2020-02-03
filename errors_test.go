package errors_test

import (
	"testing"

	errors "github.com/micharaze/errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func TestContext(t *testing.T) {

	err := errors.InvalidArgument.New("an_error")
	errWithContext := errors.AddErrorContext(err, "a_field", "the field is empty")

	expectedContext := map[string]string{"field": "a_field", "message": "the field is empty"}

	assert.Equal(t, errors.InvalidArgument, errors.GetType(errWithContext))
	assert.Equal(t, expectedContext, errors.GetErrorContext(errWithContext))
	assert.Equal(t, err.Error(), errWithContext.Error())
}

func TestContextInUnknownError(t *testing.T) {
	err := errors.New("a custom error")

	errWithContext := errors.AddErrorContext(err, "a_field", "the field is empty")

	expectedContext := map[string]string{"field": "a_field", "message": "the field is empty"}

	assert.Equal(t, errors.Unknown, errors.GetType(errWithContext))
	assert.Equal(t, expectedContext, errors.GetErrorContext(errWithContext))
	assert.Equal(t, err.Error(), errWithContext.Error())
}

func TestWrapf(t *testing.T) {
	err := errors.New("an_error")
	wrappedError := errors.InvalidArgument.Wrapf(err, "error %s", "1")

	assert.Equal(t, errors.InvalidArgument, errors.GetType(wrappedError))
	assert.EqualError(t, wrappedError, "error 1: an_error")
}

func TestWrapfInUnknownError(t *testing.T) {
	err := errors.Newf("an_error %s", "2")
	wrappedError := errors.Wrapf(err, "error %s", "1")

	assert.Equal(t, errors.Unknown, errors.GetType(wrappedError))
	assert.EqualError(t, wrappedError, "error 1: an_error 2")
}

func TestNewCodeError(t *testing.T) {
	err := errors.NewCode(codes.Unauthenticated, "error 1")
	err2 := errors.NewCodef(codes.PermissionDenied, "error 2")

	assert.Equal(t, errors.Unauthenticated, errors.GetType(err))
	assert.Equal(t, errors.PermissionDenied, errors.GetType(err2))
}
