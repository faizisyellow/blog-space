package services

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"faissal.com/blogSpace/internal/repository"
)

type CategorisServices struct {
	Repo repository.Repository
}

type CategoryRequest struct {
	Content string `json:"content" validate:"required,min=3,max=255"`
}

type UpdateCategoryRequest struct {
	Content string `json:"content" validate:"max=255"`
}

var (
	ErrDuplicateCategory = errors.New("category already exist")
	ErrNotFoundCategory  = errors.New("category not found")
)

func (cs *CategorisServices) CreateNewCategory(ctx context.Context, req CategoryRequest) error {

	newCat := repository.Category{}
	newCat.Content = req.Content

	err := cs.Repo.Categories.Create(ctx, newCat)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), repository.DUPLICATE_CODE):
			return ErrDuplicateCategory
		default:
			return err
		}
	}

	return nil
}

func (cs *CategorisServices) GetCategory(ctx context.Context, id int) (repository.Category, error) {

	cat, err := cs.Repo.Categories.GetById(ctx, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return repository.Category{}, ErrNotFoundCategory
		default:
			return repository.Category{}, err
		}
	}

	return cat, nil
}

func (cs *CategorisServices) GetCategories(ctx context.Context) ([]*repository.Category, error) {

	cats, err := cs.Repo.Categories.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return cats, nil
}

func (cs *CategorisServices) UpdateCategory(ctx context.Context, extCat repository.Category, req UpdateCategoryRequest) error {

	// if the request still the same, dont update it.
	if extCat.Content == req.Content {
		return nil
	}

	extCat.Content = req.Content

	err := cs.Repo.Categories.Update(ctx, extCat)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), repository.DUPLICATE_CODE):
			return ErrDuplicateCategory
		default:
			return err
		}
	}

	return nil

}

func (cs *CategorisServices) DeleteCategory(ctx context.Context, catId int) error {

	err := cs.Repo.Categories.Delete(ctx, catId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrNotFoundCategory
		default:
			return err

		}
	}

	return nil
}
