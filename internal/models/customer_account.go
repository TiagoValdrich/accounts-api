package models

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/uptrace/bun"
)

type CustomerAccount struct {
	bun.BaseModel `bun:"table:customer_account"`
	ID            *uuid.UUID `bun:"id,pk"`
	CustomerID    *uuid.UUID `bun:"customer_id"`
	CreatedAt     time.Time  `bun:"created_at"`
	UpdatedAt     time.Time  `bun:"updated_at"`
}

var _ bun.BeforeAppendModelHook = (*CustomerAccount)(nil)

func (c *CustomerAccount) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		genID, err := uuid.NewV6()
		if err != nil {
			return err
		}

		c.ID = &genID
		c.CreatedAt = time.Now()
		c.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		c.UpdatedAt = time.Now()
	}
	return nil
}
