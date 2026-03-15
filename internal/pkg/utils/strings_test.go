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
