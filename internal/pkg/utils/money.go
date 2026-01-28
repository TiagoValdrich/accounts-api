package utils

import (
	"net/http"

	"github.com/shopspring/decimal"
	"github.com/tiagovaldrich/accounts-api/internal/models"
	"github.com/tiagovaldrich/accounts-api/internal/pkg/cerror"
)

func positiveAmount(amount int64) int64 {
	return amount
}

func negativeAmount(amount int64) int64 {
	return amount * -1
}

func ApplyMoneyDirection(amount int64, operation models.OperationType) (int64, error) {
	formatterMap := map[models.OperationType]func(amount int64) int64{
		models.CreditVoucher:             positiveAmount,
		models.NormalPurchase:            negativeAmount,
		models.PurcharseWithInstallments: negativeAmount,
		models.Withdrawal:                negativeAmount,
	}

	moneyFormatter, ok := formatterMap[operation]
	if !ok {
		return 0, cerror.New(cerror.Params{
			Status:  http.StatusInternalServerError,
			Message: "Unsupported operation",
		})
	}

	return moneyFormatter(amount), nil
}

func ToCents(amount float64) int64 {
	return decimal.NewFromFloat(amount).Mul(decimal.NewFromFloat(100)).IntPart()
}

func FromCents(cents int64) float64 {
	result, _ := decimal.NewFromInt(cents).Div(decimal.NewFromInt(100)).Float64()
	return result
}
