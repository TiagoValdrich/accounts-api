package transactions

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tiagovaldrich/accounts-api/internal/models"
	"github.com/tiagovaldrich/accounts-api/test/integration/testutils"
)

var (
	testSuite *testutils.TestSuite
)

func TestMain(m *testing.M) {
	testSuite = testutils.Setup()
	code := m.Run()
	testutils.Teardown(testSuite)
	os.Exit(code)
}

func TestCreateTransaction(t *testing.T) {
	t.Run("POST /transactions", func(t *testing.T) {
		t.Run("should create a transaction record with positive balance when the operation type is credit_voucher", func(t *testing.T) {
			testutils.CleanupTables(t, testSuite)

			accountID := createTestAccount(t, testutils.TestDocument)

			resp, body := testutils.POST(t, testSuite.App, "/transactions", map[string]any{
				"account_id":     accountID,
				"operation_type": models.CreditVoucher,
				"amount":         100.00,
			})
			require.Equal(t, http.StatusOK, resp.StatusCode)

			var response map[string]any
			testutils.ParseJSON(t, body, &response)
			assert.NotEmpty(t, response["id"])

			creditAmount := int64(10000)
			tx := AssertTransactionExists(t, accountID, models.CreditVoucher, creditAmount)
			assert.Equal(t, accountID, tx.CustomerAccountID.String())
		})

		t.Run("should return bad request when the payload is not provided", func(t *testing.T) {
			testutils.CleanupTables(t, testSuite)

			resp, _ := testutils.POST(t, testSuite.App, "/transactions", map[string]any{})

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("should return error when the account does not exists", func(t *testing.T) {
			testutils.CleanupTables(t, testSuite)

			resp, _ := testutils.POST(t, testSuite.App, "/transactions", map[string]any{
				"account_id":     "00000000-0000-0000-0000-000000000000",
				"operation_type": models.PurcharseWithInstallments,
				"amount":         50.00,
			})

			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		})
	})
}

func TestTransactionIdempotency(t *testing.T) {
	t.Run("POST /transactions", func(t *testing.T) {
		t.Run("should reject duplicated transactions(same idempotency key)", func(t *testing.T) {
			testutils.CleanupTables(t, testSuite)

			accountID := createTestAccount(t, testutils.TestDocument)
			idempotencyKey := "unique-key-123"

			payload := map[string]any{
				"account_id":      accountID,
				"operation_type":  models.CreditVoucher,
				"amount":          200.00,
				"idempotency_key": idempotencyKey,
			}

			firstResp, _ := testutils.POST(t, testSuite.App, "/transactions", payload)
			require.Equal(t, http.StatusOK, firstResp.StatusCode)

			tx := AssertTransactionExistsWithIdempotencyKey(t, idempotencyKey)
			assert.Equal(t, accountID, tx.CustomerAccountID.String())
			assert.Equal(t, models.CreditVoucher, tx.OperationType)
			assert.Equal(t, int64(20000), tx.Amount)

			secondResp, _ := testutils.POST(t, testSuite.App, "/transactions", payload)
			assert.Equal(t, http.StatusConflict, secondResp.StatusCode)

			assert.Equal(t, 1, CountTransactionsWithIdempotencyKey(t, idempotencyKey))
		})

		t.Run("should be able to create multiple transactions with different idempotency keys", func(t *testing.T) {
			testutils.CleanupTables(t, testSuite)

			accountID := createTestAccount(t, testutils.TestDocument)

			firstResp, _ := testutils.POST(t, testSuite.App, "/transactions", map[string]any{
				"account_id":      accountID,
				"operation_type":  models.CreditVoucher,
				"amount":          100.00,
				"idempotency_key": "key-1",
			})
			require.Equal(t, http.StatusOK, firstResp.StatusCode)

			secondResp, _ := testutils.POST(t, testSuite.App, "/transactions", map[string]any{
				"account_id":      accountID,
				"operation_type":  models.CreditVoucher,
				"amount":          100.00,
				"idempotency_key": "key-2",
			})
			require.Equal(t, http.StatusOK, secondResp.StatusCode)

			assert.Equal(t, 1, CountTransactionsWithIdempotencyKey(t, "key-1"))
			assert.Equal(t, 1, CountTransactionsWithIdempotencyKey(t, "key-2"))
			assert.Equal(t, 2, CountTransactionsForAccount(t, accountID))
		})
	})
}
