package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/charmbracelet/log"
)

var (
	ErrUnauthorize = errors.New("request unuthorize")
	ErrForbidden   = errors.New("request is forbidden")
)

func (app *Application) InternalServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {

	log.Error("Internal Server Error", "Path", r.URL, "Method", r.Method, "Status", http.StatusInternalServerError, "Message", err)

	if err := WriteErrorJson(w, fmt.Errorf("server encounter errors"), http.StatusInternalServerError); err != nil {

		log.Errorf("error sending json response: %v", err)

	}
}

func (app *Application) BadRequestErrorResponse(w http.ResponseWriter, r *http.Request, err error) {

	log.Error("Bad Request Error", "Path", r.URL, "Method", r.Method, "Status", http.StatusBadRequest, "Message", err)

	if err := WriteErrorJson(w, err, http.StatusBadRequest); err != nil {
		log.Errorf("error sending json response: %v", err)
	}

}

func (app *Application) NotFoundErrorResponse(w http.ResponseWriter, r *http.Request, err error) {

	log.Error("Not Found Error", "Path", r.URL, "Method", r.Method, "Status", http.StatusNotFound, "Message", err)

	if err := WriteErrorJson(w, err, http.StatusNotFound); err != nil {
		log.Errorf("error sending json response: %v", err)
	}
}

func (app *Application) ConflictErrorResponse(w http.ResponseWriter, r *http.Request, err error) {

	log.Error("Conflict Error:", "Path", r.URL, "Method", r.Method, "Status", http.StatusConflict, "Message", err)

	if err := WriteErrorJson(w, err, http.StatusConflict); err != nil {
		log.Errorf("error sending json response: %v", err)
	}
}

func (app *Application) UnAuthorizeErrorResponse(w http.ResponseWriter, r *http.Request, err error) {

	var errMsg strings.Builder
	errMsg.WriteString(ErrUnauthorize.Error())
	errMsg.WriteString(":")
	errMsg.WriteString(err.Error())

	log.Error("UnAuthorize Error:", "Path", r.URL, "Method", r.Method, "Status", http.StatusUnauthorized, "Message", errMsg.String())

	if err := WriteErrorJson(w, err, http.StatusUnauthorized); err != nil {
		log.Errorf("error sending json response: %v", err)
	}
}

func (app *Application) ForbiddenErrorResponse(w http.ResponseWriter, r *http.Request) {

	log.Error("Forbidden Error:", "Path", r.URL, "Method", r.Method, "Status", http.StatusForbidden, "Message", "request is forbidden")

	if err := WriteErrorJson(w, ErrForbidden, http.StatusForbidden); err != nil {
		log.Errorf("error sending json response: %v", err)
	}
}
