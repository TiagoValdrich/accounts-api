package repository

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/tiagovaldrich/accounts-api/internal/models"
	"github.com/uptrace/bun"
)

type CustomerAccountRepository interface {
	Base
	CreateCustomerAccount(
		ctx context.Context,
		customerAccount models.CustomerAccount,
	) (*models.CustomerAccount, error)
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

	if err != nil {
		log.Err(err).
			Str("customer_id", customerAccount.ID.String()).
			Msg("failed to create customer account")

		return nil, err
	}

	return &customerAccount, err
}
