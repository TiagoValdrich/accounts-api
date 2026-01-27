package transactions

import (
	"github.com/gofrs/uuid/v5"
	"github.com/tiagovaldrich/accounts-api/internal/models"
	"github.com/tiagovaldrich/accounts-api/internal/pkg/utils"
)

type CreateTransactionResponse struct {
	ID                *uuid.UUID           `json:"id"`
	CustomerAccountID *uuid.UUID           `json:"customer_account_id"`
	OperationType     models.OperationType `json:"operation_type"`
	Amount            float64              `json:"amount"`
}

func DomainToCreateTransactionResponse(result CreateTransactionResult) CreateTransactionResponse {
	return CreateTransactionResponse{
		ID:                result.ID,
		CustomerAccountID: result.CustomerAccountID,
		OperationType:     result.OperationType,
		Amount:            utils.FromCents(result.Amount),
	}
}
