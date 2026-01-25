package accounts

import (
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/tiagovaldrich/accounts-api/internal/models"
	"github.com/tiagovaldrich/accounts-api/internal/repository"
)

type CustomerAccountResult struct {
	Customer        *models.Customer
	CustomerAccount *models.CustomerAccount
}

type SearchCustomerAccountResult struct {
	CustomerID *uuid.UUID
	Document   string
	CreatedAt  time.Time
}

func DatabaseToSearchCustomerAccountResult(dbResult repository.CustomerAccountByIDResult) SearchCustomerAccountResult {
	return SearchCustomerAccountResult{
		CustomerID: dbResult.ID,
		Document:   dbResult.Document,
		CreatedAt:  dbResult.CreatedAt,
	}
}
