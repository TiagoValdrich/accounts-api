package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gofrs/uuid/v5"
	"github.com/tiagovaldrich/accounts-api/internal/models"
	"github.com/uptrace/bun"
)

type BalanceRepository interface {
	Base
	CreateCustomerBalance(ctx context.Context, customerAccountBalance models.Balance) (*models.Balance, error)
	GetCustomerAccountBalance(ctx context.Context, customerAccountID *uuid.UUID) (*models.Balance, error)
	UpdateCustomerAccountBalance(
		ctx context.Context, newBalance models.Balance,
	) (models.Balance, error)
}

type balanceRepository struct {
	BaseRepo
}

func NewBalanceRepository(db bun.IDB) BalanceRepository {
	repo := &balanceRepository{}
	repo.SetDB(db)

	return repo
}

func (br *balanceRepository) CreateCustomerBalance(ctx context.Context, customerAccountBalance models.Balance) (*models.Balance, error) {
	_, err := br.GetDB(ctx).
		NewInsert().
		Model(&customerAccountBalance).
		Exec(ctx)

	return &customerAccountBalance, err
}

func (br *balanceRepository) GetCustomerAccountBalance(ctx context.Context, customerAccountID *uuid.UUID) (*models.Balance, error) {
	var result models.Balance

	err := br.GetDB(ctx).
		NewSelect().
		Model(&result).
		Where("customer_account_id = ?", customerAccountID).
		For("UPDATE").
		Scan(ctx, &result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &result, nil
}

func (br *balanceRepository) UpdateCustomerAccountBalance(
	ctx context.Context, newBalance models.Balance,
) (models.Balance, error) {
	_, err := br.GetDB(ctx).
		NewUpdate().
		Model(&newBalance).
		Where("customer_account_id = ?", newBalance.CustomerAccountID).
		Exec(ctx)

	return newBalance, err
}
