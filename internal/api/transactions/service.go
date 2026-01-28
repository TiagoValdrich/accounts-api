package transactions

import (
	"context"
	"net/http"

	"github.com/gofrs/uuid/v5"
	"github.com/rs/zerolog/log"
	"github.com/tiagovaldrich/accounts-api/internal/models"
	"github.com/tiagovaldrich/accounts-api/internal/pkg/cerror"
	"github.com/tiagovaldrich/accounts-api/internal/pkg/utils"
	"github.com/tiagovaldrich/accounts-api/internal/repository"
)

type Servicer interface {
	CreateTransaction(context.Context, createTransactionRequest) (CreateTransactionResult, error)
}

type service struct {
	transactionRepository     repository.TransactionRepository
	customerAccountRepository repository.CustomerAccountRepository
	balanceRepository         repository.BalanceRepository
}

func NewService(
	transactionRepository repository.TransactionRepository,
	customerAccountRepository repository.CustomerAccountRepository,
	balanceRepository repository.BalanceRepository,
) Servicer {
	return &service{
		transactionRepository:     transactionRepository,
		customerAccountRepository: customerAccountRepository,
		balanceRepository:         balanceRepository,
	}
}

func (s *service) CreateTransaction(ctx context.Context, request createTransactionRequest) (CreateTransactionResult, error) {
	if err := s.validateIdempotency(ctx, request); err != nil {
		return CreateTransactionResult{}, err
	}

	customerAccount, err := s.findCustomerAccount(ctx, request.CustomerAccountID)
	if err != nil {
		return CreateTransactionResult{}, err
	}

	transaction, err := s.processTransaction(ctx, customerAccount, request)
	if err != nil {
		return CreateTransactionResult{}, err
	}

	return CreateTransactionResult{
		ID:                transaction.ID,
		CustomerAccountID: transaction.CustomerAccountID,
		OperationType:     transaction.OperationType,
		Amount:            transaction.Amount,
	}, nil
}

func (s *service) validateIdempotency(ctx context.Context, request createTransactionRequest) error {
	if request.IdempotencyKey == nil || *request.IdempotencyKey == "" {
		return nil
	}

	transaction, err := s.transactionRepository.GetTransactionByIdempotencyKey(ctx, *request.IdempotencyKey)
	if err != nil {
		log.Err(err).
			Str("idempotency_key", *request.IdempotencyKey).
			Msg("failed to check transaction by idempotency key")

		return err
	}

	if transaction != nil {
		return cerror.New(cerror.Params{
			Status:  http.StatusConflict,
			Message: "Transaction is already created with that idempotency key",
		})
	}

	return nil
}

func (s *service) findCustomerAccount(ctx context.Context, customerAccountID *uuid.UUID) (*repository.CustomerAccountByIDResult, error) {
	customerAccount, err := s.customerAccountRepository.SearchCustomerAccountByID(ctx, customerAccountID)
	if err != nil {
		log.Err(err).
			Str("customer_account_id", customerAccountID.String()).
			Msg("failed to find customer account")

		return nil, cerror.New(cerror.Params{
			Status:  http.StatusNotFound,
			Message: "Customer account not found",
		})
	}

	if customerAccount == nil {
		return nil, cerror.New(cerror.Params{
			Status:  http.StatusNotFound,
			Message: "Customer account not found",
		})
	}

	return customerAccount, nil
}

func (s *service) processTransaction(
	ctx context.Context,
	customerAccount *repository.CustomerAccountByIDResult,
	request createTransactionRequest,
) (*models.Transaction, error) {
	var transactionCreated *models.Transaction

	err := s.transactionRepository.WithTransaction(ctx, func(txCtx context.Context) error {
		accountBalance, err := s.getAccountBalance(txCtx, customerAccount.ID, request)
		if err != nil {
			return err
		}

		amountCents, err := s.calculateTransactionAmount(request, accountBalance.Balance)
		if err != nil {
			return err
		}

		transactionCreated, err = s.createTransaction(txCtx, customerAccount, request, amountCents)
		if err != nil {
			return err
		}

		if err := s.updateBalance(txCtx, customerAccount, accountBalance.Balance, amountCents); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return transactionCreated, nil
}

func (s *service) getAccountBalance(
	ctx context.Context,
	customerAccountID *uuid.UUID,
	request createTransactionRequest,
) (*models.Balance, error) {
	accountBalance, err := s.balanceRepository.GetCustomerAccountBalance(ctx, customerAccountID)
	if err != nil {
		log.Err(err).
			Str("customer_account_id", customerAccountID.String()).
			Str("idempotency_key", utils.SafeStringPointerValue(request.IdempotencyKey)).
			Float64("amount", request.Amount).
			Msg("failed to get account balance")

		return nil, err
	}

	return accountBalance, nil
}

func (s *service) calculateTransactionAmount(request createTransactionRequest, currentBalance int64) (int64, error) {
	amountCents := utils.ToCents(request.Amount)

	if err := s.isValidOperation(request.OperationType, amountCents, currentBalance); err != nil {
		return 0, err
	}

	amountWithDirection, err := utils.ApplyMoneyDirection(amountCents, request.OperationType)
	if err != nil {
		return 0, err
	}

	return amountWithDirection, nil
}

func (s *service) createTransaction(
	ctx context.Context,
	customerAccount *repository.CustomerAccountByIDResult,
	request createTransactionRequest,
	amountCents int64,
) (*models.Transaction, error) {
	transaction, err := s.transactionRepository.CreateTransaction(ctx, models.Transaction{
		CustomerAccountID: customerAccount.ID,
		OperationType:     request.OperationType,
		Amount:            amountCents,
		IdempotencyKey:    request.IdempotencyKey,
	})
	if err != nil {
		log.Err(err).
			Str("customer_account_id", customerAccount.ID.String()).
			Str("idempotency_key", utils.SafeStringPointerValue(request.IdempotencyKey)).
			Int64("amount", amountCents).
			Msg("failed to create transaction")

		return nil, err
	}

	return transaction, nil
}

func (s *service) updateBalance(
	ctx context.Context,
	customerAccount *repository.CustomerAccountByIDResult,
	currentBalance int64,
	amountCents int64,
) error {
	newBalance := currentBalance + amountCents

	_, err := s.balanceRepository.UpdateCustomerAccountBalance(ctx, models.Balance{
		CustomerAccountID: customerAccount.ID,
		Balance:           newBalance,
	})
	if err != nil {
		log.Err(err).
			Str("customer_account_id", customerAccount.ID.String()).
			Int64("new_balance", newBalance).
			Msg("failed to update customer balance")

		return err
	}

	return nil
}

func (s *service) isValidOperation(operation models.OperationType, amount int64, customerAccountBalance int64) error {
	if s.isCreditOperation(operation) {
		return nil
	}

	if amount > customerAccountBalance {
		return cerror.New(cerror.Params{
			Status:  http.StatusBadRequest,
			Message: "Insufficient funds to perform operation",
		})
	}

	return nil
}

func (s *service) isCreditOperation(operation models.OperationType) bool {
	return operation == models.CreditVoucher
}
