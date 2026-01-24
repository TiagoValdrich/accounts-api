package accounts

import "github.com/tiagovaldrich/accounts-api/internal/models"

type CustomerAccountResult struct {
	Customer        *models.Customer
	CustomerAccount *models.CustomerAccount
}
