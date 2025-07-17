package services

import (
	"context"
	"database/sql"
	"errors"

	"faissal.com/blogSpace/internal/repository"
)

type CommentsServices struct {
	Repo repository.Repository
}

type CommentRequest struct {
	Content string `json:"content"`
	BlogId  int    `json:"blog_id"`
}

var (
	ErrCommentNotFound = errors.New("comment with this id not found")
)

func (cs *CommentsServices) CreateNewComment(ctx context.Context, req CommentRequest, authorId int) error {

	newCmt := repository.Comment{}
	newCmt.Content = req.Content
	newCmt.BlogId = req.BlogId
	newCmt.UserId = authorId

	err := cs.Repo.Comments.Create(ctx, newCmt)
	if err != nil {
		return err
	}

	return nil

}

func (cs *CommentsServices) GetCommentById(ctx context.Context, id int) (*repository.Comment, error) {

	comment, err := cs.Repo.Comments.GetById(ctx, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrCommentNotFound
		default:
			return nil, err
		}
	}

	return comment, nil
}

func (cs *CommentsServices) DeleteComment(ctx context.Context, id int) error {

	err := cs.Repo.Comments.DeleteByUserId(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
