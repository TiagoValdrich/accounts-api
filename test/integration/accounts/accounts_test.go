package accounts

import (
	"net/http"
	"os"
	"testing"

	"github.com/tiagovaldrich/accounts-api/test/integration/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestCreateAccount(t *testing.T) {
	t.Run("POST /accounts", func(t *testing.T) {
		t.Run("should create a customer account when provided a valid document", func(t *testing.T) {
			testutils.CleanupTables(t, testSuite)

			document := testutils.TestDocument
			payload := map[string]any{"document_number": document}

			resp, body := testutils.POST(t, testSuite.App, "/accounts", payload)

			require.Equal(t, http.StatusOK, resp.StatusCode)

			var response map[string]any
			testutils.ParseJSON(t, body, &response)

			accountID := response["account_id"].(string)
			assert.NotEmpty(t, accountID)
			assert.Equal(t, document, response["document_number"])

			customer := AssertCustomerExists(t, document)

			account := AssertCustomerAccountExists(t, *customer.ID)
			assert.Equal(t, accountID, account.ID.String())
		})

		t.Run("should return bad request when provided an empty payload ", func(t *testing.T) {
			testutils.CleanupTables(t, testSuite)

			resp, _ := testutils.POST(t, testSuite.App, "/accounts", map[string]any{})

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("should return bad request when document_number is missing", func(t *testing.T) {
			testutils.CleanupTables(t, testSuite)

			resp, _ := testutils.POST(t, testSuite.App, "/accounts", map[string]any{"other_field": "value"})

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("should return bad request when provided a document_number that is not a cpf or cnpj", func(t *testing.T) {
			testutils.CleanupTables(t, testSuite)

			resp, _ := testutils.POST(t, testSuite.App, "/accounts", map[string]any{"document_number": "11122233344"})

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})
}

func TestGetAccountByID(t *testing.T) {
	t.Run("GET /accounts/:id", func(t *testing.T) {
		t.Run("should return account details", func(t *testing.T) {
			testutils.CleanupTables(t, testSuite)

			document := testutils.TestDocument
			createResp, createBody := testutils.POST(t, testSuite.App, "/accounts", map[string]any{"document_number": document})
			require.Equal(t, http.StatusOK, createResp.StatusCode)

			var createResponse map[string]any
			testutils.ParseJSON(t, createBody, &createResponse)
			accountID := createResponse["account_id"].(string)

			resp, body := testutils.GET(t, testSuite.App, "/accounts/"+accountID)

			require.Equal(t, http.StatusOK, resp.StatusCode)

			var response map[string]any
			testutils.ParseJSON(t, body, &response)

			assert.Equal(t, accountID, response["account_id"])
			assert.Equal(t, document, response["document_number"])
		})

		t.Run("should return not found when the account is not found", func(t *testing.T) {
			testutils.CleanupTables(t, testSuite)

			resp, _ := testutils.GET(t, testSuite.App, "/accounts/00000000-0000-0000-0000-000000000000")

			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		})

		t.Run("should return bad request when provided an invalid account ID", func(t *testing.T) {
			resp, _ := testutils.GET(t, testSuite.App, "/accounts/invalid-uuid")

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})
}
