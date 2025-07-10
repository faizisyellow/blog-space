package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type TransFnc func(db *sql.DB, ctx context.Context, fnc func(*sql.Tx) error) error

func New(addr string, maxOpenConn, maxIdleConn int, maxIdleTime string) (*sql.DB, error) {

	db, err := sql.Open("mysql", addr)
	if err != nil {
		return nil, fmt.Errorf("error open db=%v", err)
	}

	db.SetMaxIdleConns(maxIdleConn)
	db.SetMaxOpenConns(maxOpenConn)

	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, fmt.Errorf("error parse time=%v", err)
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error connect to db=%v", err)
	}

	return db, nil
}

func WithTx(db *sql.DB, ctx context.Context, fnc func(*sql.Tx) error) error {

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error begin tx=%v", err)
	}

	err = fnc(tx)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
