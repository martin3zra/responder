package respond_test

import (
	"testing"

	"github.com/martin3zra/respond"
)

func TestErroFormatter(t *testing.T) {

	var instance interface{}
	instance = newNotFound()

	if _, ok := instance.(respond.ErrorFormatter); !ok {
		t.Errorf("handler returned wrong error: got %v want notFound", instance)
	}
}

func TestErrorDescriptor(t *testing.T) {
	var instance interface{}
	instance = newBadRequest()

	if _, ok := instance.(respond.ErrorFormatter); !ok {
		t.Errorf("handler returned wrong error: got %v want notFound", instance)
	}
}

func newNotFound() *notFound {
	return new(notFound)
}

type notFound struct {
}

func (notFound) Code() int {
	return 5
}

func (notFound) Error() string {
	return "resource not found"
}

func (notFound) Description() *string {
	description := "resource not found description"
	return &description
}

func (notFound) InfoURL() *string {
	info := "resource not found URL"
	return &info
}

func newBadRequest() *badRequest {
	return new(badRequest)
}

type badRequest struct {
	respond.ErrorDescriptor
}

func (badRequest) Code() int {
	return 3
}

func (badRequest) Error() string {
	return "bad Request"
}
