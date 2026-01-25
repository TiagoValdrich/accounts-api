package accounts

import (
	"context"

	"github.com/paemuri/brdoc"
	"github.com/rs/zerolog/log"
	"github.com/tiagovaldrich/accounts-api/internal/models"
	"github.com/tiagovaldrich/accounts-api/internal/pkg/cerror"
	"github.com/tiagovaldrich/accounts-api/internal/repository"
)

type Servicer interface {
	CreateAccount(context.Context, createAccountRequest) (CustomerAccountResult, error)
	SearchCustomerAccountByID(
		ctx context.Context, req searchAccountRequest,
	) (SearchCustomerAccountResult, error)
}

type service struct {
	customerRepository        repository.CustomerRepository
	customerAccountRepository repository.CustomerAccountRepository
}

func NewService(
	customerRepository repository.CustomerRepository,
	customerAccountRepository repository.CustomerAccountRepository,
) Servicer {
	return &service{
		customerRepository:        customerRepository,
		customerAccountRepository: customerAccountRepository,
	}
}

func (s *service) CreateAccount(ctx context.Context, accountReq createAccountRequest) (CustomerAccountResult, error) {
	var customerAccountResult CustomerAccountResult

	if !s.isValidDocumentNumber(accountReq.Document) {
		return customerAccountResult, cerror.New(cerror.Params{
			Status:  400,
			Message: "Invalid document",
		})
	}

	err := s.customerRepository.WithTransaction(ctx, func(txCtx context.Context) error {
		customer, err := s.customerRepository.CreateCustomer(txCtx, models.Customer{
			Document: accountReq.Document,
		})
		if err != nil {
			log.Err(err).Msg("failed to create customer")

			return err
		}

		customerAccount, err := s.customerAccountRepository.CreateCustomerAccount(txCtx, models.CustomerAccount{
			CustomerID: customer.ID,
		})
		if err != nil {
			log.Err(err).
				Str("customer_id", customerAccount.ID.String()).
				Msg("failed to create customer account")

			return err
		}

		customerAccountResult.Customer = customer
		customerAccountResult.CustomerAccount = customerAccount

		return nil
	})

	if err != nil {
		return customerAccountResult, err
	}

	return customerAccountResult, nil
}

func (s *service) isValidDocumentNumber(documentNumber string) bool {
	return brdoc.IsCPF(documentNumber) || brdoc.IsCNPJ(documentNumber)
}

func (s *service) SearchCustomerAccountByID(
	ctx context.Context, searchAccountReq searchAccountRequest,
) (SearchCustomerAccountResult, error) {
	customerAccount, err := s.customerAccountRepository.SearchCustomerAccountByID(ctx, searchAccountReq.CustomerAccountID)
	if err != nil {
		log.Err(err).
			Str("customer_account_id", searchAccountReq.CustomerAccountID.String()).
			Msg("failed to search for customer account")

		return SearchCustomerAccountResult{}, err
	}

	if customerAccount == nil {
		return SearchCustomerAccountResult{}, cerror.New(cerror.Params{
			Status:  404,
			Message: "Customer account not found",
		})
	}

	return DatabaseToSearchCustomerAccountResult(*customerAccount), nil
}
