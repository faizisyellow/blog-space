package main

import (
	"net/http"
	"strconv"

	"faissal.com/blogSpace/internal/repository"
	"faissal.com/blogSpace/internal/services"
	"faissal.com/blogSpace/internal/utils"
	"github.com/go-chi/chi/v5"
)

//	@Summary		Create New Comment
//	@Description	Create New Comment By User authentication
//	@Tags			Comments
//	@Accept			json
//	@Produce		json
//	@Security		JWT
//	@Param			payload	body		services.CommentRequest	true	"Payload To Create Comments"
//	@Success		201		{object}	main.Envelope{data=string,error=nil}
//	@Failure		400		{object}	main.Envelope{data=nil,error=string}
//	@Failure		401		{object}	main.Envelope{data=nil,error=string}
//	@Failure		500		{object}	main.Envelope{data=nil,error=string}
//	@Router			/comments [post]
func (app *Application) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {

	req := services.CommentRequest{}

	if err := ReadHttpJson(w, r, &req); err != nil {
		app.BadRequestErrorResponse(w, r, err)
		return
	}

	author, err := utils.GetContentFromContext[repository.User](r, UsrCtx)
	if err != nil {
		app.InternalServerErrorResponse(w, r, err)
		return
	}

	err = app.Services.Comments.CreateNewComment(r.Context(), req, author.Id)
	if err != nil {
		app.InternalServerErrorResponse(w, r, err)
		return
	}

	if err := app.JsonSuccessReponse(w, r, "comment created successfully", http.StatusCreated); err != nil {

		app.InternalServerErrorResponse(w, r, err)
		return
	}

}

//	@Summary		Delete Comment
//	@Description	Delete Comment Either By comment's author or Blog's author
//	@Tags			Comments
//	@Accept			json
//	@Produce		json
//	@Security		JWT
//	@Param			ID	path	int	true	"Id Comment"
//	@Success		204
//	@Failure		400	{object}	main.Envelope{data=nil,error=string}
//	@Failure		401	{object}	main.Envelope{data=nil,error=string}
//	@Failure		403	{object}	main.Envelope{data=nil,error=string}
//	@Failure		500	{object}	main.Envelope{data=nil,error=string}
//	@Router			/comments/{ID} [delete]
func (app *Application) DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "ID"))
	if err != nil {
		app.BadRequestErrorResponse(w, r, err)
		return
	}

	err = app.Services.Comments.DeleteComment(r.Context(), id)
	if err != nil {
		app.InternalServerErrorResponse(w, r, err)
		return
	}

	if err := app.JsonReponseNoContent(w, r); err != nil {
		app.InternalServerErrorResponse(w, r, err)
		return
	}
}
