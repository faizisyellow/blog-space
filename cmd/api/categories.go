package main

import (
	"net/http"

	"faissal.com/blogSpace/internal/repository"
	"faissal.com/blogSpace/internal/services"
	"faissal.com/blogSpace/internal/utils"
)

// @Summary		Create New Category
// @Description	Create New Category
// @Tags			Categories
// @Accept			json
// @Produce		json
// @Security		JWT
// @Param			payload	body		services.CategoryRequest	true	"Payload New Category"
// @Success		201		{object}	main.Envelope{data=string,error=nil}
// @Failure		400		{object}	main.Envelope{data=nil,error=string}
// @Failure		401		{object}	main.Envelope{data=nil,error=string}
// @Failure		409		{object}	main.Envelope{data=nil,error=string}
// @Failure		500		{object}	main.Envelope{data=nil,error=string}
// @Router			/categories [post]
func (app *Application) CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {

	newCat := services.CategoryRequest{}

	if err := ReadHttpJson(w, r, &newCat); err != nil {
		app.BadRequestErrorResponse(w, r, err)
		return
	}

	if err := Validate.Struct(newCat); err != nil {
		app.BadRequestErrorResponse(w, r, err)
		return
	}

	err := app.Services.Categories.CreateNewCategory(r.Context(), newCat)
	if err != nil {
		switch err {
		case services.ErrDuplicateCategory:
			app.ConflictErrorResponse(w, r, err)
		default:
			app.InternalServerErrorResponse(w, r, err)
		}
		return
	}

	if err := app.JsonSuccessReponse(w, r, "create new category successfully", http.StatusCreated); err != nil {
		app.InternalServerErrorResponse(w, r, err)
		return
	}

}

// @Summary		Get Category
// @Description	Get Category By Id
// @Tags			Categories
// @Accept			json
// @Produce		json
// @Security		JWT
// @Param			ID	path		int	true	"Category ID"
// @Success		200	{object}	main.Envelope{data=repository.Category,error=nil}
// @Failure		400	{object}	main.Envelope{data=nil,error=string}
// @Failure		404	{object}	main.Envelope{data=nil,error=string}
// @Failure		500	{object}	main.Envelope{data=nil,error=string}
// @Router			/categories/{ID} [get]
func (app *Application) GetCategoryByIdHandler(w http.ResponseWriter, r *http.Request) {

	cat, err := utils.GetContentFromContext[repository.Category](r, CatCtx)
	if err != nil {

		app.InternalServerErrorResponse(w, r, err)
		return
	}

	if err := app.JsonSuccessReponse(w, r, cat, http.StatusOK); err != nil {

		app.InternalServerErrorResponse(w, r, err)
		return
	}
}

// @Summary		Get Categories
// @Description	Get All Categories
// @Tags			Categories
// @Accept			json
// @Produce		json
// @Success		200	{object}	main.Envelope{data=[]repository.Category,error=nil}
// @Failure		400	{object}	main.Envelope{data=nil,error=string}
// @Failure		404	{object}	main.Envelope{data=nil,error=string}
// @Failure		500	{object}	main.Envelope{data=nil,error=string}
// @Router			/categories [get]
func (app *Application) GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {

	cats, err := app.Services.Categories.GetCategories(r.Context())
	if err != nil {
		app.InternalServerErrorResponse(w, r, err)
		return
	}

	if err := app.JsonSuccessReponse(w, r, cats, http.StatusOK); err != nil {
		app.InternalServerErrorResponse(w, r, err)
		return
	}

}

// @Summary		Update Category
// @Description	Update Category
// @Tags			Categories
// @Accept			json
// @Produce		json
// @Security		JWT
// @Param			ID		path		int								true	"Category ID"
// @Param			paylod	body		services.UpdateCategoryRequest	true	"Payload To Update Category"
// @Success		200		{object}	main.Envelope{data=string,error=nil}
// @Failure		400		{object}	main.Envelope{data=nil,error=string}
// @Failure		401		{object}	main.Envelope{data=nil,error=string}
// @Failure		404		{object}	main.Envelope{data=nil,error=string}
// @Failure		409		{object}	main.Envelope{data=nil,error=string}
// @Failure		500		{object}	main.Envelope{data=nil,error=string}
// @Router			/categories/{ID} [patch]
func (app *Application) UpdateCategoryHandler(w http.ResponseWriter, r *http.Request) {

	updateCat := services.UpdateCategoryRequest{}

	if err := ReadHttpJson(w, r, &updateCat); err != nil {
		app.BadRequestErrorResponse(w, r, err)
		return
	}

	if err := Validate.Struct(updateCat); err != nil {
		app.BadRequestErrorResponse(w, r, err)
		return
	}

	cat, err := utils.GetContentFromContext[repository.Category](r, CatCtx)
	if err != nil {

		app.InternalServerErrorResponse(w, r, err)
		return
	}

	err = app.Services.Categories.UpdateCategory(r.Context(), cat, updateCat)
	if err != nil {
		switch err {
		case services.ErrDuplicateCategory:
			app.ConflictErrorResponse(w, r, err)
		default:
			app.InternalServerErrorResponse(w, r, err)
		}
		return
	}

	if err := app.JsonSuccessReponse(w, r, "update category successfully", http.StatusCreated); err != nil {
		app.InternalServerErrorResponse(w, r, err)
		return
	}
}

// @Summary		Delete Category
// @Description	Delete Category
// @Tags			Categories
// @Accept			json
// @Produce		json
// @Security		JWT
// @Param			ID	path	int	true	"Category ID"
// @Success		204
// @Failure		400	{object}	main.Envelope{data=nil,error=string}
// @Failure		401	{object}	main.Envelope{data=nil,error=string}
// @Failure		404	{object}	main.Envelope{data=nil,error=string}
// @Failure		500	{object}	main.Envelope{data=nil,error=string}
// @Router			/categories [delete]
func (app *Application) DeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {

	cat, err := utils.GetContentFromContext[repository.Category](r, CatCtx)
	if err != nil {

		app.InternalServerErrorResponse(w, r, err)
		return
	}

	err = app.Services.Categories.DeleteCategory(r.Context(), cat.Id)
	if err != nil {
		switch err {
		case services.ErrNotFoundCategory:
			app.NotFoundErrorResponse(w, r, err)
		default:
			app.InternalServerErrorResponse(w, r, err)
		}
		return
	}

	if err := app.JsonReponseNoContent(w, r); err != nil {
		app.InternalServerErrorResponse(w, r, err)
		return
	}
}
