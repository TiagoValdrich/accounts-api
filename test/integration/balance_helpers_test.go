package integration

import (
	"context"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tiagovaldrich/accounts-api/internal/models"
)

func AssertBalanceExists(t *testing.T, accountID uuid.UUID) models.Balance {
	t.Helper()

	var balance models.Balance
	err := DB.NewSelect().
		Model(&balance).
		Where("customer_account_id = ?", accountID).
		Scan(context.Background())

	require.NoError(t, err, "balance for account %s should exist", accountID)
	assert.NotNil(t, balance.ID)

	return balance
}

func AssertBalanceEquals(t *testing.T, accountID string, expectedBalance int64) {
	t.Helper()

	var balance models.Balance
	err := DB.NewSelect().
		Model(&balance).
		Where("customer_account_id = ?", accountID).
		Scan(context.Background())

	require.NoError(t, err, "balance for account %s should exist", accountID)
	assert.Equal(t, expectedBalance, balance.Balance)
}

func GetBalance(t *testing.T, accountID string) int64 {
	t.Helper()

	var balance models.Balance
	err := DB.NewSelect().
		Model(&balance).
		Where("customer_account_id = ?", accountID).
		Scan(context.Background())

	require.NoError(t, err, "balance for account %s should exist", accountID)
	return balance.Balance
}
