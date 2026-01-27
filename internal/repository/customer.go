package repository

import (
	"context"

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

	return &customer, err
}
