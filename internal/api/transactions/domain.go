package transactions

import (
	"github.com/gofrs/uuid/v5"
	"github.com/tiagovaldrich/accounts-api/internal/models"
)

type CreateTransactionResult struct {
	ID                *uuid.UUID
	CustomerAccountID *uuid.UUID
	OperationType     models.OperationType
	Amount            int64
}
