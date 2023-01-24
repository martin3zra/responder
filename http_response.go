package responder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

func newHttpResponse(w http.ResponseWriter, attributes map[string]interface{}) *HttpResponse {
	return &HttpResponse{writer: w, attributes: attributes}
}

type HttpResponse struct {
	writer     http.ResponseWriter
	attributes map[string]interface{}
}

func (response *HttpResponse) setAttributes(attributes map[string]interface{}) *HttpResponse {
	response.attributes = attributes
	return response
}

// emptyStatus a collection of http status code doesn't need a body as response
func (response *HttpResponse) emptyStatus() []int {
	return []int{
		http.StatusCreated,
		http.StatusNoContent,
		http.StatusResetContent,
	}
}

// OK respond with http.StatusOK
func (response *HttpResponse) OK(payload map[string]interface{}) {
	data := map[string]interface{}{
		"data":  payload,
		"flash": response.attributes,
	}

	res, err := json.Marshal(data)
	if err != nil {
		response.writer.WriteHeader(http.StatusInternalServerError)
		response.writer.Write([]byte(err.Error()))
		return
	}

	response.asJSON(http.StatusOK, []byte(res))
}

// NoContent ...
func (response *HttpResponse) NoContent() {
	response.asJSON(http.StatusNoContent, nil)
}

// Created ...
func (response *HttpResponse) Created(r *http.Request, resource interface{}) {
	response.writer.Header().Set("Location", response.buildLocationURL(r, resource))
	response.asJSON(http.StatusCreated, nil)
}

// Plain ...
func (response *HttpResponse) Plain(stream []byte, fileName string) {
	response.writer.Header().Set("content-disposition", "attachment; filename=\""+fileName+"\"")
	response.file(stream, "application/plain")
}

// PDF ...
func (response *HttpResponse) PDF(stream []byte) {
	response.file(stream, "application/pdf")
}

// Excel ...
func (response *HttpResponse) Excel(stream []byte) {
	response.writer.WriteHeader(http.StatusOK)
	// stream straight to client(browser)
	response.writer.Header().Set("Content-Description", "File Transfer")
	response.writer.Header().Set("Content-Disposition", "attachment;")
	response.writer.Header().Set("Content-type", "application/octet-stream")
	b := bytes.NewBuffer(stream)

	if _, err := b.WriteTo(response.writer); err != nil {
		response.writer.WriteHeader(http.StatusInternalServerError)
		response.writer.Write([]byte(err.Error()))
		return
	}

	response.writer.Write([]byte(""))
}

// NotFound is returned when the resource requested by your application does not exist
func (response *HttpResponse) NotFound(err error) {
	response.asJSON(http.StatusNotFound, response.getMessage(err))
}

// Unauthorized is returned when there is a problem with the credentials provided by your application.
// This code indicates that your application tried to operate on a protected resource without
// providing the proper authorization. It may have provided the wrong credentials or none at all
func (response *HttpResponse) Unauthorized(err error) {
	response.asJSON(http.StatusUnauthorized, response.getMessage(err))
}

// Forbidden is returned when your application is not authorized to access the requested resource,
// or when your application is being rate limited
func (response *HttpResponse) Forbidden(err error) {
	response.asJSON(http.StatusForbidden, response.getMessage(err))
}

// BadRequest is returned when the request entity sent by your application could not
// be understood by the server due to malformed syntax (e.g. invalid payload, data type mismatch)
func (response *HttpResponse) BadRequest(err error) {
	response.asJSON(http.StatusBadRequest, response.getMessage(err))
}

// UnprocessableEntity ...
func (response *HttpResponse) UnprocessableEntity(err error) {
	response.asJSON(http.StatusUnprocessableEntity, response.getMessage(err))
}

// Conflict is returned when the request sent by your application could not be completed due to a conflict
// with the current state of the resource
func (response *HttpResponse) Conflict(err error) {
	response.asJSON(http.StatusConflict, response.getMessage(err))
}

// Error is returned when the server encountered an unexpected condition which prevented it from fulfilling
// the request sent by your application
func (response *HttpResponse) InternalServerError(err error) {
	response.writer.Header().Set("Content-Type", "application/json")
	response.writer.WriteHeader(http.StatusInternalServerError)
	if err == nil {
		return
	}
	response.writer.Write([]byte(err.Error()))
}

// Error is returned when the server encountered an unexpected condition which prevented it from fulfilling
// the request sent by your application
func (response *HttpResponse) Error(err error) {

	errValue, ok := err.(ErrorFormatter)
	if !ok {
		response.InternalServerError(err)
		return
	}

	response.composeCustomError(errValue)
}

func (response *HttpResponse) asJSON(statusCode int, stream []byte) {
	response.writer.Header().Set("Content-Type", "application/json")
	response.writer.WriteHeader(statusCode)

	if response.doesNotRequireContent(statusCode) {
		return
	}

	if stream == nil {
		return
	}

	response.writer.Write(stream)
}

func (response *HttpResponse) file(stream []byte, contentType string) {
	response.writer.WriteHeader(http.StatusOK)
	// stream straight to client(browser)
	response.writer.Header().Set("Content-type", contentType)
	b := bytes.NewBuffer(stream)

	if _, err := b.WriteTo(response.writer); err != nil {
		response.writer.WriteHeader(http.StatusInternalServerError)
		response.writer.Write([]byte(err.Error()))
		return
	}
}

func (response *HttpResponse) doesNotRequireContent(statusCode int) bool {
	return response.in(response.emptyStatus(), statusCode)
}

// in checks if a given value exists in an array
func (response *HttpResponse) in(items, item interface{}) bool {
	arr := reflect.ValueOf(items)

	if arr.Kind() != reflect.Array && arr.Kind() != reflect.Slice {
		fmt.Println(arr.Kind())
		panic("Invalid data-type: array or slice expected")
	}

	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

func (response *HttpResponse) getMessage(err error) []byte {
	if err == nil {
		return nil
	}

	if value, ok := err.(ErrorFormatter); ok {

		data := map[string]interface{}{"code": value.Code(), "message": value.Error()}
		if value.Description() != nil {
			data["description"] = value.Description()
		}

		if value.InfoURL() != nil {
			data["info_url"] = value.InfoURL()
		}

		response, err := json.Marshal(data)
		if err != nil {
			return []byte(err.Error())
		}

		return response
	}

	return []byte(err.Error())
}

// func (response *HttpResponse) getBytes(key interface{}) ([]byte, error) {
// 	var buf bytes.Buffer
// 	enc := gob.NewEncoder(&buf)
// 	err := enc.Encode(key)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return buf.Bytes(), nil
// }

func (response *HttpResponse) composeCustomError(err ErrorFormatter) {
	switch err.Status() {
	case http.StatusUnauthorized:
		response.Unauthorized(err)
	case http.StatusForbidden:
		response.Forbidden(err)
	case http.StatusConflict:
		response.Conflict(err)
	case http.StatusUnprocessableEntity:
		response.UnprocessableEntity(err)
	case http.StatusNotFound:
		response.NotFound(err)
	case http.StatusBadRequest:
		response.BadRequest(err)
	case http.StatusInternalServerError:
		break
	default:
		response.InternalServerError(err)
	}
}

func (response *HttpResponse) buildLocationURL(r *http.Request, resource interface{}) string {
	protocol := "http"
	if r.TLS != nil {
		protocol = "https"
	}

	return fmt.Sprintf("%s://%s%s/%v",
		protocol,
		r.Host,
		r.URL.RequestURI(),
		resource,
	)
}
