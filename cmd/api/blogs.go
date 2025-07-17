package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"faissal.com/blogSpace/internal/repository"
	"faissal.com/blogSpace/internal/services"
	"faissal.com/blogSpace/internal/utils"
)

// @Summary		Create New Blog
// @Description	Create New Blog
// @Tags			Blogs
// @Accept			mpfd
// @Produce		json
// @Security		JWT
// @Param			properties	formData	string	true	"Payload To Create New Blog"
// @Param			file		formData	file	true	"Image file"
// @Success		201			{object}	main.Envelope{data=string,error=nil}
// @Failure		400			{object}	main.Envelope{data=nil,error=string}
// @Failure		401			{object}	main.Envelope{data=nil,error=string}
// @Failure		500			{object}	main.Envelope{data=nil,error=string}
// @Router			/blogs [post]
func (app *Application) CreateBlogHandler(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	errUpload := make(chan error, 1)
	errCreate := make(chan error, 1)

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		app.InternalServerErrorResponse(w, r, err)
		return
	}

	reqProps := services.BlogRequest{}

	userAuth, err := utils.GetContentFromContext[repository.User](r, UsrCtx)
	if err != nil {
		app.InternalServerErrorResponse(w, r, err)
		return
	}

	reqProps.UserId = userAuth.Id
	reqProps.FeaturedImage = fileHeader.Filename

	if err := ReadJsonMultiPartForm(r, "properties", &reqProps); err != nil {
		app.BadRequestErrorResponse(w, r, err)
		return
	}

	if err := Validate.Struct(reqProps); err != nil {
		app.BadRequestErrorResponse(w, r, err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		app.Uploading.UploadFile(ctx, fileHeader.Filename, file, fileHeader.Header.Get("Content-Type"), errUpload)
	}()

	go func() {
		defer wg.Done()
		app.Services.Blogs.CreateNewBlog(ctx, reqProps, errCreate)
	}()

	err = <-errUpload
	if err != nil {
		app.InternalServerErrorResponse(w, r, err)
		return
	}

	// Wait for one of the tasks to fail or both to succeed
	var uploadErr, createErr error
	done := make(chan struct{})

	go func() {
		uploadErr = <-errUpload
		if uploadErr != nil {
			cancel()
		}
		createErr = <-errCreate
		if createErr != nil {
			cancel()
		}
		close(done)
	}()

	<-done
	wg.Wait()

	if uploadErr != nil {
		app.InternalServerErrorResponse(w, r, fmt.Errorf("upload failed: %w", uploadErr))
		return
	}
	if createErr != nil {
		app.InternalServerErrorResponse(w, r, fmt.Errorf("create blog failed: %w", createErr))
		return
	}

	if err := app.JsonSuccessReponse(w, r, "create blog successfully", http.StatusCreated); err != nil {
		app.InternalServerErrorResponse(w, r, err)
		return
	}
}
