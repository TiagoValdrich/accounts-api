package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tiagovaldrich/accounts-api/internal/pkg/utils"
)

func TestSafeStringPointerValue(t *testing.T) {
	t.Run("should return empty string when pointer is nil", func(t *testing.T) {
		result := utils.SafeStringPointerValue(nil)

		assert.Equal(t, "", result)
	})

	t.Run("should return the string value when pointer is not nil", func(t *testing.T) {
		value := "teste"
		result := utils.SafeStringPointerValue(&value)

		assert.Equal(t, "teste", result)
	})
}

func TestGetStringPointer(t *testing.T) {
	t.Run("should return a pointer to the provided string", func(t *testing.T) {
		result := utils.GetStringPointer("hello")

		assert.NotNil(t, result)
		assert.Equal(t, "hello", *result)
	})
}

func TestSanitizeDocumentNumber(t *testing.T) {
	t.Run("should return only digits from a formatted CPF", func(t *testing.T) {
		result := utils.SanitizeDocumentNumber("551.048.870-06")

		assert.Equal(t, "55104887006", result)
	})

	t.Run("should return only digits from a formatted CNPJ", func(t *testing.T) {
		result := utils.SanitizeDocumentNumber("36.050.352/0001-07")

		assert.Equal(t, "36050352000107", result)
	})

	t.Run("should return the same string when it only contains digits", func(t *testing.T) {
		result := utils.SanitizeDocumentNumber("55104887006")

		assert.Equal(t, "55104887006", result)
	})

	t.Run("should return empty string when input is empty", func(t *testing.T) {
		result := utils.SanitizeDocumentNumber("")

		assert.Equal(t, "", result)
	})

	t.Run("should return empty string when input has no digits", func(t *testing.T) {
		result := utils.SanitizeDocumentNumber("abcdef!@#")

		assert.Equal(t, "", result)
	})

	t.Run("should strip spaces and letters mixed with digits", func(t *testing.T) {
		result := utils.SanitizeDocumentNumber("abc 123 def 456")

		assert.Equal(t, "123456", result)
	})
}
