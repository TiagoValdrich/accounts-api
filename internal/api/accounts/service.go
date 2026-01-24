package accounts

import (
	"context"

	"github.com/paemuri/brdoc"
	"github.com/tiagovaldrich/accounts-api/internal/models"
	"github.com/tiagovaldrich/accounts-api/internal/pkg/cerror"
	"github.com/tiagovaldrich/accounts-api/internal/repository"
)

type Servicer interface {
	CreateAccount(context.Context, CreateAccountRequest) (CustomerAccountResult, error)
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

func (s *service) CreateAccount(ctx context.Context, accountReq CreateAccountRequest) (CustomerAccountResult, error) {
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
			return err
		}

		customerAccount, err := s.customerAccountRepository.CreateCustomerAccount(txCtx, models.CustomerAccount{
			CustomerID: customer.ID,
		})

		customerAccountResult.Customer = customer
		customerAccountResult.CustomerAccount = customerAccount

		return err
	})

	if err != nil {
		return customerAccountResult, err
	}

	return customerAccountResult, nil
}

func (s *service) isValidDocumentNumber(documentNumber string) bool {
	return brdoc.IsCPF(documentNumber) || brdoc.IsCNPJ(documentNumber)
}
