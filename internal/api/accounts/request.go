package accounts

import "github.com/gofrs/uuid/v5"

type createAccountRequest struct {
	Document string `json:"document_number" validate:"required"`
}

type searchAccountRequest struct {
	CustomerAccountID *uuid.UUID
}
