package accounts

import "github.com/gofrs/uuid/v5"

type CreateAccountRequest struct {
	Document string `json:"document_number" validate:"required"`
}

type SearchAccountRequest struct {
	CustomerAccountID *uuid.UUID
}
