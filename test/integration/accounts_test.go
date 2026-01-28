package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	t.Run("POST /accounts", func(t *testing.T) {
		t.Run("given a valid document should create customer, account and balance with zero value", func(t *testing.T) {
			CleanupTables(t)

			document := TestDocument
			payload := map[string]any{"document_number": document}

			resp, body := POST(t, "/accounts", payload)

			require.Equal(t, http.StatusOK, resp.StatusCode)

			var response map[string]any
			ParseJSON(t, body, &response)

			accountID := response["account_id"].(string)
			assert.NotEmpty(t, accountID)
			assert.Equal(t, document, response["document_number"])

			customer := AssertCustomerExists(t, document)

			account := AssertCustomerAccountExists(t, *customer.ID)
			assert.Equal(t, accountID, account.ID.String())

			balance := AssertBalanceExists(t, *account.ID)
			assert.Equal(t, int64(0), balance.Balance)
		})

		t.Run("given an empty payload should return bad request", func(t *testing.T) {
			CleanupTables(t)

			resp, _ := POST(t, "/accounts", map[string]any{})

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("given a missing document_number should return bad request", func(t *testing.T) {
			CleanupTables(t)

			resp, _ := POST(t, "/accounts", map[string]any{"other_field": "value"})

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("given an invalid document_number that is not a cpf or cnpj should return bad request", func(t *testing.T) {
			CleanupTables(t)

			resp, _ := POST(t, "/accounts", map[string]any{"document_number": "11122233344"})

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})
}

func TestGetAccountByID(t *testing.T) {
	t.Run("GET /accounts/:id", func(t *testing.T) {
		t.Run("with existing account should return account details", func(t *testing.T) {
			CleanupTables(t)

			document := TestDocument
			createResp, createBody := POST(t, "/accounts", map[string]any{"document_number": document})
			require.Equal(t, http.StatusOK, createResp.StatusCode)

			var createResponse map[string]any
			ParseJSON(t, createBody, &createResponse)
			accountID := createResponse["account_id"].(string)

			resp, body := GET(t, "/accounts/"+accountID)

			require.Equal(t, http.StatusOK, resp.StatusCode)

			var response map[string]any
			ParseJSON(t, body, &response)

			assert.Equal(t, accountID, response["account_id"])
			assert.Equal(t, document, response["document_number"])
		})

		t.Run("with non-existent account should return not found", func(t *testing.T) {
			CleanupTables(t)

			resp, _ := GET(t, "/accounts/00000000-0000-0000-0000-000000000000")

			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		})

		t.Run("with invalid UUID should return bad request", func(t *testing.T) {
			resp, _ := GET(t, "/accounts/invalid-uuid")

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})
}
