package services

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"faissal.com/blogSpace/internal/db"
	"faissal.com/blogSpace/internal/repository"
	"faissal.com/blogSpace/internal/utils"
)

type UsersServices struct {
	Repo repository.Repository

	Db *sql.DB

	TransFnc db.TransFnc
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=16"`

	Email string `json:"email" validate:"required,email,min=6,max=32"`

	FirstName string `json:"first_name" validate:"required,min=3,max=16"`

	LastName string `json:"last_name" validate:"max=32"`

	Password string `json:"password" validate:"required,max=18"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}

type LoginRequest struct {
	Email string `json:"email" validate:"required,email"`

	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type ActivatedRequest struct {
	Token string `json:"token"`
}

var (
	ErrTokenInvitationNotFound = errors.New("invitation not found, please register first")
	ErrUserRegisteredNotFound  = errors.New("user not found, please register first")
	ErrUserNotFound            = errors.New("user not found")
	ErrUserNotActivated        = errors.New("user not activated, please activate first")
	ErrUserAlreadyExist        = errors.New("this user already exists")
)

func (us *UsersServices) RegisterAccount(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {

	res := &RegisterResponse{}

	err := utils.IsPasswordValid(req.Password)
	if err != nil {
		return nil, err
	}

	err = us.TransFnc(us.Db, ctx, func(tx *sql.Tx) error {

		newAccount := repository.User{}
		newAccount.Email = req.Email
		newAccount.Username = req.Username
		newAccount.FirstName = req.FirstName
		newAccount.LastName = req.LastName
		newAccount.Password.Parse(req.Password)

		usrId, err := us.Repo.Users.Create(ctx, tx, newAccount)
		if err != nil {

			if strings.Contains(err.Error(), repository.DUPLICATE_CODE) {
				return ErrUserAlreadyExist
			}

			return err
		}

		tokenIvt := utils.GenerateTokenUuid()

		invt := repository.Invitation{
			UserId:   usrId,
			Token:    tokenIvt,
			ExpireAt: time.Hour * 24,
		}

		err = us.Repo.Invitations.Create(ctx, tx, invt)
		if err != nil {
			return err
		}

		// register and invite success, get the token
		res.Token = tokenIvt

		return nil
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (us *UsersServices) ActivateAccount(ctx context.Context, token string) error {

	err := us.TransFnc(us.Db, ctx, func(tx *sql.Tx) error {

		usrId, err := us.Repo.Invitations.GetByUserId(ctx, tx, token)
		if err != nil {
			return ErrTokenInvitationNotFound
		}

		user, err := us.Repo.Users.GetById(ctx, usrId)
		if err != nil {
			return ErrUserRegisteredNotFound
		}

		user.IsActive = true
		err = us.Repo.Users.Update(ctx, tx, user.Id, *user)
		if err != nil {
			return err
		}

		err = us.Repo.Invitations.DeleteByUserId(ctx, tx, user.Id)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil

}

func (us *UsersServices) Login(ctx context.Context, req LoginRequest) (*repository.User, error) {

	user, err := us.Repo.Users.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if !user.IsActive {
		return nil, ErrUserNotActivated
	}

	err = user.Password.Compare(req.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UsersServices) DeleteAccount(ctx context.Context, usrid int) error {

	return us.TransFnc(us.Db, ctx, func(tx *sql.Tx) error {

		err := us.Repo.Users.Delete(ctx, tx, usrid)
		if err != nil {
			return err
		}

		return nil
	})

}

func (us *UsersServices) GetUseById(ctx context.Context, usrid int) (repository.User, error) {

	user, err := us.Repo.Users.GetById(ctx, usrid)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return repository.User{}, ErrUserNotFound
		default:
			return repository.User{}, err
		}

	}

	return *user, nil
}
