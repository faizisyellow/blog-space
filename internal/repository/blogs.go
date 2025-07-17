package repository

import (
	"context"
	"database/sql"
)

type Blog struct {
	Id            int
	Title         string
	Content       string
	Description   string
	FeaturedImage string
	CategoryId    int
	UserId        int
	CreatedAt     string
}

type BlogsRepository struct {
	Db *sql.DB
}

func (bg *BlogsRepository) Create(ctx context.Context, newBg Blog) error {

	query := `INSERT blogs(title,content,description,featured_image,category_id,user_id)
	VALUES(?,?,?,?,?,?)
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	_, err := bg.Db.ExecContext(ctx, query, newBg.Title, newBg.Content, newBg.Description, newBg.FeaturedImage, newBg.CategoryId, newBg.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (bg *BlogsRepository) GetAll(ctx context.Context) ([]*Blog, error) {

	// query := `SELECT b.title,b.content,b.description,b.featured_image,b.created_at,u.username,c.content
	// FROM blogs b JOIN users u ON b.user_id = u.id LEFT JOIN categories c ON c.id = b.category_id
	// `

	return nil, nil
}

func (bg *BlogsRepository) GetById(ctx context.Context, blogId int) (*Blog, error) {

	return nil, nil

}

func (bg *BlogsRepository) Update(ctx context.Context, blgId int) error {

	return nil
}

func (bg *BlogsRepository) Delete(ctx context.Context, blgId int) error {

	query := `DELETE FROM blogs WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	_, err := bg.Db.ExecContext(ctx, query, blgId)
	if err != nil {
		return err
	}

	return nil
}
