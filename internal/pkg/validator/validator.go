package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/tiagovaldrich/accounts-api/internal/pkg/cerror"
)

var validate *validator.Validate

func init() {
	if validate == nil {
		validate = validator.New(validator.WithRequiredStructEnabled())
	}
}

func ValidateStruct(value any) *cerror.Error {
	var errors []cerror.FieldError

	if err := validate.Struct(value); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return cerror.New(cerror.Params{
				Message: "invalid struct",
			})
		}

		for _, validationErr := range err.(validator.ValidationErrors) {
			errors = append(errors, cerror.FieldError{
				Field:   validationErr.Field(),
				Message: validationErr.Error(),
			})
		}

		return cerror.New(cerror.Params{}, errors...)
	}

	return nil
}
