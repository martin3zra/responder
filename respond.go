package respond

import (
	"bytes"
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

//NotFound ...
func NotFound(w http.ResponseWriter, err error) {
	asJSON(w, http.StatusNotFound, nil)
}

//Unauthorized ...
func Unauthorized(w http.ResponseWriter, err error) {
	asJSON(w, http.StatusUnauthorized, nil)
}

//Forbidden ...
func Forbidden(w http.ResponseWriter, err error) {
	asJSON(w, http.StatusForbidden, nil)
}

//BadRequest ...
func BadRequest(w http.ResponseWriter, err error) {
	asJSON(w, http.StatusBadRequest, nil)
}

//UnprocessableEntity ...
func UnprocessableEntity(w http.ResponseWriter, err error) {
	asJSON(w, http.StatusUnprocessableEntity, nil)
}

//Conflict ...
func Conflict(w http.ResponseWriter, err error) {
	asJSON(w, http.StatusConflict, nil)
}

//Error ...
func Error(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	if err == nil {
		return
	}
	w.Write([]byte(err.Error()))
}

func asJSON(w http.ResponseWriter, statusCode int, stream []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if doesRequireContent(statusCode) {
		return
	}
	w.Write(stream)
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
