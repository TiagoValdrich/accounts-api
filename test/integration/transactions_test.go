package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tiagovaldrich/accounts-api/internal/models"
)

func TestCreateTransaction(t *testing.T) {
	t.Run("POST /transactions", func(t *testing.T) {
		t.Run("credit voucher should increase balance and create transaction record", func(t *testing.T) {
			CleanupTables(t)

			accountID := createTestAccount(t, TestDocument)

			AssertBalanceEquals(t, accountID, 0)

			resp, body := POST(t, "/transactions", map[string]any{
				"account_id":     accountID,
				"operation_type": models.CreditVoucher,
				"amount":         100.00,
			})
			require.Equal(t, http.StatusOK, resp.StatusCode)

			var response map[string]any
			ParseJSON(t, body, &response)
			assert.NotEmpty(t, response["id"])

			creditAmount := int64(10000)
			tx := AssertTransactionExists(t, accountID, models.CreditVoucher, creditAmount)
			assert.Equal(t, accountID, tx.CustomerAccountID.String())

			AssertBalanceEquals(t, accountID, creditAmount)
		})

		t.Run("purchase should decrease balance and create transaction record", func(t *testing.T) {
			CleanupTables(t)

			accountID := createTestAccount(t, TestDocument)

			creditResp, _ := POST(t, "/transactions", map[string]any{
				"account_id":     accountID,
				"operation_type": models.CreditVoucher,
				"amount":         100.00,
			})
			require.Equal(t, http.StatusOK, creditResp.StatusCode)

			initialBalance := GetBalance(t, accountID)

			purchaseResp, body := POST(t, "/transactions", map[string]any{
				"account_id":     accountID,
				"operation_type": models.NormalPurchase,
				"amount":         50.00,
			})
			require.Equal(t, http.StatusOK, purchaseResp.StatusCode)

			var response map[string]any
			ParseJSON(t, body, &response)
			assert.NotEmpty(t, response["id"])

			purchaseAmount := int64(5000) * -1
			AssertTransactionExists(t, accountID, models.NormalPurchase, purchaseAmount)

			AssertBalanceEquals(t, accountID, initialBalance+purchaseAmount)

			assert.Equal(t, 2, CountTransactionsForAccount(t, accountID))
		})

		t.Run("with insufficient funds should reject and not create transaction", func(t *testing.T) {
			CleanupTables(t)

			accountID := createTestAccount(t, TestDocument)

			AssertBalanceEquals(t, accountID, 0)

			resp, _ := POST(t, "/transactions", map[string]any{
				"account_id":     accountID,
				"operation_type": models.Withdrawal,
				"amount":         100.00,
			})

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

			assert.Equal(t, 0, CountTransactionsForAccount(t, accountID))

			AssertBalanceEquals(t, accountID, 0)
		})

		t.Run("with empty payload should return bad request", func(t *testing.T) {
			CleanupTables(t)

			resp, _ := POST(t, "/transactions", map[string]any{})

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("with non-existent account should return error", func(t *testing.T) {
			CleanupTables(t)

			resp, _ := POST(t, "/transactions", map[string]any{
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
		t.Run("duplicate request with same idempotency key should be rejected and only one transaction created", func(t *testing.T) {
			CleanupTables(t)

			accountID := createTestAccount(t, TestDocument)
			idempotencyKey := "unique-key-123"

			payload := map[string]any{
				"account_id":      accountID,
				"operation_type":  models.CreditVoucher,
				"amount":          200.00,
				"idempotency_key": idempotencyKey,
			}

			firstResp, _ := POST(t, "/transactions", payload)
			require.Equal(t, http.StatusOK, firstResp.StatusCode)

			tx := AssertTransactionExistsWithIdempotencyKey(t, idempotencyKey)
			assert.Equal(t, accountID, tx.CustomerAccountID.String())
			assert.Equal(t, models.CreditVoucher, tx.OperationType)
			assert.Equal(t, int64(20000), tx.Amount)

			balanceAfterFirst := GetBalance(t, accountID)

			secondResp, _ := POST(t, "/transactions", payload)
			assert.Equal(t, http.StatusConflict, secondResp.StatusCode)

			assert.Equal(t, 1, CountTransactionsWithIdempotencyKey(t, idempotencyKey))

			AssertBalanceEquals(t, accountID, balanceAfterFirst)
		})

		t.Run("different idempotency keys should create separate transactions", func(t *testing.T) {
			CleanupTables(t)

			accountID := createTestAccount(t, TestDocument)

			firstResp, _ := POST(t, "/transactions", map[string]any{
				"account_id":      accountID,
				"operation_type":  models.CreditVoucher,
				"amount":          100.00,
				"idempotency_key": "key-1",
			})
			require.Equal(t, http.StatusOK, firstResp.StatusCode)

			secondResp, _ := POST(t, "/transactions", map[string]any{
				"account_id":      accountID,
				"operation_type":  models.CreditVoucher,
				"amount":          100.00,
				"idempotency_key": "key-2",
			})
			require.Equal(t, http.StatusOK, secondResp.StatusCode)

			assert.Equal(t, 1, CountTransactionsWithIdempotencyKey(t, "key-1"))
			assert.Equal(t, 1, CountTransactionsWithIdempotencyKey(t, "key-2"))
			assert.Equal(t, 2, CountTransactionsForAccount(t, accountID))

			AssertBalanceEquals(t, accountID, 20000)
		})
	})
}
