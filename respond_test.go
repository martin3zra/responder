package respond_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/martin3zra/respond"
)

func TestMain(m *testing.M) {
	exitValue := m.Run()
	os.Exit(exitValue)
}

func TestResponses(t *testing.T) {
	cases := []struct {
		handler    http.HandlerFunc
		expectCode func(t *testing.T, w *httptest.ResponseRecorder)
		name       string
	}{
		{
			handler:    func(w http.ResponseWriter, r *http.Request) { respond.OK(w, nil) },
			expectCode: assertOK,
			name:       "it returns http status 200 when respond OK",
		},
		{
			handler:    func(w http.ResponseWriter, r *http.Request) {
				respond.Created(w,  fmt.Sprintf("%s://%s", r.URL.Scheme, r.Host))
			},
			expectCode: assertCreated,
			name:       "it returns http status 201 when respond created",
		},
		{
			handler:    func(w http.ResponseWriter, r *http.Request) { respond.NoContent(w) },
			expectCode: assertNoContent,
			name:       "it returns http status 204 when respond no content",
		},
		{
			handler:    func(w http.ResponseWriter, r *http.Request) { respond.BadRequest(w, nil) },
			expectCode: assertBadRequest,
			name:       "it returns http status 400 when respond bad request",
		},
		{
			handler:    func(w http.ResponseWriter, r *http.Request) { respond.Conflict(w, nil) },
			expectCode: assertConflict,
			name:       "it returns http status 409 when respond conflict",
		},
		{
			handler:    func(w http.ResponseWriter, r *http.Request) { respond.NotFound(w, nil) },
			expectCode: assertNotFound,
			name:       "it returns http status 404 when respond not found",
		},
		{
			handler:    func(w http.ResponseWriter, r *http.Request) { respond.UnprocessableEntity(w, nil) },
			expectCode: assertUnprocessableEntity,
			name:       "it returns http status 422 when respond unprocessable entity",
		},
		{
			handler:    func(w http.ResponseWriter, r *http.Request) { respond.Error(w, errors.New("some error")) },
			expectCode: assertInternalError,
			name:       "it returns http status 500 when respond error",
		},
	}

	for _, item := range cases {
		t.Run(item.name, func(t *testing.T) {

			rr := httptest.NewRecorder()
			item.handler(rr, buildRequest(t))
			item.expectCode(t, rr)
			assertIsJSON(t, rr)
		})
	}
}

func TestBadRequestResponse_ErrorFormatter(t *testing.T) {

	handler := func(w http.ResponseWriter, r *http.Request) { respond.BadRequest(w, newFormatterBadRequest()) }

	rr := httptest.NewRecorder()
	handler(rr, buildRequest(t))
	assertBadRequest(t, rr)
	assertIsJSON(t, rr)

	responseMap := transform(t, rr)

	if code, ok := responseMap["code"]; !ok {
		t.Errorf("expected key `code`: got nil")
	} else {
		if code.(float64) != 3 {
			t.Errorf("handler returned wrong code: got %v want %v", code, 3)
		}
	}

}

func TestResponseContentType(t *testing.T) {

	cases := []struct {
		handler           http.HandlerFunc
		expectContentType string
		name              string
	}{
		{
			handler:           func(w http.ResponseWriter, r *http.Request) { respond.PDF(w, nil) },
			expectContentType: "application/pdf",
			name:              "it returns application/pdf when respond PDF",
		},
		{
			handler:           func(w http.ResponseWriter, r *http.Request) { respond.Plain(w, nil, "name") },
			expectContentType: "application/plain",
			name:              "it returns application/plain when respond Plain",
		},
		{
			handler:           func(w http.ResponseWriter, r *http.Request) { respond.Excel(w, nil) },
			expectContentType: "application/octet-stream",
			name:              "it returns application/octet-stream when respond Excel",
		},
	}

	for _, item := range cases {
		t.Run(item.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			item.handler(rr, buildRequest(t))

			if rr.Header().Get("Content-Type") != item.expectContentType {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Header().Get("Content-Type"), item.expectContentType)
			}

			assertOK(t, rr)
		})
	}
}

func buildRequest(t *testing.T) *http.Request {
	req, err := http.NewRequest("GET", "http://localhost/ok", nil)
	if err != nil {
		t.Fatal(err)
	}

	return req
}

func assertOK(t *testing.T, w *httptest.ResponseRecorder) {
	assertStatusCode(t, http.StatusOK, w.Code)
}

func assertCreated(t *testing.T, w *httptest.ResponseRecorder) {
	assertStatusCode(t, http.StatusCreated, w.Code)
	_, err := url.ParseRequestURI(w.Header().Get("location"))
	if err != nil {
		t.Fatalf("URL on Location header is not valid.")
	}
}

func assertNoContent(t *testing.T, w *httptest.ResponseRecorder) {
	assertStatusCode(t, http.StatusNoContent, w.Code)
}

func assertConflict(t *testing.T, w *httptest.ResponseRecorder) {
	assertStatusCode(t, http.StatusConflict, w.Code)
}

func assertBadRequest(t *testing.T, w *httptest.ResponseRecorder) {
	assertStatusCode(t, http.StatusBadRequest, w.Code)
}

func assertUnprocessableEntity(t *testing.T, w *httptest.ResponseRecorder) {
	assertStatusCode(t, http.StatusUnprocessableEntity, w.Code)
}

func assertNotFound(t *testing.T, w *httptest.ResponseRecorder) {
	assertStatusCode(t, http.StatusNotFound, w.Code)
}

func assertInternalError(t *testing.T, w *httptest.ResponseRecorder) {
	assertStatusCode(t, http.StatusInternalServerError, w.Code)
}

func assertStatusCode(t *testing.T, expected, given int) {
	if given != expected {
		t.Errorf("handler returned wrong status code: got %v want %v",
			given, expected)
	}
}

func assertIsJSON(t *testing.T, w http.ResponseWriter) {
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("handler returned wrong status code: got %v want %v",
			w.Header().Get("Content-Type"), "application/json")
	}
}

func transform(t *testing.T, rr *httptest.ResponseRecorder) map[string]interface{} {
	responseMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(rr.Body.String()), &responseMap)
	if err != nil {
		t.Errorf("Cannot convert to json: %v", err)
	}

	return responseMap
}

func newFormatterBadRequest() *formatterBadRequest {
	return new(formatterBadRequest)
}

type formatterBadRequest struct {
	respond.ErrorDescriptor
}

func (formatterBadRequest) Code() int {
	return 3
}

func (formatterBadRequest) Error() string {
	return "bad Request"
}

func (formatterBadRequest) Description() *string {
	val := "some description here"
	return &val
}
