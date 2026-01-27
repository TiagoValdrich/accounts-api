package transactions

import (
	"github.com/gofrs/uuid/v5"
	"github.com/tiagovaldrich/accounts-api/internal/models"
)

type createTransactionRequest struct {
	CustomerAccountID *uuid.UUID           `json:"account_id" validate:"required"`
	OperationType     models.OperationType `json:"operation_type" validate:"required,oneof=normal_purchase installment_purchase withdrawal credit_voucher"`
	Amount            float64              `json:"amount" validate:"required,gt=0"`
	IdempotencyKey    *string              `json:"idempotency_key" validate:"omitempty"`
}
