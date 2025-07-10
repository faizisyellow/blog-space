package repository

import (
	"context"
	"database/sql"
	"time"
)

type Invitation struct {
	UserId   int           `json:"user_id"`
	Token    string        `json:"token"`
	ExpireAt time.Duration `json:"expire_at"`
}

type InvitationsRepository struct {
	Db *sql.DB
}

func (ir *InvitationsRepository) Create(ctx context.Context, tx *sql.Tx, ivt Invitation) error {

	query := `INSERT INTO invitations(user_id,token,expire_at)
	VALUES(?,?,?)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, ivt.UserId, ivt.Token, time.Now().Add(ivt.ExpireAt))
	if err != nil {
		return err
	}

	return nil
}

func (ir *InvitationsRepository) GetByUserId(ctx context.Context, tx *sql.Tx, token string) (id int, err error) {

	// query := `SELECT user_id FROM invitations WHERE token = "0fc3ef081ea829c30401c17df1f204e5c780ffc3fe25fef0ecb8553f0d3eb6bb";`
	query := `SELECT user_id FROM invitations WHERE token = ? AND expire_at > ?;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	userId := 0

	row := tx.QueryRowContext(ctx, query, token, time.Now())
	err = row.Scan(&userId)

	if err != nil {
		return 0, err
	}

	return userId, nil
}

func (ir *InvitationsRepository) DeleteByUserId(ctx context.Context, tx *sql.Tx, usrId int) error {

	query := `DELETE FROM invitations WHERE user_id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	res, err := tx.ExecContext(ctx, query, usrId)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrNotAffected
	}

	return nil
}
