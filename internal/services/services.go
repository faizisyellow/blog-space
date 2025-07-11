package services

import (
	"context"
	"database/sql"

	"faissal.com/blogSpace/internal/db"
	"faissal.com/blogSpace/internal/repository"
)

type Services struct {
	Users interface {
		RegisterAccount(ctx context.Context, req RegisterRequest) (*RegisterResponse, error)

		ActivateAccount(ctx context.Context, token string) error

		Login(ctx context.Context, req LoginRequest) (*repository.User, error)

		DeleteAccount(ctx context.Context, usrid int) error

		GetUseById(ctx context.Context, usrid int) (repository.User, error)
	}

	Categories interface {
		CreateNewCategory(ctx context.Context, req CategoryRequest) error

		GetCategory(ctx context.Context, id int) (repository.Category, error)

		GetCategories(ctx context.Context) ([]*repository.Category, error)

		UpdateCategory(ctx context.Context, extCat repository.Category, req UpdateCategoryRequest) error

		DeleteCategory(ctx context.Context, catId int) error
	}
}

func NewServices(store repository.Repository, txfnc db.TransFnc, Db *sql.DB) *Services {
	return &Services{
		Users: &UsersServices{Repo: store, TransFnc: txfnc, Db: Db},

		Categories: &CategorisServices{Repo: store},
	}
}
