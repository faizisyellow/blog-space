package main

import "net/http"

func (app *Application) PingHandler(w http.ResponseWriter, r *http.Request) {

	if err := app.JsonSuccessReponse(w, r, "ping", http.StatusOK); err != nil {
		app.InternalServerErrorResponse(w, r, err)
		return
	}
}
