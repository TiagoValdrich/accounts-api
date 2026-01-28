package models

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/uptrace/bun"
)

type OperationType string

const (
	NormalPurchase            OperationType = "normal_purchase"
	PurcharseWithInstallments OperationType = "installment_purchase"
	Withdrawal                OperationType = "withdrawal"
	CreditVoucher             OperationType = "credit_voucher"
)

type Transaction struct {
	bun.BaseModel     `bun:"table:transactions"`
	ID                *uuid.UUID    `bun:"id,pk"`
	CustomerAccountID *uuid.UUID    `bun:"customer_account_id"`
	OperationType     OperationType `bun:"operation_type"`
	Amount            int64         `bun:"amount"`
	IdempotencyKey    *string       `bun:"idempotency_key"`
	CreatedAt         time.Time     `bun:"created_at"`
	UpdatedAt         time.Time     `bun:"updated_at"`
}

var _ bun.BeforeAppendModelHook = (*Transaction)(nil)

func (t *Transaction) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		genID, err := uuid.NewV6()
		if err != nil {
			return err
		}

		t.ID = &genID
		t.CreatedAt = time.Now()
		t.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		t.UpdatedAt = time.Now()
	}
	return nil
}
