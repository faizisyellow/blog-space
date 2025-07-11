package repository

import (
	"context"
	"database/sql"
)

type CategoriesRepository struct {
	Db *sql.DB
}

type Category struct {
	Id      int
	Content string
}

func (c *CategoriesRepository) Create(ctx context.Context, nwCat Category) error {

	query := `INSERT INTO categories(content) VALUES(?)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	_, err := c.Db.ExecContext(ctx, query, nwCat.Content)
	if err != nil {
		return err
	}

	return nil
}

func (c *CategoriesRepository) GetById(ctx context.Context, catId int) (Category, error) {

	query := `SELECT id,content FROM categories WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	cat := Category{}

	err := c.Db.QueryRowContext(ctx, query, catId).Scan(
		&cat.Id,
		&cat.Content,
	)

	if err != nil {
		return Category{}, err
	}

	return cat, nil
}

func (c *CategoriesRepository) GetAll(ctx context.Context) ([]*Category, error) {

	query := `SELECT id,content FROM categories ORDER BY id`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	row, err := c.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer row.Close()

	cats := []*Category{}

	for row.Next() {

		cat := Category{}

		err := row.Scan(&cat.Id, (&cat.Content))
		if err != nil {
			return nil, err
		}

		cats = append(cats, &cat)
	}

	err = row.Err()
	if err != nil {
		return nil, err
	}

	return cats, nil

}

func (c *CategoriesRepository) Update(ctx context.Context, cat Category) error {

	query := `UPDATE categories SET content = ? WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	_, err := c.Db.ExecContext(ctx, query, cat.Content, cat.Id)
	if err != nil {
		return err
	}

	return nil
}

func (c *CategoriesRepository) Delete(ctx context.Context, catId int) error {

	query := `DELETE FROM categories WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	_, err := c.Db.ExecContext(ctx, query, catId)
	if err != nil {
		return err
	}

	return nil
}
