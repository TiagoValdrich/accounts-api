package models

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/uptrace/bun"
)

type Customer struct {
	bun.BaseModel `bun:"table:customer"`
	ID            *uuid.UUID `bun:"id,pk"`
	Document      string     `bun:"document"`
	CreatedAt     time.Time  `bun:"created_at"`
	UpdatedAt     time.Time  `bun:"updated_at"`
}

var _ bun.BeforeAppendModelHook = (*Customer)(nil)

func (c *Customer) BeforeAppendModel(ctx context.Context, query bun.Query) error {
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
