package accounts

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type AccountCreatedResponse struct {
	ID        *uuid.UUID `json:"id"`
	Document  string     `json:"document"`
	CreatedAt time.Time  `json:"created_at"`
}

func toAccountCreatedResponse(customerAccountResult CustomerAccountResult) AccountCreatedResponse {
	return AccountCreatedResponse{
		ID:        customerAccountResult.CustomerAccount.ID,
		Document:  customerAccountResult.Customer.Document,
		CreatedAt: customerAccountResult.Customer.CreatedAt,
	}
}
