package utils_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tiagovaldrich/accounts-api/internal/models"
	"github.com/tiagovaldrich/accounts-api/internal/pkg/cerror"
	"github.com/tiagovaldrich/accounts-api/internal/pkg/utils"
)

func TestToCents(t *testing.T) {
	t.Run("should convert a float value to cents", func(t *testing.T) {
		inputValue := 15.67
		expectedValue := int64(1567)
		parsedValue := utils.ToCents(inputValue)

		assert.Equal(t, expectedValue, parsedValue)
	})
}

func TestFromCents(t *testing.T) {
	t.Run("should convert a value that in cents to float value", func(t *testing.T) {
		inputValue := int64(1567)
		expectedValue := 15.67
		parsedValue := utils.FromCents(inputValue)

		assert.Equal(t, expectedValue, parsedValue)
	})
}

func TestApplyMoneyDirection(t *testing.T) {
	type testTable struct {
		description   string
		operationType models.OperationType
		inputValue    int64
		expectedValue int64
	}

	testCases := []testTable{
		{
			description:   "credit voucher operations should return positive values",
			operationType: models.CreditVoucher,
			inputValue:    int64(10000),
			expectedValue: int64(10000),
		},
		{
			description:   "normal purchase operations should return negative values",
			operationType: models.NormalPurchase,
			inputValue:    int64(10000),
			expectedValue: int64(-10000),
		},
		{
			description:   "installment purchase operations should return negative values",
			operationType: models.PurcharseWithInstallments,
			inputValue:    int64(10000),
			expectedValue: int64(-10000),
		},
		{
			description:   "withdrawal operations should return negative values",
			operationType: models.Withdrawal,
			inputValue:    int64(10000),
			expectedValue: int64(-10000),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			result, err := utils.ApplyMoneyDirection(tc.inputValue, tc.operationType)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedValue, result)
		})
	}

	t.Run("should return unsupported operation error if operation is not mapped", func(t *testing.T) {
		inputValue := int64(10000)
		operationNotExists := models.OperationType("random_op")

		result, err := utils.ApplyMoneyDirection(inputValue, operationNotExists)

		require.NotNil(t, err)
		assert.Equal(t, int64(0), result)

		var outputError *cerror.Error
		if assert.ErrorAs(t, err, &outputError) {
			assert.Equal(t, http.StatusInternalServerError, outputError.Status)
			assert.Equal(t, "Unsupported operation", outputError.Message)
		} else {
			t.Errorf("should have received a *cerror.Error type but received %T", err)
		}
	})
}
