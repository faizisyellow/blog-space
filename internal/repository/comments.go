package repository

import (
	"context"
	"database/sql"
)

type Comment struct {
	Id        int
	Content   string
	UserId    int
	BlogId    int
	CreatedAt string
}

type CommentsRepository struct {
	Db *sql.DB
}

func (cr *CommentsRepository) Create(ctx context.Context, newCmnt Comment) error {

	query := `INSERT INTO comments(content,user_id,blog_id) (?,?,?)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	_, err := cr.Db.ExecContext(ctx, query, newCmnt.Content, newCmnt.UserId, newCmnt.BlogId)
	if err != nil {
		return err
	}

	return nil
}

func (cr *CommentsRepository) GetById(ctx context.Context, id int) (*Comment, error) {

	query := `SELECT id,content,user_id,blog_id FROM comments WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	comment := Comment{}
	err := cr.Db.QueryRowContext(ctx, query, id).Scan(&comment.Id, &comment.Content, &comment.UserId, &comment.BlogId)

	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (cr *CommentsRepository) DeleteByUserId(ctx context.Context, id int) error {

	query := `DELETE FROM comments WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	_, err := cr.Db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
