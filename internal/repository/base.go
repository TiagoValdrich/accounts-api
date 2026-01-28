package repository

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
)

type ContextKey string

const (
	TxKey ContextKey = "database_tx"
)

type Base interface {
	WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
	GetDB(ctx context.Context) bun.IDB
	SetDB(db bun.IDB)
}

type BaseRepo struct {
	db bun.IDB
}

func (br *BaseRepo) WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	return br.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		txCtx := context.WithValue(ctx, TxKey, tx)

		return fn(txCtx)
	})
}

func (br *BaseRepo) GetDB(ctx context.Context) bun.IDB {
	if txCtx, ok := ctx.Value(TxKey).(bun.Tx); ok {
		return txCtx
	}

	return br.db
}

func (br *BaseRepo) SetDB(db bun.IDB) {
	br.db = db
}
