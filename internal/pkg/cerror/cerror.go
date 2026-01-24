package cerror

type (
	FieldError struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}

	Params struct {
		Status  int
		Message string
	}

	Error struct {
		Status      int          `json:"status"`
		Message     string       `json:"message"`
		FieldErrors []FieldError `json:"field_errors,omitempty"`
	}
)

func New(params Params, validationErrors ...FieldError) *Error {
	return &Error{
		Status:      params.Status,
		Message:     params.Message,
		FieldErrors: validationErrors,
	}
}

func (e *Error) Error() string {
	return e.Message
}
