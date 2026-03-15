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
	CustomerAccountID *uuid.UUID
	Document          string
	CreatedAt         time.Time
}

func DatabaseToSearchCustomerAccountResult(dbResult repository.CustomerAccountByIDResult) SearchCustomerAccountResult {
	return SearchCustomerAccountResult{
		CustomerAccountID: dbResult.ID,
		Document:          dbResult.Document,
		CreatedAt:         dbResult.CreatedAt,
	}
}
