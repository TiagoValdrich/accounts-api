package validator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tiagovaldrich/accounts-api/internal/pkg/validator"
)

type validStruct struct {
	Name string `validate:"required"`
}

func TestValidateStruct(t *testing.T) {
	t.Run("should return nil when struct is valid", func(t *testing.T) {
		input := validStruct{
			Name: "Titi",
		}

		err := validator.ValidateStruct(input)

		assert.Nil(t, err)
	})

	t.Run("should return field errors when required fields are missing", func(t *testing.T) {
		input := validStruct{}

		err := validator.ValidateStruct(input)

		require.NotNil(t, err)
		assert.Len(t, err.FieldErrors, 1)
	})
}
