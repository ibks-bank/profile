package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

const errViolatesUnique = "violates unique"

type store struct {
	db *sql.DB
}

func New(db *sql.DB) *store {
	return &store{
		db: db,
	}
}

type TxFunc func(ctx context.Context, tx *sql.Tx) error

func (st *store) WithTransaction(ctx context.Context, fn TxFunc) (err error) {
	tx, err := st.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(ctx, tx)
	return err
}
