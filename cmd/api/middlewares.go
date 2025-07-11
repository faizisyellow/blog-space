package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"faissal.com/blogSpace/internal/keys"
	"faissal.com/blogSpace/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

var UsrCtx keys.User = "user"
var CatCtx keys.Category = "category"

func (app *Application) AuthMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			app.BadRequestErrorResponse(w, r, fmt.Errorf("authorization is missing"))
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 {
			app.BadRequestErrorResponse(w, r, fmt.Errorf("authorization is malformed"))
			return
		}

		if parts[0] != "Bearer" {
			app.BadRequestErrorResponse(w, r, fmt.Errorf("authorization is malformed: authentication use Bearer"))
			return
		}

		token := parts[1]

		jwtToken, err := app.Authentication.VerifyToken(token)
		if err != nil {
			app.BadRequestErrorResponse(w, r, err)
			return
		}

		claim, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok {
			app.InternalServerErrorResponse(w, r, fmt.Errorf("error while parsing token claim type"))
			return
		}

		usrId, ok := claim["id"].(float64)
		if !ok {
			app.InternalServerErrorResponse(w, r, fmt.Errorf("error while parsing field claim type"))
			return
		}

		ctx := r.Context()

		user, err := app.Services.Users.GetUseById(ctx, int(usrId))
		if err != nil {
			switch err {
			case services.ErrUserNotFound:
				app.UnAuthorizeErrorResponse(w, r, err)
			default:
				app.InternalServerErrorResponse(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, UsrCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func (app *Application) GetCategoryMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		id := chi.URLParam(r, "ID")

		catId, err := strconv.Atoi(id)
		if err != nil {
			app.InternalServerErrorResponse(w, r, err)
			return
		}

		ctx := r.Context()

		cat, err := app.Services.Categories.GetCategory(ctx, catId)
		if err != nil {
			switch err {
			case services.ErrNotFoundCategory:
				app.NotFoundErrorResponse(w, r, err)
			default:
				app.InternalServerErrorResponse(w, r, err)
			}

			return
		}

		ctx = context.WithValue(ctx, CatCtx, cat)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
