package transactions

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tiagovaldrich/accounts-api/internal/models"
	"github.com/tiagovaldrich/accounts-api/test/integration/testutils"
)

func createTestAccount(t *testing.T, document string) string {
	t.Helper()

	resp, body := testutils.POST(t, testSuite.App, "/accounts", map[string]any{"document_number": document})
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]any
	testutils.ParseJSON(t, body, &response)

	return response["account_id"].(string)
}

func AssertTransactionExists(t *testing.T, accountID string, operationType models.OperationType, amount int64) models.Transaction {
	t.Helper()

	var transaction models.Transaction
	err := testSuite.DB.NewSelect().
		Model(&transaction).
		Where("customer_account_id = ?", accountID).
		Where("operation_type = ?", operationType).
		Where("amount = ?", amount).
		Scan(context.Background())

	require.NoError(t, err)
	assert.NotNil(t, transaction.ID)

	return transaction
}

func AssertTransactionExistsWithIdempotencyKey(t *testing.T, idempotencyKey string) models.Transaction {
	t.Helper()

	var transaction models.Transaction
	err := testSuite.DB.NewSelect().
		Model(&transaction).
		Where("idempotency_key = ?", idempotencyKey).
		Scan(context.Background())

	require.NoError(t, err)
	assert.NotNil(t, transaction.ID)

	return transaction
}

func CountTransactionsWithIdempotencyKey(t *testing.T, idempotencyKey string) int {
	t.Helper()

	count, err := testSuite.DB.NewSelect().
		Model((*models.Transaction)(nil)).
		Where("idempotency_key = ?", idempotencyKey).
		Count(context.Background())

	require.NoError(t, err)
	return count
}

func CountTransactionsForAccount(t *testing.T, accountID string) int {
	t.Helper()

	count, err := testSuite.DB.NewSelect().
		Model((*models.Transaction)(nil)).
		Where("customer_account_id = ?", accountID).
		Count(context.Background())

	require.NoError(t, err)
	return count
}
