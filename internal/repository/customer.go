package repository

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/tiagovaldrich/accounts-api/internal/models"
	"github.com/uptrace/bun"
)

type CustomerRepository interface {
	Base
	CreateCustomer(context.Context, models.Customer) (*models.Customer, error)
}

type customerRepository struct {
	BaseRepo
}

func NewCustomerRepository(db bun.IDB) CustomerRepository {
	repo := &customerRepository{}
	repo.SetDB(db)

	return repo
}

func (r *customerRepository) CreateCustomer(ctx context.Context, customer models.Customer) (*models.Customer, error) {
	_, err := r.GetDB(ctx).
		NewInsert().
		Model(&customer).
		Exec(ctx)

	if err != nil {
		log.Err(err).Msg("failed to create customer")

		return nil, err
	}

	return &customer, err
}
