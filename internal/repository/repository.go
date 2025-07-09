package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Users interface {
	Create(ctx context.Context, tx *sql.Tx, usr User) (usrId int, err error)

	GetById(ctx context.Context, usrId int) (*User, error)

	GetByEmail(ctx context.Context, usrEmail string) (*User, error)

	Update(ctx context.Context, tx *sql.Tx, usrId int, usr User) error

	Delete(ctx context.Context, tx *sql.Tx, usrId int) error
}

type Invitations interface {
	Create(ctx context.Context, tx *sql.Tx, ivt Invitation) error

	GetByUserId(ctx context.Context, tx *sql.Tx, usrId int) (id int, err error)

	DeleteByUserId(ctx context.Context, tx *sql.Tx, usrId int) error
}

type Repository struct {
	Users Users

	Invitations Invitations
}

var (
	QueryTimeout   = time.Second * 5
	ErrNotAffected = errors.New("errors rows not affected")
)

func NewRepostory(db *sql.DB) *Repository {

	return &Repository{
		Users:       &UsersRepository{Db: db},
		Invitations: &InvitationsRepository{Db: db},
	}
}
