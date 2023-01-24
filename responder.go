package responder

import (
	"fmt"
	"net/http"
	"strings"
)

// NEW return a new instance of the Respond object
func New(w http.ResponseWriter) *Respond {
	response := newHttpResponse(w, make(map[string]interface{}))
	return &Respond{w: w, response: response}
}

// Respond object struct
type Respond struct {
	w        http.ResponseWriter
	response *HttpResponse
}

// With allow you set flash message
func (res *Respond) With(name string, value interface{}) *HttpResponse {
	res.response.setAttributes(map[string]interface{}{
		name: value,
	})
	return res.response
}

// OK respond with http.StatusOK
func (res *Respond) OK(payload interface{}) {
	res.response.OK(payload)
}

// NoContent respond with http.StatusNoContent
func (res *Respond) NoContent() {
	res.response.NoContent()
}

// Created respond with http.StatusCreated
func (res *Respond) Created(r *http.Request, resource interface{}) {
	res.response.Created(r, resource)
}

// NotFound respond with http.StatusNotFound
func (res *Respond) NotFound(err error) {
	res.response.NotFound(err)
}

// Unauthorized respond with http.StatusUnauthorized
func (res *Respond) Unauthorized(err error) {
	res.response.Unauthorized(err)
}

// Forbidden respond with http.StatusForbidden
func (res *Respond) Forbidden(err error) {
	res.response.Forbidden(err)
}

// BadRequest respond with http.StatusBadRequest
func (res *Respond) BadRequest(err error) {
	res.response.BadRequest(err)
}

// NotFound respond with
func (res *Respond) UnprocessableEntity(err error) {
	res.response.UnprocessableEntity(err)
}

// Conflict respond with http.StatusConflict
func (res *Respond) Conflict(err error) {
	res.response.Conflict(err)
}

// InternalServerError respond with http.StatusInternalServerError
func (res *Respond) InternalServerError(err error) {
	res.response.InternalServerError(err)
}

// Error respond with an error based on the ErrorFormatter object
func (res *Respond) Error(err error) {
	res.response.Error(err)
}

// Plain stream a plain text file
func (res *Respond) Plain(stream []byte, fileName string) {
	res.response.Plain(stream, fileName)
}

// PDF stream an PDF file
func (res *Respond) PDF(stream []byte) {
	res.response.PDF(stream)
}

// Excel stream an Excel file
func (res *Respond) Excel(stream []byte) {
	res.response.Excel(stream)
}

func buildHost(r *http.Request) string {
	return fmt.Sprintf("%s://%s", r.URL.Scheme, r.Host)
}

type uriComponentsBuilder struct {
	request *http.Request
	path    string
}

func NewUriComponentsBuilder(r *http.Request) *uriComponentsBuilder {
	return &uriComponentsBuilder{request: r}
}

func (u *uriComponentsBuilder) Path(path string) {
	u.path = path
}

func (u *uriComponentsBuilder) ToURI() string {

	if !strings.HasPrefix(u.path, "/") {
		u.path = "/" + u.path
	}

	return fmt.Sprintf("%s%s", buildHost(u.request), u.path)
}
