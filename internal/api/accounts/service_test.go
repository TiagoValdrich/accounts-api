package accounts_test

import (
	"context"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tiagovaldrich/accounts-api/internal/api/accounts"
	"github.com/tiagovaldrich/accounts-api/internal/models"
	"github.com/tiagovaldrich/accounts-api/internal/pkg/cerror"
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
			assert.Equal(t, &customerAccountID, createdAccount.CustomerAccount.ID)
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
				assert.Equal(t, 400, resultError.Status)
				assert.Equal(t, "Invalid document", resultError.Message)
			} else {
				t.Errorf("should have received a *cerror.Error type but received %T", err)
			}
		})
	}
}
