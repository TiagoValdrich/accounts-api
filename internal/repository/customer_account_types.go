package repository

import (
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/uptrace/bun"
)

type CustomerAccountByIDResult struct {
	bun.BaseModel `bun:"table:customer_account"`
	ID            *uuid.UUID `bun:"id"`
	Document      string     `bun:"document"`
	CreatedAt     time.Time  `bun:"created_at"`
}
