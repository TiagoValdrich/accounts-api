package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gofrs/uuid/v5"
	"github.com/tiagovaldrich/accounts-api/internal/models"
	"github.com/uptrace/bun"
)

type CustomerAccountRepository interface {
	Base
	CreateCustomerAccount(
		ctx context.Context,
		customerAccount models.CustomerAccount,
	) (*models.CustomerAccount, error)
	SearchCustomerAccountByID(
		ctx context.Context, customerAccountID *uuid.UUID,
	) (*CustomerAccountByIDResult, error)
}

type customerAccountRepository struct {
	BaseRepo
}

func NewCustomerAccountRepository(db bun.IDB) CustomerAccountRepository {
	repo := &customerAccountRepository{}
	repo.SetDB(db)

	return repo
}

func (r *customerAccountRepository) CreateCustomerAccount(
	ctx context.Context,
	customerAccount models.CustomerAccount,
) (*models.CustomerAccount, error) {
	_, err := r.GetDB(ctx).
		NewInsert().
		Model(&customerAccount).
		Exec(ctx)

	return &customerAccount, err
}

func (r *customerAccountRepository) SearchCustomerAccountByID(
	ctx context.Context, customerAccountID *uuid.UUID,
) (*CustomerAccountByIDResult, error) {
	var result CustomerAccountByIDResult

	err := r.GetDB(ctx).
		NewSelect().
		Model(&result).
		ColumnExpr("c.document AS document").
		ColumnExpr("ca.id AS id").
		ColumnExpr("ca.created_at AS created_at").
		TableExpr("customer_account ca").
		Join("JOIN customer c ON c.id = ca.customer_id").
		Where("ca.id = ?", customerAccountID).
		Scan(ctx, &result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &result, nil
}
