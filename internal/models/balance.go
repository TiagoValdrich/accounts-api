package models

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/uptrace/bun"
)

type Balance struct {
	bun.BaseModel `bun:"table:balance"`

	ID                *uuid.UUID `bun:"id,pk"`
	CustomerAccountID *uuid.UUID `bun:"account_id"`
	Amount            int64      `bun:"amount"`
	CreatedAt         time.Time  `bun:"created_at"`
	UpdatedAt         time.Time  `bun:"updated_at"`
}

var _ bun.BeforeAppendModelHook = (*Balance)(nil)

func (b *Balance) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		genID, err := uuid.NewV6()
		if err != nil {
			return err
		}

		b.ID = &genID
		b.CreatedAt = time.Now()
		b.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		b.UpdatedAt = time.Now()
	}
	return nil
}
