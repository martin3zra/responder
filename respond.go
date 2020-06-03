package respond

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

//emptyStatus a collection of http status code doesn't need a body as response
var emptyStatus = []int{
	http.StatusCreated,
	http.StatusNoContent,
	http.StatusResetContent,
}

//OK respond with http.StatusOK
func OK(w http.ResponseWriter, payload interface{}) {

	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	asJSON(w, http.StatusOK, []byte(response))
}

//NoContent ...
func NoContent(w http.ResponseWriter) {
	asJSON(w, http.StatusNoContent, nil)
}

// Created ...
func Created(w http.ResponseWriter, r *http.Request, resource interface{}) {
	w.Header().Set("Location", fmt.Sprintf("%s%s/%v", r.Host, r.URL.RequestURI(), resource))
	asJSON(w, http.StatusCreated, nil)
}

//Plain ...
func Plain(w http.ResponseWriter, stream []byte, fileName string) {
	w.Header().Set("content-disposition", "attachment; filename=\""+fileName+"\"")
	file(w, stream, "application/plain")
}

//PDF ...
func PDF(w http.ResponseWriter, stream []byte) {
	file(w, stream, "application/pdf")
}

//Excel ...
func Excel(w http.ResponseWriter, stream []byte) {
	w.WriteHeader(http.StatusOK)
	// stream straight to client(browser)
	w.Header().Set("Content-Description", "File Transfer")
	w.Header().Set("Content-Disposition", "attachment;")
	w.Header().Set("Content-type", "application/octet-stream")
	b := bytes.NewBuffer(stream)

	if _, err := b.WriteTo(w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(""))
}

//NotFound is returned when the resource requested by your application does not exist
func NotFound(w http.ResponseWriter, err error) {
	asJSON(w, http.StatusNotFound, getMessage(err))
}

//Unauthorized is returned when there is a problem with the credentials provided by your application.
//This code indicates that your application tried to operate on a protected resource without
//providing the proper authorization. It may have provided the wrong credentials or none at all
func Unauthorized(w http.ResponseWriter, err error) {
	asJSON(w, http.StatusUnauthorized, getMessage(err))
}

//Forbidden is returned when your application is not authorized to access the requested resource,
//or when your application is being rate limited
func Forbidden(w http.ResponseWriter, err error) {
	asJSON(w, http.StatusForbidden, getMessage(err))
}

//BadRequest is returned when the request entity sent by your application could not
//be understood by the server due to malformed syntax (e.g. invalid payload, data type mismatch)
func BadRequest(w http.ResponseWriter, err error) {
	asJSON(w, http.StatusBadRequest, getMessage(err))
}

//UnprocessableEntity ...
func UnprocessableEntity(w http.ResponseWriter, err error) {
	asJSON(w, http.StatusUnprocessableEntity, getMessage(err))
}

//Conflict is returned when the request sent by your application could not be completed due to a conflict
//with the current state of the resource
func Conflict(w http.ResponseWriter, err error) {
	asJSON(w, http.StatusConflict, getMessage(err))
}

//Error is returned when the server encountered an unexpected condition which prevented it from fulfilling
//the request sent by your application
func Error(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	if err == nil {
		return
	}
	w.Write([]byte(err.Error()))
}

func asJSON(w http.ResponseWriter, statusCode int, stream interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if doesRequireContent(statusCode) {
		return
	}

	if stream == nil {
		return
	}

	bytes, err := getBytes(stream)
	if err != nil {
		return
	}

	w.Write(bytes)
}

func file(w http.ResponseWriter, stream []byte, contentType string) {
	w.WriteHeader(http.StatusOK)
	// stream straight to client(browser)
	w.Header().Set("Content-type", contentType)
	b := bytes.NewBuffer(stream)

	if _, err := b.WriteTo(w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func doesRequireContent(statusCode int) bool {
	return !in(emptyStatus, statusCode)
}

// in checks if a given value exists in an array
func in(items, item interface{}) bool {
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

func getMessage(err error) *string {
	if err == nil {
		return nil
	}

	message := err.Error()
	return &message
}

func getBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
