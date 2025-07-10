package repository

import (
	"context"
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type UsersRepository struct {
	Db *sql.DB
}

type User struct {
	Id        int            `json:"id"`
	Username  string         `json:"username"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Email     string         `json:"email"`
	Password  HashedPassword `json:"-"`
	IsActive  bool           `json:"is_active"`
}

type HashedPassword struct {
	Text   *string
	Hashed []byte
}

func (hp *HashedPassword) Parse(text string) error {

	hashed, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	hp.Text = &text
	hp.Hashed = hashed

	return nil

}

func (hp *HashedPassword) Compare(pass string) error {

	err := bcrypt.CompareHashAndPassword(hp.Hashed, []byte(pass))
	if err != nil {
		return err
	}

	return nil

}

func (ur *UsersRepository) Create(ctx context.Context, tx *sql.Tx, usr User) (usrId int, err error) {

	query := `INSERT INTO users(username,first_name,last_name,email,password)
	VALUES(?,?,?,?,?)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	res, err := tx.ExecContext(ctx, query, usr.Username, usr.FirstName, usr.LastName, usr.Email, usr.Password.Hashed)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (ur *UsersRepository) GetById(ctx context.Context, usrId int) (*User, error) {

	query := `SELECT id,username,first_name,last_name,email,is_active,password FROM users WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	user := &User{}

	err := ur.Db.QueryRowContext(ctx, query, usrId).Scan(
		&user.Id,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.IsActive,
		&user.Password.Hashed,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UsersRepository) GetByEmail(ctx context.Context, usrEmail string) (*User, error) {

	query := `SELECT id,username,first_name,last_name,email,password,is_active FROM users WHERE email = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	user := &User{}

	err := ur.Db.QueryRowContext(ctx, query, usrEmail).Scan(
		&user.Id,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password.Hashed,
		&user.IsActive,
	)

	if err != nil {
		return nil, err
	}

	return user, nil

}

func (ur *UsersRepository) Update(ctx context.Context, tx *sql.Tx, usrId int, usr User) error {

	query := `UPDATE users SET username = ?, first_name = ?, last_name = ?, email = ?, password = ?, is_active = ? WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	res, err := tx.ExecContext(ctx, query, &usr.Username, &usr.FirstName, &usr.LastName, &usr.Email, &usr.Password.Hashed, &usr.IsActive, usrId)
	if err != nil {
		return err
	}

	afct, err := res.RowsAffected()
	if err != nil {
		return err
	}

	// TODO: refactor https://pkg.go.dev/database/sql#Conn.ExecContext
	if afct == 0 {
		return ErrNotAffected
	}

	return nil
}

func (ur *UsersRepository) Delete(ctx context.Context, tx *sql.Tx, usrId int) error {

	query := `DELETE FROM users WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	res, err := tx.ExecContext(ctx, query, usrId)
	if err != nil {
		return err
	}

	afct, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if afct == 0 {
		return ErrNotAffected
	}

	return nil
}
