package accounts

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type AccountCreatedResponse struct {
	ID        *uuid.UUID `json:"account_id"`
	Document  string     `json:"document_number"`
	CreatedAt time.Time  `json:"created_at"`
}

type SearchCustomerAccountByIDResponse struct {
	ID       *uuid.UUID `json:"account_id"`
	Document string     `json:"document_number"`
}

func DomainToAccountCreatedResponse(customerAccountResult CustomerAccountResult) AccountCreatedResponse {
	return AccountCreatedResponse{
		ID:        customerAccountResult.CustomerAccount.ID,
		Document:  customerAccountResult.Customer.Document,
		CreatedAt: customerAccountResult.Customer.CreatedAt,
	}
}

func DomainToSearchAccountByIDResponse(searchCustomerAccountResult SearchCustomerAccountResult) SearchCustomerAccountByIDResponse {
	return SearchCustomerAccountByIDResponse{
		ID:       searchCustomerAccountResult.CustomerID,
		Document: searchCustomerAccountResult.Document,
	}
}
