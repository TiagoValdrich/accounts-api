package accounts

type CreateAccountRequest struct {
	Document string `json:"document_number" validate:"required"`
}
