package services

import (
	"context"

	"faissal.com/blogSpace/internal/repository"
)

type BlogsServices struct {
	Repo repository.Repository
}

type BlogRequest struct {
	Title         string `json:"title" validate:"required"`
	Content       string `json:"content" validate:"required"`
	Description   string `json:"description" validate:"required"`
	FeaturedImage string `json:"-"`
	CategoryId    int    `json:"category_id"`
	UserId        int    `json:"-"`
}

func (bs *BlogsServices) CreateNewBlog(ctx context.Context, req BlogRequest, errCreate chan<- error) {

	defer close(errCreate)

	select {
	case <-ctx.Done():
		errCreate <- ctx.Err()
		return
	default:
	}

	newBlog := repository.Blog{}
	newBlog.Title = req.Title
	newBlog.Content = req.Content
	newBlog.Description = req.Description
	newBlog.FeaturedImage = req.FeaturedImage
	newBlog.CategoryId = req.CategoryId
	newBlog.UserId = req.UserId

	err := bs.Repo.Blogs.Create(ctx, newBlog)
	if err != nil {
		errCreate <- err
		return
	}

	errCreate <- nil

}

func (bs *BlogsServices) GetBlogById() {

}

func (bs *BlogsServices) DeleteBlog() {

}
