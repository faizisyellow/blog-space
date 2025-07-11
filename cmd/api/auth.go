package main

import (
	"fmt"
	"net/http"

	"faissal.com/blogSpace/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// @Summary		Sign Up Account
// @Description	Sign Up New Account
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			payload	body		services.RegisterRequest	true	"Payload to Sign Up"
// @Success		201		{object}	main.Envelope{data=services.RegisterResponse,error=nil}
// @Failure		400		{object}	main.Envelope{data=nil,error=string}
// @Failure		409		{object}	main.Envelope{data=nil,error=string}
// @Failure		500		{object}	main.Envelope{data=nil,error=string}
// @Router			/authentication/sign-up [post]
func (app *Application) SignUpHandler(w http.ResponseWriter, r *http.Request) {

	req := services.RegisterRequest{}

	if err := ReadHttpJson(w, r, &req); err != nil {
		app.BadRequestErrorResponse(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.BadRequestErrorResponse(w, r, err)
		return
	}

	res, err := app.Services.Users.RegisterAccount(r.Context(), req)
	if err != nil {
		switch err {
		case services.ErrUserAlreadyExist:

			app.ConflictErrorResponse(w, r, fmt.Errorf("email not available"))
		default:
			app.InternalServerErrorResponse(w, r, err)
		}
		return
	}

	if err := app.JsonSuccessReponse(w, r, res, http.StatusCreated); err != nil {
		app.InternalServerErrorResponse(w, r, err)
		return
	}

}

// @Summary		Activate Account
// @Description	Activate  New Account
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			token	path		string	true	"Token Inivitation to Activate Account"
// @Success		200		{object}	main.Envelope{data=string,error=nil}
// @Failure		400		{object}	main.Envelope{data=nil,error=string}
// @Failure		404		{object}	main.Envelope{data=nil,error=string}
// @Failure		500		{object}	main.Envelope{data=nil,error=string}
// @Router			/authentication/activation/{token} [post]
func (app *Application) ActivateAccountHandler(w http.ResponseWriter, r *http.Request) {

	req := chi.URLParam(r, "token")

	if req == "" {
		app.BadRequestErrorResponse(w, r, fmt.Errorf("token invition empty"))
	}

	err := app.Services.Users.ActivateAccount(r.Context(), req)
	if err != nil {
		switch err {
		case services.ErrTokenInvitationNotFound:
			app.NotFoundErrorResponse(w, r, err)
		case services.ErrUserRegisteredNotFound:
			app.NotFoundErrorResponse(w, r, err)
		default:
			app.InternalServerErrorResponse(w, r, err)
		}

		return
	}

	if err := app.JsonSuccessReponse(w, r, "user activated successfully", http.StatusOK); err != nil {
		app.InternalServerErrorResponse(w, r, err)
		return
	}

}

// @Summary		Sign in Account
// @Description	Sign in  New Account
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			payload	body		services.LoginRequest	true	"Token to Sign in 	Account"
// @Success		200		{object}	main.Envelope{data=services.LoginResponse,error=nil}
// @Failure		400		{object}	main.Envelope{data=nil,error=string}
// @Failure		404		{object}	main.Envelope{data=nil,error=string}
// @Failure		500		{object}	main.Envelope{data=nil,error=string}
// @Router			/authentication/sign-in [post]
func (app *Application) SignInHandler(w http.ResponseWriter, r *http.Request) {

	req := services.LoginRequest{}

	if err := ReadHttpJson(w, r, &req); err != nil {
		app.BadRequestErrorResponse(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.BadRequestErrorResponse(w, r, err)
		return
	}

	user, err := app.Services.Users.Login(r.Context(), req)
	if err != nil {
		switch err {
		case services.ErrUserNotFound:

			app.NotFoundErrorResponse(w, r, err)
		case services.ErrUserNotActivated:

			app.BadRequestErrorResponse(w, r, err)
		case bcrypt.ErrMismatchedHashAndPassword:

			app.BadRequestErrorResponse(w, r, fmt.Errorf("email or password incorrect"))
		default:

			app.InternalServerErrorResponse(w, r, err)
		}

		return
	}

	claims := jwt.MapClaims{
		"iss": app.JwtAuth.Iss,
		"sub": app.JwtAuth.Sub,
		"exp": app.JwtAuth.Exp,
		"id":  user.Id,
	}

	token, err := app.Authentication.GenerateToken(claims)
	if err != nil {
		app.InternalServerErrorResponse(w, r, err)
		return
	}

	if err := app.JsonSuccessReponse(w, r, services.LoginResponse{Token: token}, http.StatusOK); err != nil {
		app.InternalServerErrorResponse(w, r, err)
		return
	}
}
