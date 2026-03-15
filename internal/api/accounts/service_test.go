package accounts_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tiagovaldrich/accounts-api/internal/api/accounts"
	"github.com/tiagovaldrich/accounts-api/internal/models"
	"github.com/tiagovaldrich/accounts-api/internal/pkg/cerror"
	"github.com/tiagovaldrich/accounts-api/internal/repository"
	repoMock "github.com/tiagovaldrich/accounts-api/test/mock/repository"
	"go.uber.org/mock/gomock"
)

func TestCreateAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customerRepoMock := repoMock.NewMockCustomerRepository(ctrl)
	customerAccountRepoMock := repoMock.NewMockCustomerAccountRepository(ctrl)
	service := accounts.NewService(customerRepoMock, customerAccountRepoMock)

	customerRepoMock.EXPECT().
		WithTransaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(txCtx context.Context) error) error {
			return fn(ctx)
		}).
		AnyTimes()

	validDocumentsToTest := []string{
		"76713659047",
		"551.048.870-06",
		"36.050.352/0001-07",
		"73980588000160",
	}

	for _, documentNumber := range validDocumentsToTest {
		t.Run("should create a customer and customer account when provided a valid document number", func(t *testing.T) {
			customerID, err := uuid.NewV6()
			require.Nil(t, err)
			customerAccountID, err := uuid.NewV6()
			require.Nil(t, err)
			expectedTime := time.Date(2026, time.March, 25, 0, 0, 0, 0, time.UTC)

			customerRepoMock.EXPECT().
				CreateCustomer(gomock.Any(), models.Customer{
					Document: documentNumber,
				}).
				AnyTimes().
				Return(&models.Customer{
					ID:        &customerID,
					Document:  documentNumber,
					CreatedAt: expectedTime,
				}, nil)

			customerAccountRepoMock.EXPECT().
				CreateCustomerAccount(gomock.Any(), models.CustomerAccount{
					CustomerID: &customerID,
				}).
				AnyTimes().
				Return(&models.CustomerAccount{
					ID: &customerAccountID,
				}, nil)

			createdAccount, err := service.CreateAccount(t.Context(), accounts.CreateAccountRequest{
				Document: documentNumber,
			})

			assert.Nil(t, err)
			assert.Equal(t, documentNumber, createdAccount.Customer.Document)
			assert.Equal(t, expectedTime, createdAccount.Customer.CreatedAt)
			assert.Equal(t, customerAccountID.String(), createdAccount.CustomerAccount.ID.String())
		})
	}
}

func TestInvalidDocumentNumberOnCreateAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customerRepoMock := repoMock.NewMockCustomerRepository(ctrl)
	customerAccountRepoMock := repoMock.NewMockCustomerAccountRepository(ctrl)
	service := accounts.NewService(customerRepoMock, customerAccountRepoMock)

	wrongDocumentList := []string{
		"12345678910",
		"123.456.789-10",
		"blablabla",
		"123blabla123123",
		"",
	}

	for _, wrongDoc := range wrongDocumentList {
		t.Run("should return an invalid document error when provided an invalid document number", func(t *testing.T) {
			_, err := service.CreateAccount(t.Context(), accounts.CreateAccountRequest{
				Document: wrongDoc,
			})

			require.NotNil(t, err, "Should raise an error when document is invalid")

			var resultError *cerror.Error
			if assert.ErrorAs(t, err, &resultError) {
				assert.Equal(t, http.StatusBadRequest, resultError.Status)
				assert.Equal(t, "Invalid document", resultError.Message)
			} else {
				t.Errorf("should have received a *cerror.Error type but received %T", err)
			}
		})
	}
}

func TestSearchCustomerByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customerRepoMock := repoMock.NewMockCustomerRepository(ctrl)
	customerAccountRepoMock := repoMock.NewMockCustomerAccountRepository(ctrl)
	service := accounts.NewService(customerRepoMock, customerAccountRepoMock)

	expectedCustomerAccountID, err := uuid.NewV6()
	require.Nil(t, err)
	expectedDocument := "76713659047"
	expectedTime := time.Date(2026, time.March, 25, 0, 0, 0, 0, time.UTC)

	customerAccountRepoMock.EXPECT().
		SearchCustomerAccountByID(gomock.Any(), &expectedCustomerAccountID).
		Times(1).
		Return(&repository.CustomerAccountByIDResult{
			ID:        &expectedCustomerAccountID,
			Document:  expectedDocument,
			CreatedAt: expectedTime,
		}, nil)

	t.Run("should return a customer account when provided a valid customer account id", func(t *testing.T) {
		customerAccount, err := service.SearchCustomerAccountByID(t.Context(), accounts.SearchAccountRequest{
			CustomerAccountID: &expectedCustomerAccountID,
		})

		assert.Nil(t, err)
		assert.Equal(t, expectedCustomerAccountID.String(), customerAccount.CustomerID.String())
		assert.Equal(t, expectedDocument, customerAccount.Document)
	})
}

func TestSearchCustomerByIDNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customerRepoMock := repoMock.NewMockCustomerRepository(ctrl)
	customerAccountRepoMock := repoMock.NewMockCustomerAccountRepository(ctrl)
	service := accounts.NewService(customerRepoMock, customerAccountRepoMock)

	expectedCustomerAccountID, err := uuid.NewV6()
	require.Nil(t, err)

	customerAccountRepoMock.EXPECT().
		SearchCustomerAccountByID(gomock.Any(), &expectedCustomerAccountID).
		Times(1).
		Return(nil, nil)

	t.Run("should return an error with not found status", func(t *testing.T) {
		_, err := service.SearchCustomerAccountByID(t.Context(), accounts.SearchAccountRequest{
			CustomerAccountID: &expectedCustomerAccountID,
		})

		require.NotNil(t, err)

		var resultError *cerror.Error
		if assert.ErrorAs(t, err, &resultError) {
			assert.Equal(t, http.StatusNotFound, resultError.Status)
			assert.Equal(t, "Customer account not found", resultError.Message)
		} else {
			t.Errorf("should have received a *cerror.Error type but received %T", err)
		}
	})
}
