package transactions_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tiagovaldrich/accounts-api/internal/api/transactions"
	"github.com/tiagovaldrich/accounts-api/internal/models"
	"github.com/tiagovaldrich/accounts-api/internal/pkg/cerror"
	"github.com/tiagovaldrich/accounts-api/internal/pkg/utils"
	"github.com/tiagovaldrich/accounts-api/internal/repository"
	repoMock "github.com/tiagovaldrich/accounts-api/internal/repository/mock"
	"go.uber.org/mock/gomock"
)

func TestCreateTransaction(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	transactionRepoMock := repoMock.NewMockTransactionRepository(mockCtrl)
	customerAccountRepoMock := repoMock.NewMockCustomerAccountRepository(mockCtrl)

	customerAccountID, err := uuid.NewV6()
	require.Nil(t, err)

	transactionsToTest := []transactions.CreateTransactionRequest{
		{
			CustomerAccountID: &customerAccountID,
			OperationType:     models.CreditVoucher,
			Amount:            400.00,
			IdempotencyKey:    utils.GetStringPointer("id-1"),
		},
		{
			CustomerAccountID: &customerAccountID,
			OperationType:     models.NormalPurchase,
			Amount:            100.00,
			IdempotencyKey:    utils.GetStringPointer("id-2"),
		},
		{
			CustomerAccountID: &customerAccountID,
			OperationType:     models.PurchaseWithInstallments,
			Amount:            50.00,
			IdempotencyKey:    utils.GetStringPointer("id-3"),
		},
		{
			CustomerAccountID: &customerAccountID,
			OperationType:     models.Withdrawal,
			Amount:            25.00,
			IdempotencyKey:    utils.GetStringPointer("id-4"),
		},
	}

	transactionRepoMock.EXPECT().
		GetTransactionByIdempotencyKey(gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(nil, nil)

	customerAccountRepoMock.EXPECT().
		SearchCustomerAccountByID(gomock.Any(), &customerAccountID).
		AnyTimes().
		Return(&repository.CustomerAccountByIDResult{
			ID: &customerAccountID,
		}, nil)

	transactionRepoMock.EXPECT().
		WithTransaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(txCtx context.Context) error) error {
			return fn(ctx)
		}).
		AnyTimes()

	for _, transaction := range transactionsToTest {
		t.Run("should create a transaction successfully when provided all the required data", func(t *testing.T) {
			transactionRepoMock.EXPECT().
				CreateTransaction(gomock.Any(), gomock.Any()).
				Times(1).
				DoAndReturn(func(ctx context.Context, transactionModel models.Transaction) (*models.Transaction, error) {
					expectedAmount := utils.ToCents(transaction.Amount)
					if transactionModel.OperationType != models.CreditVoucher {
						expectedAmount = expectedAmount * -1
					}

					assert.Equal(t, expectedAmount, transactionModel.Amount)
					assert.Equal(t, transaction.IdempotencyKey, transactionModel.IdempotencyKey)
					assert.Equal(t, transaction.CustomerAccountID.String(), transactionModel.CustomerAccountID.String())
					assert.Equal(t, transaction.OperationType, transactionModel.OperationType)

					return &transactionModel, nil
				})

			service := transactions.NewService(transactionRepoMock, customerAccountRepoMock)

			_, err := service.CreateTransaction(t.Context(), transaction)

			require.Nil(t, err)
		})
	}
}

func TestCreateTransactionDuplicatedIdempotencyKey(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	transactionRepoMock := repoMock.NewMockTransactionRepository(mockCtrl)
	customerAccountRepoMock := repoMock.NewMockCustomerAccountRepository(mockCtrl)

	customerAccountID, err := uuid.NewV6()
	require.Nil(t, err)

	transactionRepoMock.EXPECT().
		GetTransactionByIdempotencyKey(gomock.Any(), "duplicated-key").
		Times(1).
		Return(&models.Transaction{
			CustomerAccountID: &customerAccountID,
			OperationType:     models.NormalPurchase,
			Amount:            -10000,
			IdempotencyKey:    utils.GetStringPointer("duplicated-key"),
		}, nil)

	service := transactions.NewService(transactionRepoMock, customerAccountRepoMock)

	t.Run("should return a conflict error when a transaction has a duplicated idempotency key", func(t *testing.T) {
		_, err := service.CreateTransaction(t.Context(), transactions.CreateTransactionRequest{
			CustomerAccountID: &customerAccountID,
			OperationType:     models.NormalPurchase,
			Amount:            100.00,
			IdempotencyKey:    utils.GetStringPointer("duplicated-key"),
		})

		require.NotNil(t, err)

		var resultError *cerror.Error
		if assert.ErrorAs(t, err, &resultError) {
			assert.Equal(t, http.StatusConflict, resultError.Status)
			assert.Equal(t, "Transaction is already created with that idempotency key", resultError.Message)
		} else {
			t.Errorf("should have received a *cerror.Error type but received %T", err)
		}
	})
}

func TestCreateTransactionCustomerAccountNotFound(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	transactionRepoMock := repoMock.NewMockTransactionRepository(mockCtrl)
	customerAccountRepoMock := repoMock.NewMockCustomerAccountRepository(mockCtrl)

	customerAccountID, err := uuid.NewV6()
	require.Nil(t, err)

	transactionRepoMock.EXPECT().
		GetTransactionByIdempotencyKey(gomock.Any(), "id-1").
		Times(1).
		Return(nil, nil)

	customerAccountRepoMock.EXPECT().
		SearchCustomerAccountByID(gomock.Any(), &customerAccountID).
		Times(1).
		Return(nil, nil)

	service := transactions.NewService(transactionRepoMock, customerAccountRepoMock)

	t.Run("should return a not found error when customer account does not exist", func(t *testing.T) {
		_, err := service.CreateTransaction(t.Context(), transactions.CreateTransactionRequest{
			CustomerAccountID: &customerAccountID,
			OperationType:     models.NormalPurchase,
			Amount:            100.00,
			IdempotencyKey:    utils.GetStringPointer("id-1"),
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
