package main

import (
	"net/http"

	"faissal.com/blogSpace/internal/repository"
	"faissal.com/blogSpace/internal/utils"
)

// @Summary		Get User Profile
// @Description	Get User Profile Who's log in
// @Tags			Users
// @Accept			json
// @Produce		json
// @Security		JWT
// @Success		200	{object}	main.Envelope{data=repository.User,error=nil}
// @Failure		400	{object}	main.Envelope{data=nil,error=string}
// @Failure		401	{object}	main.Envelope{data=nil,error=string}
// @Failure		500	{object}	main.Envelope{data=nil,error=string}
// @Router			/users/profile [get]
func (app *Application) GetUserProfileHandler(w http.ResponseWriter, r *http.Request) {

	user, err := utils.GetContentFromContext[repository.User](r, UsrCtx)
	if err != nil {
		app.InternalServerErrorResponse(w, r, err)
		return
	}

	if err := app.JsonSuccessReponse(w, r, user, http.StatusOK); err != nil {
		app.InternalServerErrorResponse(w, r, err)
		return
	}
}

// @Summary		Delete User Account
// @Description	Delete User Account
// @Tags			Users
// @Accept			json
// @Produce		json
// @Security		JWT
// @Success		204
// @Failure		400	{object}	main.Envelope{data=nil,error=string}
// @Failure		401	{object}	main.Envelope{data=nil,error=string}
// @Failure		500	{object}	main.Envelope{data=nil,error=string}
// @Router			/users/delete [delete]
func (app *Application) DeleteUserAccountHandler(w http.ResponseWriter, r *http.Request) {

	user, err := utils.GetContentFromContext[repository.User](r, UsrCtx)
	if err != nil {
		app.InternalServerErrorResponse(w, r, err)
		return
	}

	err = app.Services.Users.DeleteAccount(r.Context(), user.Id)
	if err != nil {
		app.InternalServerErrorResponse(w, r, err)
		return
	}

	if err := app.JsonReponseNoContent(w, r); err != nil {
		app.InternalServerErrorResponse(w, r, err)
		return
	}
}
