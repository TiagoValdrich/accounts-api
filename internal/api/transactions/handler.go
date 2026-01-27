package transactions

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/tiagovaldrich/accounts-api/internal/pkg/cerror"
)

type httpHandler struct {
	service Servicer
}

func NewHTTPHandler(app *fiber.App, service Servicer) {
	httpHandler := &httpHandler{
		service: service,
	}

	routeGroup := app.Group("/transactions")
	routeGroup.Post("/", httpHandler.createTransaction)
}

func (h *httpHandler) createTransaction(c *fiber.Ctx) error {
	var createTransactionReq createTransactionRequest

	if err := c.BodyParser(&createTransactionReq); err != nil {
		return cerror.New(cerror.Params{
			Status:  http.StatusBadRequest,
			Message: "Invalid transaction payload",
		})
	}

	createTransactionResult, err := h.service.CreateTransaction(c.Context(), createTransactionReq)
	if err != nil {
		log.Err(err).Msg("failed to create transaction")

		return err
	}

	return c.Status(http.StatusOK).JSON(DomainToCreateTransactionResponse(createTransactionResult))
}
