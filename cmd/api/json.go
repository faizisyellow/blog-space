package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/go-playground/validator/v10"
)

type Envelope struct {
	Data  any `json:"data"`
	Error any `json:"error"`
}

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func WriteHttpJson(w http.ResponseWriter, data any, status int) error {

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(data)
}

func ReadHttpJson(w http.ResponseWriter, r *http.Request, data any) error {

	if len(r.Header["Content-Type"]) == 0 {
		return fmt.Errorf("no header content-type found")
	}

	if r.Header["Content-Type"][0] != "application/json" {
		return fmt.Errorf("header request only accept json")
	}

	// limit of the body size for 1mb
	maxBytes := 1_048_578

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decode := json.NewDecoder(r.Body)
	decode.DisallowUnknownFields()

	return decode.Decode(data)

}

func WriteHttpNoContent(w http.ResponseWriter) error {

	w.WriteHeader(http.StatusNoContent)

	_, err := fmt.Fprint(w)
	if err != nil {
		return err
	}

	return nil
}

func WriteErrorJson(w http.ResponseWriter, err error, status int) error {

	return WriteHttpJson(w, Envelope{Data: nil, Error: err.Error()}, status)

}

func (app *Application) JsonSuccessReponse(w http.ResponseWriter, r *http.Request, data any, status int) error {

	log.Info("Success", "Path", r.URL, "Method", r.Method, "Status", status)

	err := WriteHttpJson(w, Envelope{Data: data, Error: nil}, status)
	if err != nil {
		return fmt.Errorf("error sending json response: %v", err)
	}

	return nil

}

func (app *Application) JsonReponseNoContent(w http.ResponseWriter, r *http.Request) error {

	log.Info("Success", "Path", r.URL, "Method", r.Method, "Status", http.StatusNoContent)

	err := WriteHttpNoContent(w)
	if err != nil {
		return fmt.Errorf("error sending json response: %v", err)
	}

	return nil
}
