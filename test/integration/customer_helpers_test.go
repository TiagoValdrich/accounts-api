package integration

import (
	"context"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tiagovaldrich/accounts-api/internal/models"
)

func AssertCustomerExists(t *testing.T, document string) models.Customer {
	t.Helper()

	var customer models.Customer
	err := DB.NewSelect().
		Model(&customer).
		Where("document = ?", document).
		Scan(context.Background())

	require.NoError(t, err, "customer with document %s should exist", document)
	assert.NotNil(t, customer.ID)
	assert.Equal(t, document, customer.Document)

	return customer
}

func AssertCustomerAccountExists(t *testing.T, customerID uuid.UUID) models.CustomerAccount {
	t.Helper()

	var account models.CustomerAccount
	err := DB.NewSelect().
		Model(&account).
		Where("customer_id = ?", customerID).
		Scan(context.Background())

	require.NoError(t, err, "customer account for customer %s should exist", customerID)
	assert.NotNil(t, account.ID)

	return account
}

func AssertCustomerAccountExistsByID(t *testing.T, accountID string) models.CustomerAccount {
	t.Helper()

	var account models.CustomerAccount
	err := DB.NewSelect().
		Model(&account).
		Where("id = ?", accountID).
		Scan(context.Background())

	require.NoError(t, err, "customer account with ID %s should exist", accountID)

	return account
}
