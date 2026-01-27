package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tiagovaldrich/accounts-api/internal/models"
	"github.com/uptrace/bun"
)

type TransactionRepository interface {
	Base
	GetTransactionByIdempotencyKey(ctx context.Context, idempotencyKey string) (*models.Transaction, error)
	CreateTransaction(context.Context, models.Transaction) (*models.Transaction, error)
}

type transactionRepository struct {
	BaseRepo
}

func NewTransactionRepository(db bun.IDB) TransactionRepository {
	repo := &transactionRepository{}
	repo.SetDB(db)

	return repo
}

func (tr *transactionRepository) GetTransactionByIdempotencyKey(ctx context.Context, idempotencyKey string) (*models.Transaction, error) {
	var result models.Transaction

	err := tr.GetDB(ctx).
		NewSelect().
		Model(&result).
		Where("idempotency_key = ?", idempotencyKey).
		Scan(ctx, &result)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &result, nil
}

func (tr *transactionRepository) CreateTransaction(ctx context.Context, transaction models.Transaction) (*models.Transaction, error) {
	_, err := tr.GetDB(ctx).
		NewInsert().
		Model(&transaction).
		Exec(ctx)

	return &transaction, err
}
