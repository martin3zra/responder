package responder_test

import (
	"net/http"
	"testing"

	"github.com/martin3zra/responder"
)

func TestErrorFormatter(t *testing.T) {

	var instance interface{} = new(notFound)

	if _, ok := instance.(responder.ErrorFormatter); !ok {
		t.Errorf("handler returned wrong error: got %v want notFound", instance)
	}
}

func TestErrorDescriptor(t *testing.T) {
	var instance interface{} = new(badRequest)

	t.Run("implement ErrorFormatter", func(t *testing.T) {

		if _, ok := instance.(responder.ErrorFormatter); !ok {
			t.Errorf("handler returned wrong error: got %v want notFound", instance)
		}
	})

	t.Run("it returns bad request when HTTP Status code is 400", func(t *testing.T) {

		if instance.(responder.ErrorFormatter).Status() != http.StatusBadRequest {
			t.Errorf("handler returned wrong error: got %d want %d", instance.(responder.ErrorFormatter).Status(), http.StatusBadRequest)
		}
	})
}

type notFound struct {
}

func (notFound) Status() int {
	return http.StatusNotFound
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

type badRequest struct {
	responder.ErrorDescriptor
}

func (badRequest) Status() int {
	return http.StatusBadRequest
}

func (badRequest) Code() int {
	return 3
}

func (badRequest) Error() string {
	return "bad Request"
}
